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
	app.dd(app.relations)
}
