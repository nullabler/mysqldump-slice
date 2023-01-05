package application

import (
	"fmt"
	"mysqldump-slice/entity"
)

func (app *App) IsIntByCol(tab, col string) (bool, error) {
	var typeCol string
	sql := `SELECT data_type  
		FROM information_schema.columns 
		WHERE table_schema = ? 
		AND table_name = ? 
		AND column_name = ?;`
	if err := app.db.QueryRow(sql, app.conf.Database, tab, col).Scan(&typeCol); err != nil {
		return false, err
	}

	if typeCol == "int" {
		return true, nil
	}

	return false, nil
}

func (app *App) ColExist(tab, col string) bool {
	sql := fmt.Sprintf(`SHOW columns FROM %s LIKE '%s'`, tab, col)
	var a,b,c,d,f,g interface{} 
	_ = app.db.QueryRow(sql).Scan(&a, &b, &c, &d, &f, &g)
	return a != nil
}

func (app *App) LoadRelations() (err error) {
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

	rows, err := app.db.Query(sql, app.conf.Database)
	if err != nil {
		return
	}

	for rows.Next() {
		rel := entity.Relation{}
		if err = rel.Parse(rows); err != nil {
			return
		}

		app.Collector().PushRelation(rel)
	}
	return
}

func (app *App) FindAllTables() (tableList []string) {
	rows, err := app.db.Query(fmt.Sprintf(`SHOW FULL TABLES FROM %s`, app.conf.Database))
	if err != nil {
		return
	}

	for rows.Next() {
		var tabName, tabType string
		if err = rows.Scan(&tabName, &tabType); err != nil {
			return
		}

		if tabType != "BASE TABLE" || app.conf.Ignore(tabName) {
			continue
		}

		tableList = append(tableList, tabName)
	}

	return
}

func (app *App) PrimaryKeys(tabName string) (keyList []string) {
	rows, err := app.db.Query(fmt.Sprintf(`SHOW KEYS FROM %s WHERE Key_name = 'PRIMARY'`, tabName))
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
