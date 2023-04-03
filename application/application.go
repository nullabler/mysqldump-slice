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
	log    service.LogInterface
}

func NewApp(conf *repository.Conf, db repository.DbInterface, cli repository.CliInterface) *App {
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

	app.log.Info("Sort......Start")
	if err := app.loader.Weight(collect); err != nil {
		app.log.Error(err)
	}

	service.CallNormalize(collect)
	app.log.Info("Sort......Done")

	app.log.Info("Load tables......Start")
	if err := app.loader.Tables(collect); err != nil {
		app.log.Error(err)
	}
	app.log.Info("Load tables......Done")

	app.log.Info("Load relations for leader flag......Start")
	app.loader.LoadRelationsForLeader(collect)
	app.log.Info("Load relations for leader flag......Done")

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

	app.log.Info("Load dependences and dump data like short......Start")
	app.runSlice(collect)
	app.log.Info("Load dependences and dump data like short......Done")

	filename, err := app.dumper.Save()
	if err != nil {
		app.log.Error(err)
	}

	app.log.State(filename)
}

func (app *App) runSlice(collect *entity.Collect) {
	isLoop := true
	for {
		if isLoop {
			isLoop = false
		} else {
			break
		}

		for _, table := range collect.Tables() {
			rows := collect.Tab(table.Name).Pull()
			if len(rows) == 0 {
				continue
			} else {
				isLoop = true
			}

			app.log.Infof("- %s......Loading", table.Name)

			for _, rel := range collect.RelList(table.Name) {
				if err := app.loader.Dependences(collect, rel, table.Name, rows); err != nil {
					app.log.Error(err)
				}
			}

			if err := app.dumper.Slice(collect, table.Name, rows); err != nil {
				app.log.Error(err)
			}

			app.log.Infof("- %s......Done", table.Name)
		}
	}
}
