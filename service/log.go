package service

import (
	"fmt"
	"log"
	"mysqldump-slice/repository"
	"os"
)

type LogInterface interface {
	Printf(string, ...any)
	Info(...string)
	Infof(string, ...any)
	Error(error)
	Dump(...interface{})
	State(string)
}

type Log struct {
	conf *repository.Conf
}

func NewLog(conf *repository.Conf) *Log {
	return &Log{
		conf: conf,
	}
}

func (l *Log) Printf(format string, params ...any) {
	log.Printf(format, params...)
}

func (l *Log) Info(msgs ...string) {
	if l.conf.Log {
		log.Println(msgs)
	}
}

func (l *Log) Infof(format string, params ...any) {
	if l.conf.Log {
		l.Info(fmt.Sprintf(format, params...))
	}
}

func (l *Log) Error(err error) {
	if l.conf.Log {
		log.Panic(err)
	}
	panic(err.Error())
}

func (l *Log) Dump(data ...interface{}) {
	log.Printf("%+v\n", data)
}

func (l *Log) Debug(data ...interface{}) {
	if l.conf.Debug {
		l.Dump(data)
	}
}

func (l *Log) State(filename string) {
	f, err := os.Stat(filename)
	if err != nil {
		l.Error(err)
	}

	l.Printf("Save dump: %s......Done (%.2f Mb)", filename, (float64)(f.Size()/1024)/1024)
}
