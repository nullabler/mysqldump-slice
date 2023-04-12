package config

func FillSimpleSpecs(conf *Conf) {
	conf.Tables.Specs = append(
		conf.Tables.Specs,
		Specs{
			Name: "test",
			Fk: []Fk{
				NewFk("cat_id", "category", "id"),
				NewFk("fil_id", "filter", "id"),
				NewFk("user_id", "user", "id"),
			},
			Limit: 2,
		},
	)

}
