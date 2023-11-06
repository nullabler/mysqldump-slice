package main

import (
	"log"
	"mysqldump-slice/addapter"
	"mysqldump-slice/application"
	"mysqldump-slice/config"
	"mysqldump-slice/module"
	"mysqldump-slice/repository"
	"os"
)

var Version = "v1.0.1-stable"

func main() {
	if len(os.Args) < 2 {
		log.Printf("Slice version: %s \n", Version)
		return
	}

	f, err := os.CreateTemp("", "")
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	conf, err := config.NewConf(Version, os.Args[1], f.Name())
	if err != nil {
		log.Fatal(err)
	}

	lg := module.NewLog(conf)

	db, err := repository.NewDb(conf, "mysql", lg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cli, err := repository.NewCli(conf, addapter.NewExec(conf.Shell()))
	if err != nil {
		log.Fatal(err)
	}

	app := application.NewApp(conf, lg, db, cli)

	app.Run()
}
