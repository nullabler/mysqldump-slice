package repository

func FillSimpleSpecs(conf *Conf) {
	conf.Tables.Specs = append(
		conf.Tables.Specs,
		Specs{
			Name: "test",
			Fk: []Fk{
				Fk{
					Col:   "cat_id",
					FkTab: "category",
					FkCol: "id",
				},
				Fk{
					Col:   "fil_id",
					FkTab: "filter",
					FkCol: "id",
				},
				Fk{
					Col:   "user_id",
					FkTab: "user",
					FkCol: "id",
				},
			},
			Limit: 2,
		},
	)

}
