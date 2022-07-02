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
}

func NewConf() *Conf {
	conf := Conf{}
	flag.StringVar(&conf.user, "u", "root", "User")
	flag.StringVar(&conf.password, "p", "123", "Password")
	flag.StringVar(&conf.host, "h", "localhost:3306", "Host:Port")
	flag.StringVar(&conf.database, "d", "test", "Database")
	flag.IntVar(&conf.limit, "l", 100, "Limit")
	flag.Parse()
	return &conf 
}

func (conf *Conf) GetDbUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s",
		conf.user, conf.password, conf.host, conf.database)
}

func (conf *Conf) Db() string {
	return conf.database
}

func (conf *Conf) Limit() int {
	return conf.limit
}
