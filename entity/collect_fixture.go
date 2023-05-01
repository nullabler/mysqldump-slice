package entity

func FillTables(collect *Collect) {
	for _, tabName := range []string{"user", "product", "category", "order"} {
		collect.tables = append(collect.tables, NewTable(tabName))
	}
}

func FillAllRelList(collect *Collect) {
	genRel := func(tab, col, refTab, refCol string, limit int) *Relation {
		rel := NewRelation()
		rel.Load(tab, col, refTab, refCol, limit, false)
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

func FillTabList(collect *Collect) {
	for _, table := range collect.Tables() {
		collect.PushTab(table.Name)
	}

	// User
	user := collect.Tab("user")
	for _, id := range []string{"1", "2", "3"} {
		valList := []*Value{
			NewValue("id", id),
		}
		user.Push(valList)
	}

	// Product
	product := collect.Tab("product")
	for _, id := range []string{"1", "2", "3"} {
		valList := []*Value{
			NewValue("uuid", id),
		}
		product.Push(valList)
	}

	// Category
	category := collect.Tab("product")
	for _, id := range []string{"1", "2", "3"} {
		valList := []*Value{
			NewValue("id", id),
		}
		category.Push(valList)
	}

	// Order
	order := collect.Tab("order")
	for _, id := range []string{"1", "2", "3"} {
		valList := []*Value{
			NewValue("id", id),
		}
		order.Push(valList)
	}
}
