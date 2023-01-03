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

	app.LoadRelations()
	tableList := app.FindAllTables()
	for _, tabName := range tableList {
		prKeyList := app.PrimaryKeys(tabName)

		app.Collector().PushTable(tabName)	
		app.LoadIds(tabName, prKeyList)
	}

	for _, tabName := range tableList {
		app.LoadDeps(tabName)
	}

	//app.RemoveFile()
	app.DumpStruct()
	app.DumpFullData()

	for _, tabName := range tableList {
		app.DumpSliceData(tabName, app.Collector().WhereId(tabName))
	}

	log.Println("Finish dump")
}

