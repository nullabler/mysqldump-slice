package config

type Tables struct {
	Limit  int      `yaml:"limit"`
	Full   []string `yaml:"full"`
	Ignore []string `yaml:"ignore"`
	Specs  []Specs  `yaml:"specs"`
}
