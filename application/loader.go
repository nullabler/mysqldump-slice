package application

import (
	"fmt"
	"strings"
)

func (app *App) LoadIds(tabName string, prKeyList []string) {
	specs, ok := app.conf.Specs(tabName)
	if len(prKeyList) == 0 {
		if !ok || len(specs.Pk) == 0 {  
			return
		}
		prKeyList = specs.Pk
	}

	var sort string
	if ok && len(specs.Sort) > 0 {
		sort = strings.Join(specs.Sort, ", ")
	} else {
		sort = strings.Join(prKeyList, ", ")
	}

	limit := app.conf.Tables.Limit
	if ok && specs.Limit > 0 {
		limit = specs.Limit
	}

	for _, key := range prKeyList {
		rows, err := app.db.Query(fmt.Sprintf("SELECT %s FROM %s ORDER BY %s DESC LIMIT %d", 
			key, tabName, sort, limit))

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
		whereId, ok := app.Collector().WhereId(tabName);
		if !ok {
			continue
		}

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

