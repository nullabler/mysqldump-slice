package application

import "fmt"

func (app *App) isIntByCol(tab, col string) (bool, error) {
	var typeCol string
	sql := `SELECT data_type  
		FROM information_schema.columns 
		WHERE table_schema = ? 
		AND table_name = ? 
		AND column_name = ?;`
	err := app.db.QueryRow(sql, app.conf.Db(), tab, col).Scan(&typeCol)
	if err != nil {
		return false, err
	}

	if typeCol == "int" {
		return true, nil
	}

	return false, nil
}

func (app *App) colExist(tab, col string) bool {
	sql := fmt.Sprintf(`SHOW columns FROM %s LIKE '%s'`, tab, col)
	var a,b,c,d,f,g interface{} 
	_ = app.db.QueryRow(sql).Scan(&a, &b, &c, &d, &f, &g)
	return a != nil
}
