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
}

func NewConf() *Conf {
	conf := Conf{}
	flag.StringVar(&conf.user, "u", "root", "User")
	flag.StringVar(&conf.password, "p", "1234", "Password")
	flag.StringVar(&conf.host, "h", "db:3306", "Host:Port")
	flag.StringVar(&conf.database, "d", "test", "Database")
	flag.IntVar(&conf.limit, "l", 10, "Limit")
	flag.StringVar(&conf.filename, "f", "dump.sql", "Filename")
	flag.Parse()
	return &conf 
}

func (conf *Conf) GetDbUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s",
		conf.user, conf.password, conf.host, conf.database)
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
