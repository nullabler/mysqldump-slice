package application

import (
	"fmt"
	"mysqldump-slice/relationship"
)

func (app *App) LoadTables() (err error) {
	rows, err := app.db.Query(fmt.Sprintf(`SHOW FULL TABLES FROM %s`, app.conf.Db()))
	if err != nil {
		return
	}

	for rows.Next() {
		var tabName, tabType string
		err = rows.Scan(&tabName, &tabType)

		if err != nil {
			return
		}

		if tabType != "BASE TABLE" {
			continue
		}

		app.tables[tabName] = relationship.NewTable()
	}

	return nil
}

func (app *App) LoadDependence() (err error) {
	if err = app.getRelations(); err != nil {
		return
	}

	for _, rel := range app.relations {
		err = app.addDependence(rel)
		if err != nil {
			return
		}
		return
	}

	return
}

func (app *App) addDependence(rel relationship.Relation) (err error) {
	if !app.colExist(rel.FrTab(), "id") {
		return
	}

	isInt, err := app.isIntByCol(rel.FrTab(), "id")
	if err != nil {
		return
	}

	sqlOrder := ""
	if isInt {
		sqlOrder += "ORDER BY id DESC"
	}

	tab := app.tables[rel.FrTab()]
	rows, err := app.db.Query(fmt.Sprintf("SELECT %s FROM %s GROUP BY %s %s LIMIT %d", 
		rel.FkCol(), rel.FrTab(), rel.FkCol(), sqlOrder, app.conf.Limit()))
	if err != nil {
		return
	}
	
	isInt, err = app.isIntByCol(rel.FrTab(), rel.FkCol())
	if err != nil {
		return
	}

	for rows.Next() {
		err = tab.Parse(rel, isInt, rows)
		if err != nil {
			continue
		}
	}
	app.tables[rel.FrTab()] = tab
	return
}

func (app *App) getRelations() (err error) {
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

	rows, err := app.db.Query(sql, app.conf.Db())
	if err != nil {
		return
	}

	for rows.Next() {
		rel := relationship.Relation{}
		err = rel.Parse(rows)
		if err != nil {
			return
		}
		app.relations = append(app.relations, rel)
	}
	return
}
