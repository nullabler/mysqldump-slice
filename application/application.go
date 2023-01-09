package application

import (
	"fmt"
	"mysqldump-slice/entity"
	"mysqldump-slice/repository"
	"mysqldump-slice/service"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	loader    *service.Loader
	dumper    *service.Dumper
	normalize *service.Normalize
}

func NewApp(conf *repository.Conf, db *repository.Db, cli *repository.Cli) *App {
	return &App{
		loader:    service.NewLoader(conf, db, cli),
		dumper:    service.NewDumper(conf, cli, db),
		normalize: service.NewNormalize(),
	}
}

func (app *App) Run() {
	collect := entity.NewCollect()
	if err := app.loader.Relations(collect); err != nil {
		app.Panic(err)
	}

	if err := app.loader.Tables(collect); err != nil {
		app.Panic(err)
	}

	if err := app.loader.Dependences(collect); err != nil {
		app.Panic(err)
	}

	if err := app.dumper.RmFile(); err != nil {
		app.Panic(err)
	}

	if err := app.dumper.Struct(); err != nil {
		app.Panic(err)
	}

	if err := app.dumper.Full(); err != nil {
		app.Panic(err)
	}

	if err := app.dumper.Slice(collect); err != nil {
		app.Panic(err)
	}

	if err := app.dumper.Save(); err != nil {
		app.Panic(err)
	}
}

func (app *App) Panic(err error) {
	panic(err.Error())
}

func (app *App) dd(data ...interface{}) {
	fmt.Printf("%+v\n", data)
}
