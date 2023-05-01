package config

type Specs struct {
	Name      string   `yaml:"name"`
	Pk        []string `yaml:"pk"`
	Fk        []Fk     `yaml:"fk"`
	Sort      []string `yaml:"sort"`
	Limit     int      `yaml:"limit"`
	DepLimit  int      `yaml:"dep-limit"`
	Condition string   `yaml:"condition"`
	IsLeader  bool     `yaml:"is-leader"`
}
