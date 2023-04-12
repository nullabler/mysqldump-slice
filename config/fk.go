package config

type Fk struct {
	Col   string `yaml:"col"`
	FkTab string `yaml:"fk_tab"`
	FkCol string `yaml:"fk_col"`
	Limit int    `yaml:"limit"`
}

func NewFk(col, fkTab, fkCol string) Fk {
	return Fk{
		Col:   col,
		FkTab: fkTab,
		FkCol: fkCol,
	}
}
