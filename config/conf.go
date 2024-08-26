package config

import (
	"fmt"
	"io/ioutil"
	"mysqldump-slice/helper"

	"gopkg.in/yaml.v3"
)

type Conf struct {
	ConfVersion      string `yaml:"conf-version"`
	User             string `yaml:"user"`
	Password         string `yaml:"password"`
	Host             string `yaml:"host"`
	Port             string `yaml:"port"`
	Database         string `yaml:"database"`
	DefaultExtraFile string `yaml:"default-extra-file"`

	MaxConnectCount          int      `yaml:"max-connect"`
	MaxLifetimeConnectMinute int      `yaml:"max-lifetime-connect-minute"`
	MaxLifetimeQuerySecond   int      `yaml:"max-lifetime-query-second"`
	Log                      bool     `yaml:"log"`
	Debug                    bool     `yaml:"debug"`
	Profiler                 Profiler `yaml:"profiler"`

	File   File   `yaml:"filename"`
	Tables Tables `yaml:"tables"`

	version string
	shell   string
	def     Default

	Tmp      string
	LimitCli int
}

func NewConf(version, pathToFile, tmpFilename string) (*Conf, error) {
	conf := &Conf{
		version:  version,
		shell:    "/bin/sh",
		Debug:    false,
		Tmp:      tmpFilename,
		LimitCli: 30,
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

	return conf, nil
}

func (conf *Conf) Version() string {
	return conf.version
}

func (conf *Conf) DbName() string {
	return conf.Database
}

func (conf *Conf) DbUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
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
	return helper.SliceIsExist(tabName, conf.Tables.Ignore)
}

func (conf *Conf) IsFull(tabName string) bool {
	return helper.SliceIsExist(tabName, conf.Tables.Full)
}

func (conf *Conf) DateFormat() string {
	format := conf.def.dateFormat
	if len(conf.File.DateFormat) > 0 {
		format = conf.File.DateFormat
	}

	return format
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
