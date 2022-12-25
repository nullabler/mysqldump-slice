package application

import (
	"database/sql"
	"fmt"
	"mysqldump-slice/config"
	"mysqldump-slice/relationship"
	"mysqldump-slice/tmpl"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	conf *config.Conf

	db       *sql.DB
	tables   map[string]*relationship.Table
	relations []relationship.Relation
	templateParams *tmpl.TemplateParams
}

func NewApp() *App {
	app := &App{
		conf: config.NewConf(),
		tables: make(map[string]*relationship.Table),
	}

	app.templateParams = tmpl.NewTemplateParams(app.conf.Host(), app.conf.Db())

	return app
}

func (app *App) InitApp() (err error) {
	app.db, err = sql.Open("mysql", app.conf.GetDbUrl())
	return
}

func (app *App) Close() {
	app.db.Close()
}

func (app *App) Panic(err error) {
	panic(err.Error())
}

func (app *App) ExecDump(call string) {
	app.exec(fmt.Sprintf(
		"mysqldump -u%s -p%s -h %s %s >> %s",
		app.conf.User(),
		app.conf.Passwd(),
		app.conf.Host(),
		call,
		app.conf.Filename(),
	))
}

func (app *App) RemoveFile() {
	app.exec(fmt.Sprintf("rm -f %s 2> /dev/null", app.conf.Filename()))
}

func (app *App) exec(call string) {
	cmd := exec.Command(app.conf.Shell(), "-c", call)
	if err := cmd.Run(); err != nil {
		app.Panic(err)
	}
}

func (app *App) getTab(tabName string) *relationship.Table {
	tab := app.tables[tabName]
	if tab == nil {
		tab = relationship.NewTable(tabName)
	}

	return tab 
}

func (app *App) setTab(tab *relationship.Table) {
	app.tables[tab.Name()] = tab
}

func (app *App) dd(data ...interface{}) {
	fmt.Printf("%+v\n", data)
}
