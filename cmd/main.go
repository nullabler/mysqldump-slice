package main

import (
	"log"
	"mysqldump-slice/application"
	"mysqldump-slice/repository"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Not found yaml file")
	}

	f, err := os.CreateTemp("", "")
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	conf, err := repository.NewConf(os.Args[1], f.Name())
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewDb(conf, "mysql")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cli, err := repository.NewCli(conf)
	if err != nil {
		log.Fatal(err)
	}

	app := application.NewApp(conf, db, cli)

	app.Run()
}
