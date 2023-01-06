package application

import (
	"database/sql"
	"fmt"
	"log"
	"mysqldump-slice/config"
	"mysqldump-slice/entity"
	"mysqldump-slice/service"
	"os"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	conf *config.Conf

	db        *sql.DB
	collector *service.Collector
	relations []entity.Relation
}

func NewApp() *App {
	if len(os.Args) < 2 {
		log.Fatal("Not found yaml file")
	}

	f, err := os.CreateTemp("", "")
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	app := &App{
		conf: config.NewConf(os.Args[1], f.Name()),
	}

	app.collector = service.NewCollector()

	return app
}

func (app *App) InitApp() (err error) {
	app.db, err = sql.Open("mysql", app.conf.GetDbUrl())
	return
}

func (app *App) Close() {
	app.db.Close()
	os.Remove(app.conf.Tmp)
}

func (app *App) Panic(err error) {
	panic(err.Error())
}

func (app *App) ExecDump(call string) {
	app.exec(fmt.Sprintf(
		"mysqldump -u%s -p%s -h %s %s >> %s",
		app.conf.User,
		app.conf.Password,
		app.conf.Host,
		call,
		app.conf.Tmp,
	))
}

func (app *App) RemoveFile() {
	app.exec(fmt.Sprintf("rm -f %s 2> /dev/null", app.conf.Filename()))
}

func (app *App) Save() {
	action := "cp %s %s"
	if app.conf.File.Gzip {
		action = "cat %s | gzip > %s.gz"
	}

	app.exec(fmt.Sprintf(
		action,
		app.conf.Tmp,
		app.conf.Filename(),
	))
}

func (app *App) Collector() *service.Collector {
	return app.collector
}

func (app *App) exec(call string) {
	cmd := exec.Command(app.conf.Shell(), "-c", call)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		app.dd(call)
		app.Panic(err)
	}

}

func (app *App) dd(data ...interface{}) {
	fmt.Printf("%+v\n", data)
}
