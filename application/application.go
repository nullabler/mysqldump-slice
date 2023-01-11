package application

import (
	"mysqldump-slice/entity"
	"mysqldump-slice/repository"
	"mysqldump-slice/service"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	loader *service.Loader
	dumper *service.Dumper
	log    *service.Log
}

func NewApp(conf *repository.Conf, db *repository.Db, cli *repository.Cli) *App {
	log := service.NewLog(conf)

	return &App{
		loader: service.NewLoader(conf, db, cli, log),
		dumper: service.NewDumper(conf, cli, db, log),
		log:    log,
	}
}

func (app *App) Run() {
	collect := entity.NewCollect()
	app.log.Info("Load relations......Start")
	if err := app.loader.Relations(collect); err != nil {
		app.log.Error(err)
	}
	app.log.Info("Load relations......Done")

	app.log.Info("Load tables......Start")
	if err := app.loader.Tables(collect); err != nil {
		app.log.Error(err)
	}
	app.log.Info("Load tables......Done")

	app.log.Info("Sort......Start")
	if err := app.loader.Weight(collect); err != nil {
		app.log.Error(err)
	}

	service.CallNormalize(collect)
	app.log.Info("Sort......Done")

	app.log.Info("Load dependences......Start")
	if err := app.loader.Dependences(collect); err != nil {
		app.log.Error(err)
	}
	app.log.Info("Load dependences......Done")

	if err := app.dumper.RmFile(); err != nil {
		app.log.Error(err)
	}

	app.log.Info("Dump struct......Start")
	if err := app.dumper.Struct(); err != nil {
		app.log.Error(err)
	}
	app.log.Info("Dump struct......Done")

	app.log.Info("Dump data like full......Start")
	if err := app.dumper.Full(); err != nil {
		app.log.Error(err)
	}
	app.log.Info("Dump data like full......Done")

	app.log.Info("Dump data like short......Start")
	if err := app.dumper.Slice(collect); err != nil {
		app.log.Error(err)
	}
	app.log.Info("Dump data like short......Done")

	filename, err := app.dumper.Save()
	if err != nil {
		app.log.Error(err)
	}
	app.log.Printf("Save dump: %s......Done", filename)
}
