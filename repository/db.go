package repository

import (
	"context"
	"database/sql"
	"fmt"
	"mysqldump-slice/config"
	"mysqldump-slice/entity"
	"mysqldump-slice/module"
	"time"
)

type DbInterface interface {
	Close()
	LoadRelations(entity.CollectInterface) error
	LoadTables(entity.CollectInterface) error
	PrimaryKeys(string) ([]string, error)
	LoadIds(string, *config.Specs, []string) ([]entity.ValList, error)
	LoadDeps(string, string, entity.RelationInterface) ([]string, error)
	LoadPkByCol(string, string, []string, []string, bool) ([]entity.ValList, error)

	Sql() SqlInterface
}

type Db struct {
	name string
	con  *sql.DB
	conf *config.Conf
	sql  SqlInterface
	log  module.LogInterface

	isClose bool
}

func NewDb(conf *config.Conf, driver string, log module.LogInterface) (*Db, error) {
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
		log:  log,

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

func (db *Db) LoadRelations(collect entity.CollectInterface) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.conf.MaxLifetimeQuery())*time.Second)
	defer cancel()

	sql := db.Sql().QueryRelations()
	db.log.Debug("LoadRelations", sql)
	rows, err := db.con.QueryContext(ctx, sql, db.name)

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

	sql := db.Sql().QueryFullTables(db.name)
	db.log.Debug("LoadTables", sql)

	rows, err := db.con.QueryContext(ctx, sql)
	if err != nil {
		return err
	}

	for rows.Next() {
		var tabName, tabType string
		if err = rows.Scan(&tabName, &tabType); err != nil {
			return err
		}

		if tabType != "BASE TABLE" {
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

	sql := db.Sql().QueryPrimaryKeys()
	db.log.Debug("PrimaryKeys", sql)

	rows, err := db.con.QueryContext(ctx,
		sql,
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

func (db *Db) LoadIds(tabName string, specs *config.Specs, prKeyList []string) ([]entity.ValList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.conf.MaxLifetimeQuery())*time.Second)
	defer cancel()

	list := []entity.ValList{}
	sql, errSql := db.Sql().QueryLoadIds(tabName, specs, prKeyList)
	db.log.Debug("LoadIds", sql)

	if errSql != nil {
		return list, errSql
	}

	rows, err := db.con.QueryContext(ctx, sql)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		valList, err := db.Scan(rows)
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
	spec, ok := db.conf.Specs(tabName)
	if ok {
		if spec.DepLimit > 0 {
			limit = fmt.Sprintf("LIMIT %d", spec.DepLimit)
		}
	}

	sql := db.Sql().QueryLoadDeps(
		rel.Col(),
		tabName,
		where,
		limit,
	)

	db.log.Debug("LoadDeps", sql)

	rows, err := db.con.QueryContext(ctx, sql)
	if err != nil {
		return
	}

	for rows.Next() {
		val, err := db.singleScan(rows)
		if err != nil {
			break
		}

		if len(val) > 0 {
			list = append(list, val)
		}
	}

	return
}

func (db *Db) LoadPkByCol(tabName, tabCol string, pkList, valList []string, isGreedy bool) ([]entity.ValList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.conf.MaxLifetimeQuery())*time.Second)
	defer cancel()

	list := []entity.ValList{}
	sql := db.Sql().QueryLoadPkByCol(pkList, tabName, tabCol, valList, isGreedy)
	db.log.Debug("LoadPkByCol", sql)

	rows, err := db.con.QueryContext(ctx, sql)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		valList, err := db.Scan(rows)
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

func (db *Db) Scan(rows *sql.Rows) ([]*entity.Value, error) {
	list := []*entity.Value{}

	cols, err := rows.Columns()
	if err != nil {
		return list, err
	}

	length := len(cols)
	buf := make([]interface{}, length)

	for i := range cols {
		buf[i] = new(sql.RawBytes)
	}

	if err := rows.Scan(buf...); err != nil {
		return list, err
	}

	for i := 0; i < length; i++ {
		if v, ok := buf[i].(*sql.RawBytes); ok {
			list = append(list, entity.NewValue(cols[i], string(*v)))
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
