package main

import (
	"log"
	"mysqldump-slice/application"
)

func main() {
	log.Println("Start dump")
	app := application.NewApp()
	defer app.Close()
	
	if err := app.InitApp(); err != nil {
		app.Panic(err)
	}

	if err := app.LoadTables(); err != nil {
		app.Panic(err)
	}

	if err := app.LoadDependence(); err != nil {
		app.Panic(err)
	}

	//app.NormilizeTable()

	app.DumpStructDB()
	app.DumpDataDb()

	log.Println("Finish dump")
}
