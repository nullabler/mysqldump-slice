package repository

import (
	"context"
	"database/sql"
	"fmt"
	"mysqldump-slice/entity"
	"strings"
	"time"
)

type DbInterface interface {
	Close()
	IsIntByCol(string, string) (bool, error)
	LoadRelations(entity.CollectInterface) error
	LoadTables(entity.CollectInterface) error
	PrimaryKeys(string) ([]string, error)
	LoadIds(string, bool, Specs, []string, int) ([]*entity.Value, error)
	LoadDeps(string, string, entity.RelationInterface) ([]string, error)
	LoadPkByCol(string, string, []string, []string) ([]*entity.Value, error)

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
	sql := `SELECT data_type  
		FROM information_schema.columns 
		WHERE table_schema = ? 
		AND table_name = ? 
		AND column_name = ?;`
	if err := db.con.QueryRowContext(ctx, sql, db.name, tab, col).Scan(&typeCol); err != nil {
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

	sql := `SELECT fks.table_name as foreign_table,
			fks.referenced_table_name as primary_table,
			kcu.column_name as fk_column,
			kcu.referenced_column_name as ref_column
		FROM information_schema.referential_constraints fks
		JOIN information_schema.key_column_usage kcu
			ON fks.constraint_schema = kcu.table_schema
			AND fks.table_name = kcu.table_name
			AND fks.constraint_name = kcu.constraint_name
		WHERE fks.constraint_schema = ?
		GROUP BY fks.constraint_schema,
			fks.table_name,
			fks.unique_constraint_schema,
			fks.referenced_table_name,
			fks.constraint_name
		ORDER BY fks.constraint_schema, fks.table_name`

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

	rows, err := db.con.QueryContext(ctx, fmt.Sprintf("SHOW FULL TABLES FROM `%s`", db.name))
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
		`SELECT COLUMN_NAME 
		FROM information_schema.KEY_COLUMN_USAGE 
		WHERE TABLE_NAME = ? 
		AND CONSTRAINT_NAME = ?`,
		tabName,
		"PRIMARY",
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

func (db *Db) LoadIds(tabName string, okSpecs bool, specs Specs, prKeyList []string, confLimit int) ([]*entity.Value, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.conf.MaxLifetimeQuery())*time.Second)
	defer cancel()

	var sort string
	if okSpecs && len(specs.Sort) > 0 {
		sort = strings.Join(db.wrapKeys(specs.Sort, "`"), ", ")
	} else {
		sort = strings.Join(db.wrapKeys(prKeyList, "`"), ", ")
	}

	var condition, limit string
	if okSpecs && len(specs.Condition) > 0 {
		condition = "WHERE " + specs.Condition
		if specs.Limit > 0 {
			limit = fmt.Sprintf("LIMIT %d", specs.Limit)
		}
	}

	if len(condition) == 0 && confLimit > 0 {
		limit = fmt.Sprintf("LIMIT %d", confLimit)
		if okSpecs && specs.Limit > 0 {
			limit = fmt.Sprintf("LIMIT %d", specs.Limit)
		}
	}

	list := []*entity.Value{}
	for _, key := range prKeyList {
		rows, err := db.con.QueryContext(ctx, fmt.Sprintf("SELECT `%s` FROM `%s` %s ORDER BY %s DESC %s",
			key, tabName, condition, sort, limit))
		if err != nil {
			return list, err
		}

		IsIntByCol, err := db.IsIntByCol(tabName, key)
		if err != nil {
			return list, err
		}

		for rows.Next() {
			val, err := db.toString(rows, IsIntByCol)
			if err != nil {
				return list, err
			}

			if len(val) > 0 {
				list = append(list, entity.NewValue(key, val))
			}
		}
	}

	return list, nil
}

func (db *Db) LoadDeps(tabName, where string, rel entity.RelationInterface) (list []string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.conf.MaxLifetimeQuery())*time.Second)
	defer cancel()

	if tabName == "sym_order" && rel.RefTab() == "billing_balance_log" {
		fmt.Println("*****************", tabName, rel)
	}

	limit := ""
	if rel.Limit() > 0 {
		limit = fmt.Sprintf("LIMIT %d", rel.Limit())
		if tabName == "sym_order" && rel.Tab() == "billing_balance_log" {
			fmt.Println("===================", limit)
		}
	}

	rows, err := db.con.QueryContext(ctx, fmt.Sprintf("SELECT `%s` FROM `%s` WHERE %s %s",
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
		val, err := db.toString(rows, isIntDep)
		if err != nil {
			break
		}

		if len(val) > 0 {
			list = append(list, val)
		}
	}

	return
}

func (db *Db) LoadPkByCol(tabName, tabCol string, pkList, valList []string) ([]*entity.Value, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(db.conf.MaxLifetimeQuery())*time.Second)
	defer cancel()

	list := []*entity.Value{}
	for _, key := range pkList {
		rows, err := db.con.QueryContext(ctx, fmt.Sprintf("SELECT `%s` FROM `%s` WHERE `%s` IN (%s)",
			key, tabName, tabCol, strings.Join(valList, ", ")))
		if err != nil {
			return list, err
		}

		IsIntByCol, err := db.IsIntByCol(tabName, key)
		if err != nil {
			return list, err
		}

		for rows.Next() {
			val, err := db.toString(rows, IsIntByCol)
			if err != nil {
				return list, err
			}

			if len(val) > 0 {
				list = append(list, entity.NewValue(key, val))
			}
		}
	}

	return list, nil
}

func (db *Db) Sql() SqlInterface {
	return db.sql
}

func (d *Db) wrapKeys(keys []string, wrapSym string) (list []string) {
	for _, key := range keys {
		list = append(list, wrapSym+key+wrapSym)
	}

	return
}

func (db *Db) toString(rows *sql.Rows, isInt bool) (string, error) {
	var id *int
	var uid *string
	var key string

	if isInt {
		if err := rows.Scan(&id); err != nil {
			return "", err
		}

		if id != nil {
			key = fmt.Sprint(*id)
		}
	} else {
		if err := rows.Scan(&uid); err != nil {
			return "", err
		}

		if uid != nil {
			key = fmt.Sprintf("'%s'", *uid)
		}
	}

	return key, nil
}
