package application

import (
	"fmt"
	"mysqldump-slice/relationship"
)

func (app *App) LoadTables() (err error) {
	res, err := app.db.Query(`SHOW tables`)
	if err != nil {
		return
	}

	for res.Next() {
		var name string
		err = res.Scan(&name)
		if err != nil {
			return
		}
		app.tables[name] = relationship.Table{}
	}
	return nil
}

func (app *App) LoadDependence() (err error) {
	relations, err := app.getRelations()
	for _, rel := range relations {
		err = app.addDependence(rel.FrTab(), rel.FkCol())
		if err != nil {
			return
		}
		return
	}
	return
}

func (app *App) addDependence(frTab, fkCol string) (err error) {
	if !app.colExist(frTab, "id") {
		return
	}

	isInt, err := app.isIntByCol(frTab, "id")
	if err != nil {
		return
	}

	sqlOrder := ""
	if isInt {
		sqlOrder += "ORDER BY id DESC"
	}

	tab := app.tables[frTab]
	rows, err := app.db.Query(fmt.Sprintf("SELECT %s FROM %s GROUP BY %s %s LIMIT %d", 
		fkCol, frTab, fkCol, sqlOrder, app.conf.Limit()))
	if err != nil {
		return
	}
	
	isInt, err = app.isIntByCol(frTab, fkCol)
	if err != nil {
		return
	}

	for rows.Next() {
		err = tab.Parse(isInt, rows)
		if err != nil {
			continue
		}
	}
	app.tables[frTab] = tab
	return
}

func (app *App) getRelations() (list []relationship.Relation, err error) {
	sql := `select fks.table_name as foreign_table,
			fks.referenced_table_name as primary_table,
			kcu.column_name as fk_column
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
		list = append(list, rel)
	}
	return
}
