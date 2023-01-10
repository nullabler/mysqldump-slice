package application

import (
	"log"
	"mysqldump-slice/entity"
	"mysqldump-slice/repository"
	"mysqldump-slice/service"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	loader *service.Loader
	dumper *service.Dumper
}

func NewApp(conf *repository.Conf, db *repository.Db, cli *repository.Cli) *App {
	return &App{
		loader: service.NewLoader(conf, db, cli),
		dumper: service.NewDumper(conf, cli, db),
	}
}

func (app *App) Run() {
	collect := entity.NewCollect()
	log.Println("Load relations......Start")
	if err := app.loader.Relations(collect); err != nil {
		app.Panic(err)
	}
	log.Println("Load relations......Done")

	log.Println("Load tables......Start")
	if err := app.loader.Tables(collect); err != nil {
		app.Panic(err)
	}
	log.Println("Load tables......Done")

	log.Println("Sort......Start")
	if err := app.loader.Weight(collect); err != nil {
		app.Panic(err)
	}

	service.CallNormalize(collect)
	log.Println("Sort......Done")

	log.Println("Load dependences......Start")
	if err := app.loader.Dependences(collect); err != nil {
		app.Panic(err)
	}
	log.Println("Load dependences......Done")

	if err := app.dumper.RmFile(); err != nil {
		app.Panic(err)
	}

	log.Println("Dump struct......Start")
	if err := app.dumper.Struct(); err != nil {
		app.Panic(err)
	}
	log.Println("Dump struct......Done")

	log.Println("Dump data like full......Start")
	if err := app.dumper.Full(); err != nil {
		app.Panic(err)
	}
	log.Println("Dump data like full......Done")

	log.Println("Dump data like short......Start")
	if err := app.dumper.Slice(collect); err != nil {
		app.Panic(err)
	}
	log.Println("Dump data like short......Done")

	if err := app.dumper.Save(); err != nil {
		app.Panic(err)
	}
}

func (app *App) Panic(err error) {
	panic(err.Error())
}
