package main

import (
	"log"
	"mysqldump-slice/application"
)

func main() {
	log.Println("Start dump")
	app := application.NewApp()
	if err := app.InitApp(); err != nil {
		app.Panic(err)
	}

	defer app.Close()

	if err := app.LoadTables(); err != nil {
		app.Panic(err)
	}

	if err := app.LoadDependence(); err != nil {
		app.Panic(err)
	}

	app.DumpStructDB()
	app.DumpDataDb()

	log.Println("Finish dump")
}
