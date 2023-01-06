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

	var condition, limit string
	if ok && len(specs.Condition) > 0 {
		condition = "WHERE " + specs.Condition
		if specs.Limit > 0 {
			limit = fmt.Sprintf("LIMIT %d", specs.Limit)
		}
	}

	if len(condition) == 0 {
		limit = fmt.Sprintf("LIMIT %d", app.conf.Tables.Limit)
		if ok && specs.Limit > 0 {
			limit = fmt.Sprintf("LIMIT %d", specs.Limit)
		}
	}

	for _, key := range prKeyList {
		rows, err := app.db.Query(fmt.Sprintf("SELECT %s FROM %s %s ORDER BY %s DESC %s",
			key, tabName, condition, sort, limit))

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
		whereId, ok := app.Collector().WhereId(tabName)
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
