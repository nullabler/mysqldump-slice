package application

import (
	"fmt"
	"strings"
)

func (app *App) DumpStruct() {
	app.ExecDump(fmt.Sprintf("--no-data --routines %s", app.conf.Db()))
}

func (app *App) DumpFullData() {
	app.ExecDump(fmt.Sprintf(
		"--skip-triggers --no-create-info %s %s", 
		app.conf.Db(),
		strings.Join(app.conf.TablesForFullData(), " "),
	))
}

func (app *App) DumpSliceData() {
	for tabName, tab := range app.tables {
		if app.hasTabNameLikeFullData(tabName) {
			continue
		}

		if where := tab.Where(); len(where) > 0 {
			if tabName == "orders" {
				fmt.Println(tabName, where)
			}

			app.ExecDump(fmt.Sprintf(
				"--skip-triggers --no-create-info %s %s --where=\"%s\"",
				app.conf.Db(),
				tabName, 
				where,
			))
		}
	}
}

func (app *App) hasTabNameLikeFullData(val string) (ok bool) {
	for i := range app.conf.TablesForFullData() {
        if ok = app.conf.TablesForFullData()[i] == val; ok {
            return
        }
    }
    return
}
