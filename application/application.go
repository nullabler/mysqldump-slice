package application

import (
	"database/sql"
	"mysqldump-slice/config"
	"mysqldump-slice/relationship"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	conf *config.Conf

	db       *sql.DB
	tables   map[string]relationship.Table
}

func NewApp() *App {
	return &App{
		conf: config.NewConf(),
		tables: make(map[string]relationship.Table),
	}
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
