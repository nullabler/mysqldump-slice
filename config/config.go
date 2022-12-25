package config

import (
	"flag"
	"fmt"
)

type Conf struct {
	user     string
	password string
	host     string
	database string
	limit    int
	filename string
	shell string
	tablesForFullData []string
}

func NewConf() *Conf {
	conf := Conf{
		tablesForFullData: []string{"migration_versions"},
	}
	flag.StringVar(&conf.user, "u", "root", "User")
	flag.StringVar(&conf.password, "p", "1234", "Password")
	flag.StringVar(&conf.host, "h", "mysql", "Host:Port")
	flag.StringVar(&conf.database, "d", "db", "Database")
	flag.IntVar(&conf.limit, "l", 3, "Limit")
	flag.StringVar(&conf.filename, "f", "dump.sql", "Filename")
	flag.StringVar(&conf.shell, "s", "sh", "Shell")
	flag.Parse()
	return &conf 
}

func (conf *Conf) GetDbUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s",
		conf.user, conf.password, conf.host, conf.database)
}

func (conf *Conf) User() string {
	return conf.user
}

func (conf *Conf) Passwd() string {
	return conf.password
}

func (conf *Conf) Host() string {
	return conf.host
}

func (conf *Conf) Db() string {
	return conf.database
}

func (conf *Conf) Limit() int {
	return conf.limit
}

func (conf *Conf) Filename() string {
	return conf.filename
}

func (conf *Conf) Shell() string {
	return conf.shell
}

func (conf *Conf) TablesForFullData() []string {
	return conf.tablesForFullData
}
