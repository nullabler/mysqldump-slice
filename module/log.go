package module

import (
	"fmt"
	"log"
	"mysqldump-slice/config"
	"os"
)

type LogInterface interface {
	Printf(string, ...any)
	Info(...string)
	Infof(string, ...any)
	Error(error)
	Dump(data ...interface{})
	Prof(label, sql string) *Profiler
	State(string)
}

type Log struct {
	conf *config.Conf
	prof *Profiler
}

func NewLog(conf *config.Conf) *Log {
	return &Log{
		conf: conf,
		prof: NewProfiler(conf),
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

func (l *Log) Prof(label, sql string) *Profiler {
	if l.conf.Debug && !l.conf.Profiler.Active {
		log.Printf("Debug[%s]: %+v\n", label, sql)
	}

	if l.conf.Profiler.Active {
		l.prof.PushHead("++++++++++" + label + "+++++++++++++")
		l.prof.PushHead(sql)
	}

	return l.prof
}

func (l *Log) State(filename string) {
	f, err := os.Stat(filename)
	if err != nil {
		l.Error(err)
	}

	l.Printf("Save dump: %s......Done (%.2f Mb)", filename, (float64)(f.Size()/1024)/1024)
}
