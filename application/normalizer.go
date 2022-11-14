package application

import "fmt"

func (app *App) NormilizeTable() {
	for tab, data := range app.tables {
		fmt.Println(tab, data)
	}
}
