package entity

func FillTables(collect *Collect) {
	for _, tabName := range []string{"user", "product", "category"} {
		collect.tables = append(collect.tables, NewTable(tabName))
	}
}
