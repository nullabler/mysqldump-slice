package repository

import (
	"context"
	"database/sql"
	"fmt"
	"mysqldump-slice/entity"
	"time"
)

type DbInterface interface {
	Close()
	IsIntByCol(string, string) (bool, error)
	LoadRelations(entity.CollectInterface) error
	LoadTables(entity.CollectInterface) error
	PrimaryKeys(string) ([]string, error)
	LoadIds(string, *Specs, []string) ([][]*entity.Value, error)
	LoadDeps(string, string, entity.RelationInterface) ([]string, error)
	LoadPkByCol(string, string, []string, []string) ([][]*entity.Value, error)

	Sql() SqlInterface
}

type Db struct {
	name string
	con  *sql.DB
	conf *Conf
	sql  SqlInterface

	isClose bool
}

func NewDb(conf *Conf, driver string) (*Db, error) {
	con, err := sql.Open(driver, conf.DbUrl())
	if err != nil {
		return nil, err
	}
	con.SetMaxOpenConns(conf.MaxConnect())
	con.SetMaxIdleConns(conf.MaxConnect())
	con.SetConnMaxLifetime(time.Minute * time.Duration(conf.MaxLifetimeConnect()))

	return &Db{
		name: conf.Database,
		con:  con,
		conf: conf,
		sql:  NewSql(conf),

		isClose: false,
	}, nil
}

func (db *Db) Close() {
	if db.isClose {
		return
	}

	db.isClose = true
	db.con.Close()
}

func (db *Db) IsIntByCol(tab, col string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.conf.MaxLifetimeQuery())*time.Second)
	defer cancel()

	var typeCol string
	if err := db.con.QueryRowContext(ctx, db.Sql().QueryIsIntByCol(), db.name, tab, col).Scan(&typeCol); err != nil {
		return false, err
	}

	if typeCol == "int" {
		return true, nil
	}

	return false, nil
}

func (db *Db) LoadRelations(collect entity.CollectInterface) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.conf.MaxLifetimeQuery())*time.Second)
	defer cancel()

	rows, err := db.con.QueryContext(ctx, db.Sql().QueryRelations(), db.name)
	if err != nil {
		return err
	}

	for rows.Next() {
		rel := entity.NewRelation()
		if err := rel.Parse(rows); err != nil {
			return err
		}

		collect.PushRelation(rel)
	}

	return err
}

func (db *Db) LoadTables(collect entity.CollectInterface) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.conf.MaxLifetimeQuery())*time.Second)
	defer cancel()

	rows, err := db.con.QueryContext(ctx, db.Sql().QueryFullTables(db.name))
	if err != nil {
		return err
	}

	for rows.Next() {
		var tabName, tabType string
		if err = rows.Scan(&tabName, &tabType); err != nil {
			return err
		}

		if tabType != "BASE TABLE" || db.conf.IsIgnore(tabName) {
			continue
		}

		collect.PushTable(tabName)
	}

	return nil
}

func (db *Db) PrimaryKeys(tabName string) ([]string, error) {
	var keyList []string

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.conf.MaxLifetimeQuery())*time.Second)
	defer cancel()

	rows, err := db.con.QueryContext(ctx,
		db.Sql().QueryPrimaryKeys(),
		tabName,
		"PRIMARY",
		db.conf.DbName(),
	)
	if err != nil {
		return keyList, err
	}

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return keyList, err
		}

		keyList = append(keyList, key)
	}

	return keyList, nil
}

func (db *Db) LoadIds(tabName string, specs *Specs, prKeyList []string) ([][]*entity.Value, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.conf.MaxLifetimeQuery())*time.Second)
	defer cancel()

	list := [][]*entity.Value{}
	sql, errSql := db.Sql().QueryLoadIds(tabName, specs, prKeyList)
	if errSql != nil {
		return list, errSql
	}

	rows, err := db.con.QueryContext(ctx, sql)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		valList, err := db.Scan(rows, prKeyList)
		if err != nil {
			return list, err
		}

		if len(valList) > 0 {
			list = append(list, valList)
		}
	}

	return list, nil
}

func (db *Db) LoadDeps(tabName, where string, rel entity.RelationInterface) (list []string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.conf.MaxLifetimeQuery())*time.Second)
	defer cancel()

	limit := ""
	if rel.Limit() > 0 {
		limit = fmt.Sprintf("LIMIT %d", rel.Limit())
	}

	rows, err := db.con.QueryContext(ctx, db.Sql().QueryLoadDeps(
		rel.Col(),
		tabName,
		where,
		limit,
	))
	if err != nil {
		return
	}

	isIntDep, err := db.IsIntByCol(tabName, rel.Col())
	if err != nil {
		return
	}

	for rows.Next() {
		val, err := db.singleScan(rows)
		if err != nil {
			break
		}

		if !isIntDep {
			val = fmt.Sprintf("'%s'", val)
		}

		if len(val) > 0 {
			list = append(list, val)
		}
	}

	return
}

func (db *Db) LoadPkByCol(tabName, tabCol string, pkList, valList []string) ([][]*entity.Value, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.conf.MaxLifetimeQuery())*time.Second)
	defer cancel()

	list := [][]*entity.Value{}
	rows, err := db.con.QueryContext(ctx, db.Sql().QueryLoadPkByCol(pkList, tabName, tabCol, valList))
	if err != nil {
		return list, err
	}

	for rows.Next() {
		valList, err := db.Scan(rows, pkList)
		if err != nil {
			return list, err
		}

		if len(valList) > 0 {
			list = append(list, valList)
		}
	}

	return list, nil
}

func (db *Db) Sql() SqlInterface {
	return db.sql
}

func (db *Db) Scan(rows *sql.Rows, fields []string) ([]*entity.Value, error) {
	length := len(fields)
	list := []*entity.Value{}

	buf := make([]interface{}, length)
	for i := range fields {
		buf[i] = new(sql.RawBytes)
	}

	if err := rows.Scan(buf...); err != nil {
		return list, err
	}

	for i := 0; i < length; i++ {
		if v, ok := buf[i].(*sql.RawBytes); ok {
			list = append(list, entity.NewValue(fields[i], string(*v)))
		}
	}

	return list, nil
}

func (db *Db) singleScan(rows *sql.Rows) (string, error) {
	buf := new(sql.RawBytes)
	if err := rows.Scan(buf); err != nil {
		return "", err
	}

	return string(*buf), nil
}
