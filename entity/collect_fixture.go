package entity

func FillTables(collect *Collect) {
	for _, tabName := range []string{"user", "product", "category", "order"} {
		collect.tables = append(collect.tables, NewTable(tabName))
	}
}

func FillAllRelList(collect *Collect) {
	genRel := func(tab, col, refTab, refCol string, limit int) *Relation {
		rel := NewRelation()
		rel.Load(tab, col, refTab, refCol, limit)
		return rel
	}

	collect.relList["product"] = append(collect.relList["product"],
		genRel("product", "category_id", "category", "id", 2),
	)

	collect.relList["order"] = append(collect.relList["order"],
		genRel("order", "user_id", "user", "id", 3),
	)
	collect.relList["order"] = append(collect.relList["order"],
		genRel("order", "product_id", "product", "uuid", 3),
	)
}
