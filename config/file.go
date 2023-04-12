package config

type File struct {
	Path       string `yaml:"path"`
	Prefix     string `yaml:"prefix"`
	DateFormat string `yaml:"date-format"`
	Gzip       bool   `yaml:"gzip"`
}
