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

		app.tables[tabName] = relationship.NewTable(tabName)
	}

	return nil
}

func (app *App) LoadDependence() (err error) {
	if err = app.getRelations(); err != nil {
		return
	}

	for _, rel := range app.relations {
		app.collectMetadataForTable(rel)

		//err = app.addDependence(rel)
		//if err != nil {
			//return
		//}
	}

	for _, tab := range app.tables {
		app.loadIds(tab)
		app.loadRefs(tab)
		if tab.Name() == "orders" {
			fmt.Println(tab)
		}
	}

	return
}

func (app *App) collectMetadataForTable(rel relationship.Relation) {
	tab := app.getTab(rel.Tab())
	tab.PushRelation(rel)
	//app.setTab(tab)
}

func (app *App) loadIds(tab *relationship.Table) {
	if !app.colExist(tab.Name(), "id") {
		return
	}

	if tab.IsIntId() == nil {
		isInt, err := app.isIntByCol(tab.Name(), "id")
		if err != nil {
			return
		}

		tab.SetIsIntId(isInt)
	}

	rows, err := app.db.Query(fmt.Sprintf("SELECT id FROM %s LIMIT %d", 
		tab.Name(), app.conf.Limit()))

	if err != nil {
		return 
	}

	for rows.Next() {
		if err = tab.ParseId(rows); err != nil {
			continue
		}
	}
}

func (app *App) loadRefs(tab *relationship.Table) {
	for _, rel := range tab.RelationList () { 
		whereId, ok := tab.WhereId();
		if !ok {
			return	
		}
		
		rows, err := app.db.Query(fmt.Sprintf("SELECT %s FROM %s WHERE id IN (%s)", 
			rel.Col(), tab.Name(), whereId))

		if err != nil {
			continue 
		}

		isInt, err := app.isIntByCol(tab.Name(), rel.Col())
		if err != nil {
			continue	
		}

		var refId *int
		var refUid *string

		for rows.Next() {
			if isInt {
				err = relationship.ParseRelation(refId, rows)
			} else {
				err = relationship.ParseRelation(refUid, rows)
			}

			if err != nil {
				continue
			}

			app.getTab(rel.RefTab()).PushDep(rel.RefCol(), isInt, *refId, *refUid)
		}
	}

}

func (app *App) addDependence(rel relationship.Relation) (err error) {
	if !app.colExist(rel.Tab(), "id") {
		return
	}

	isInt, err := app.isIntByCol(rel.Tab(), "id")
	if err != nil {
		return
	}

	sqlOrder := ""
	if isInt {
		sqlOrder += "ORDER BY id DESC"
	}

	tab := app.tables[rel.Tab()]
	
	tab.PushRef(rel.Col())

	rows, err := app.db.Query(fmt.Sprintf("SELECT %s FROM %s GROUP BY %s %s LIMIT %d", 
		rel.Col(), rel.Tab(), rel.Col(), sqlOrder, app.conf.Limit()))

	if rel.Tab() == "orders" {
		app.dd(rel)
	}

	if err != nil {
		return
	}
	
	isInt, err = app.isIntByCol(rel.Tab(), rel.Col())
	if err != nil {
		return
	}

	for rows.Next() {
		err = tab.ParseOld(rel, isInt, rows)
		if err != nil {
			continue
		}
	}
	app.tables[rel.Tab()] = tab
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
