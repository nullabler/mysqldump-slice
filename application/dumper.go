package application

import (
	"fmt"
	"strings"
)

func (app *App) DumpStruct() {
	app.ExecDump(fmt.Sprintf("--no-data --routines %s", app.conf.Database))
}

func (app *App) DumpFullData() {
	app.ExecDump(fmt.Sprintf(
		"--skip-triggers --no-create-info %s %s", 
		app.conf.Database,
		strings.Join(app.conf.Tables.Full, " "),
	))
}

func (app *App) DumpSliceData(tabName, where string) {
		if app.hasTabNameLikeFullData(tabName) {
			return
		}

		if len(where) > 0 {
			if tabName == "orders" {
				fmt.Println(tabName, where)
			}

			app.ExecDump(fmt.Sprintf(
				"--skip-triggers --no-create-info %s %s --where=\"%s\"",
				app.conf.Database,
				tabName, 
				where,
			))
		}
}

func (app *App) hasTabNameLikeFullData(val string) (ok bool) {
	for i := range app.conf.Tables.Full {
        if ok = app.conf.Tables.Full[i] == val; ok {
            return
        }
    }
    return
}
