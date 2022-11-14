package application

import (
	"fmt"
	"html/template"
	"mysqldump-slice/tmpl"
	"os"
)

func (app *App) DumpStructDB() {
	file, _ := os.Create(app.conf.Filename())
	defer file.Close()

	tm, err := template.New("dump").Parse(tmpl.Dump())
	if err != nil {
		panic(err.Error())
	}

	var d, c string
	_ = app.db.QueryRow(fmt.Sprintf("SHOW CREATE DATABASE %s", app.conf.Db())).Scan(&d, &c)

	param := tmpl.Param{
		Host: app.conf.Host(),
		Database: app.conf.Db(),
		CreateDatabase: c,
		ServerVersion: "0.0.1",
	}

	if err = tm.Execute(file, param); err != nil {
		panic(err.Error())
	}

}

func (app *App) DumpDataDb() {

}
