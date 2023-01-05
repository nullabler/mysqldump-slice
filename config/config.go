package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v3"
)

type Conf struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
	Filename string `yaml:"filename"`
	Tables   Tables `yaml:"tables"`

	shell	 string
}

type Tables struct {
	Limit int `yaml:"limit"`
	Full []string `yaml:"full"`
	Ignore []string `yaml:"ignore"`
	Specs []Specs `yaml:"specs"`
}

type Specs struct {
	Name string `yaml:"name"`
	Pk []string `yaml:"pk"`
	Sort []string `yaml:"sort"`
	Limit int `yaml:"limit"`
}

func NewConf(pathToFile string) *Conf {
	conf := &Conf{
		shell: "/bin/sh",
	}

	yamlFile, err := ioutil.ReadFile(pathToFile)
    if err != nil {
		log.Printf("ReadFile: %v", err)
    }

	if err := yaml.Unmarshal(yamlFile, conf); err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return conf
}

func (conf *Conf) GetDbUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s",
		conf.User, conf.Password, conf.Host, conf.Database)
}

func (conf *Conf) Shell() string {
	return conf.shell
}

func (conf *Conf) Specs(tabName string) (Specs, bool) {
	for _, specs := range conf.Tables.Specs {
		if specs.Name == tabName {
			return specs, true
		}
	}

	return Specs{}, false
}

func (conf *Conf) Ignore(tabName string) bool {
	for _, ignore := range conf.Tables.Ignore {
		if ignore == tabName {
			return true
		}
	}

	return false
}
