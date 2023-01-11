package repository

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Conf struct {
	User                     string `yaml:"user"`
	Password                 string `yaml:"password"`
	Host                     string `yaml:"host"`
	Database                 string `yaml:"database"`
	MaxConnectCount          int    `yaml:"max-connect"`
	MaxLifetimeConnectMinute int    `yaml:"max-lifetime-connect-minute"`
	MaxLifetimeQuerySecond   int    `yaml:"max-lifetime-query-second"`
	Log                      bool   `yaml:"log"`

	File   File   `yaml:"filename"`
	Tables Tables `yaml:"tables"`

	shell string
	def   Default

	Tmp      string
	LimitCli int
}

type Default struct {
	dateFormat               string
	maxConnectCount          int
	maxLifetimeConnectMinute int
	maxLifetimeQuerySecond   int
}

type File struct {
	Path       string `yaml:"path"`
	Prefix     string `yaml:"prefix"`
	DateFormat string `yaml:"date-format"`
	Gzip       bool   `yaml:"gzip"`
}

type Tables struct {
	Limit  int      `yaml:"limit"`
	Full   []string `yaml:"full"`
	Ignore []string `yaml:"ignore"`
	Specs  []Specs  `yaml:"specs"`
}

type Specs struct {
	Name      string   `yaml:"name"`
	Pk        []string `yaml:"pk"`
	Sort      []string `yaml:"sort"`
	Limit     int      `yaml:"limit"`
	Condition string   `yaml:"condition"`
}

func NewConf(pathToFile, tmpFilename string) (*Conf, error) {
	conf := &Conf{
		shell:    "/bin/sh",
		Tmp:      tmpFilename,
		LimitCli: 10,
		def: Default{
			dateFormat:               "2006-01-02_15_04",
			maxConnectCount:          10,
			maxLifetimeConnectMinute: 5,
			maxLifetimeQuerySecond:   3,
		},
	}

	yamlFile, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(yamlFile, conf); err != nil {
		return nil, err
	}
	defer os.Remove(conf.Tmp)

	return conf, nil
}

func (conf *Conf) DbUrl() string {
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

func (conf *Conf) IsIgnore(tabName string) bool {
	return conf.hasTable(tabName, conf.Tables.Ignore)
}

func (conf *Conf) IsFull(tabName string) bool {
	return conf.hasTable(tabName, conf.Tables.Full)
}

func (conf *Conf) Filename() string {
	prefix := ""
	if len(conf.File.Prefix) > 0 {
		prefix = conf.File.Prefix + "_"
	}

	format := conf.def.dateFormat
	if len(conf.File.DateFormat) > 0 {
		format = conf.File.DateFormat
	}
	date := time.Now().Format(format)

	return fmt.Sprintf(
		"%s%s%s_%s.sql",
		conf.File.Path,
		prefix,
		date,
		conf.Database,
	)
}

func (conf *Conf) MaxConnect() int {
	if conf.MaxConnectCount > 0 {
		return conf.MaxConnectCount
	}

	return conf.def.maxConnectCount
}

func (conf *Conf) MaxLifetimeConnect() int {
	if conf.MaxLifetimeConnectMinute > 0 {
		return conf.MaxLifetimeConnectMinute
	}

	return conf.def.maxLifetimeConnectMinute
}

func (conf *Conf) MaxLifetimeQuery() int {
	if conf.MaxLifetimeQuerySecond > 0 {
		return conf.MaxLifetimeQuerySecond
	}

	return conf.def.maxLifetimeQuerySecond
}

func (conf *Conf) hasTable(tabName string, tabList []string) bool {
	for _, tab := range tabList {
		if tab == tabName {
			return true
		}
	}

	return false
}
