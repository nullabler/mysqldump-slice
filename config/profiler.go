package config

type Profiler struct {
	Active bool   `yaml:"active"`
	Table  string `yaml:"table"`
	Key    string `yaml:"key"`
	Val    string `yaml:"val"`
}
