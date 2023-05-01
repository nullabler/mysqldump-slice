package config

type Profiler struct {
	Active   bool   `yaml:"active"`
	Table    string `yaml:"table"`
	Key      string `yaml:"key"`
	Val      string `yaml:"val"`
	RefTab   string `yaml:"ref-tab"`
	RefKey   string `yaml:"ref-key"`
	RefVal   string `yaml:"ref-val"`
	TraceDep string `yaml:"trace-dep"`

	Trace bool
}
