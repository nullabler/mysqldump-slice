package repository

import (
	"database/sql"
	"fmt"
	"mysqldump-slice/entity"
	"strings"
)

type Db struct {
	name string
	con  *sql.DB
	conf *Conf
}

func NewDb(conf *Conf, driver string) (*Db, error) {
	con, err := sql.Open(driver, conf.DbUrl())
	if err != nil {
		return nil, err
	}
	con.SetMaxOpenConns(conf.MaxConnect)

	return &Db{
		name: conf.Database,
		con:  con,
		conf: conf,
	}, nil
}

func (db *Db) Close() {
	db.con.Close()
}

func (db *Db) IsIntByCol(tab, col string) (bool, error) {
	var typeCol string
	sql := `SELECT data_type  
		FROM information_schema.columns 
		WHERE table_schema = ? 
		AND table_name = ? 
		AND column_name = ?;`
	if err := db.con.QueryRow(sql, db.name, tab, col).Scan(&typeCol); err != nil {
		return false, err
	}

	if typeCol == "int" {
		return true, nil
	}

	return false, nil
}

func (db *Db) ColExist(tab, col string) bool {
	sql := fmt.Sprintf(`SHOW columns FROM %s LIKE '%s'`, tab, col)
	var a, b, c, d, f, g interface{}
	_ = db.con.QueryRow(sql).Scan(&a, &b, &c, &d, &f, &g)
	return a != nil
}

func (db *Db) LoadRelations(collect *entity.Collect) error {
	sql := `select fks.table_name as foreign_table,
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

	rows, err := db.con.Query(sql, db.name)
	if err != nil {
		return err
	}

	for rows.Next() {
		rel := entity.Relation{}
		if err := rel.Parse(rows); err != nil {
			return err
		}

		collect.PushRelation(rel)
	}

	return err
}

func (db *Db) LoadTables(collect *entity.Collect) {
	rows, err := db.con.Query(fmt.Sprintf(`SHOW FULL TABLES FROM %s`, db.name))
	if err != nil {
		return
	}

	for rows.Next() {
		var tabName, tabType string
		if err = rows.Scan(&tabName, &tabType); err != nil {
			return
		}

		if tabType != "BASE TABLE" || db.conf.Ignore(tabName) {
			continue
		}

		collect.PushTable(tabName)
	}
}

func (db *Db) PrimaryKeys(tabName string) (keyList []string) {
	rows, err := db.con.Query(fmt.Sprintf(`SHOW KEYS FROM %s WHERE Key_name = 'PRIMARY'`, tabName))
	if err != nil {
		return
	}

	for rows.Next() {
		var t *string
		var key string
		if err = rows.Scan(&t, &t, &t, &t, &key, &t, &t, &t, &t, &t, &t, &t, &t, &t); err != nil {
			return
		}

		keyList = append(keyList, key)
	}

	return
}

func (db *Db) LoadIds(tabName string, collect *entity.Collect, okSpecs bool, specs Specs, prKeyList []string, confLimit int) {
	var sort string
	if okSpecs && len(specs.Sort) > 0 {
		sort = strings.Join(specs.Sort, ", ")
	} else {
		sort = strings.Join(prKeyList, ", ")
	}

	var condition, limit string
	if okSpecs && len(specs.Condition) > 0 {
		condition = "WHERE " + specs.Condition
		if specs.Limit > 0 {
			limit = fmt.Sprintf("LIMIT %d", specs.Limit)
		}
	}

	if len(condition) == 0 {
		limit = fmt.Sprintf("LIMIT %d", confLimit)
		if okSpecs && specs.Limit > 0 {
			limit = fmt.Sprintf("LIMIT %d", specs.Limit)
		}
	}

	for _, key := range prKeyList {
		rows, err := db.con.Query(fmt.Sprintf("SELECT %s FROM %s %s ORDER BY %s DESC %s",
			key, tabName, condition, sort, limit))

		IsIntByCol, errIsIntByCol := db.IsIntByCol(tabName, key)
		if err != nil || errIsIntByCol != nil {
			return
		}

		for rows.Next() {
			collect.PushKey(tabName, key, IsIntByCol, rows)
		}
	}

}

func (db *Db) LoadDeps(tabName string, collect *entity.Collect, rel entity.Relation, keys map[string][]string) {
	rows, err := db.con.Query(fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		rel.Col(),
		tabName,
		db.WhereAll(keys),
	))

	isIntDep, errIsInt := db.IsIntByCol(tabName, rel.Col())
	if err != nil || errIsInt != nil {
		return
	}

	for rows.Next() {
		collect.PushKey(rel.RefTab(), rel.RefCol(), isIntDep, rows)
	}
}

func (db *Db) WhereAll(keys map[string][]string) string {
	var whereList []string
	for _, list := range db.Where(keys) {
		whereList = append(whereList, "("+strings.Join(list, ", ")+")")
	}

	return strings.Join(whereList, " AND ")
}

func (db *Db) WhereSlice(point *Point) []string {
	var where []string
	var query []string
	for n := 0; n < point.Count; n++ {
		col, i := point.Next(n)
		query = append(query, point.Keys[col][i])

		if point.Current == 0 {
			where = append(where, strings.Join(query, " AND "))
			query = []string{}
		}

	}

	return where
}

func (db *Db) Where(keys map[string][]string) map[string][]string {
	where := make(map[string][]string)
	for col, valList := range keys {
		limit := 0
		for start := 0; start < len(valList); start += db.conf.LimitCli {
			limit += db.conf.LimitCli
			if limit > len(valList) {
				limit = len(valList)
			}

			where[col] = append(where[col], fmt.Sprintf("%s IN (%s)", col, strings.Join(valList[start:limit], ", ")))
		}
	}

	return where
}
