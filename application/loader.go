package application

import (
	"fmt"
)

func (app *App) LoadIds(tabName string, prKeyList []string) {
	if len(prKeyList) == 0 {
		return
	}

	for _, key := range prKeyList {
		rows, err := app.db.Query(fmt.Sprintf("SELECT %s FROM %s LIMIT %d", 
			key, tabName, app.conf.Tables.Limit))

		IsIntByCol, errIsIntByCol := app.IsIntByCol(tabName, key) 
		if err != nil || errIsIntByCol != nil {
			return 
		}

		for rows.Next() {
			app.Collector().ParseId(tabName, key, IsIntByCol, rows)
		}
	}
}

func (app *App) LoadDeps(tabName string) {
	for _, rel := range app.Collector().RelationList(tabName) { 
		whereId := app.Collector().WhereId(tabName);
		
		rows, err := app.db.Query(fmt.Sprintf("SELECT %s FROM %s WHERE %s", 
			rel.Col(), tabName, whereId))

		isIntDep, errIsInt := app.IsIntByCol(tabName, rel.Col())
		if err != nil || errIsInt != nil {
			continue 
		}

		for rows.Next() {
			app.Collector().PushDep(rel.RefTab(), isIntDep, rel, rows)
		}
	}
}

