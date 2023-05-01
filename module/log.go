package module

import (
	"fmt"
	"log"
	"mysqldump-slice/config"
	"mysqldump-slice/entity"
	"os"
)

type LogInterface interface {
	Printf(string, ...any)
	Info(...string)
	Infof(string, ...any)
	Error(error)
	Dump(data ...interface{})
	Debug(label, sql string)
	State(string)
	ProfRowList(tabName string, rows []*entity.Row, enableUsed, enableKey, enableVal bool) bool
	ProfRel(tabName string, rel entity.RelationInterface, enableRefTab, enableRefCol bool) bool
	ProfStrList(list []string, enableVal, enableRefVal bool) bool
	ProfValList(tabName string, valList []entity.ValList, enableVal, enableRefVal bool) entity.ValList
	ProfTraceDep(tabName string, rel entity.RelationInterface)
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

func (l *Log) Debug(label, sql string) {
	if l.conf.Debug {
		log.Printf("Debug[%s]: %+v\n", label, sql)
	}
}

func (l *Log) State(filename string) {
	f, err := os.Stat(filename)
	if err != nil {
		l.Error(err)
	}

	l.Printf("Save dump: %s......Done (%.2f Mb)", filename, (float64)(f.Size()/1024)/1024)
}

func (l *Log) ProfRowList(tabName string, rows []*entity.Row, enableUsed, enableKey, enableVal bool) bool {
	if !l.prof.Active() || !l.prof.TabName(tabName) {
		return false
	}

	return l.prof.RowList(rows, enableUsed, enableKey, enableVal)
}

func (l *Log) ProfRel(tabName string, rel entity.RelationInterface, enableRefTab, enableRefCol bool) bool {
	if !l.prof.Active() || !l.prof.TabName(tabName) {
		return false
	}

	return l.prof.Rel(rel, enableRefTab, enableRefCol)
}

func (l *Log) ProfStrList(list []string, enableVal, enableRefVal bool) bool {
	if !l.prof.Active() {
		return false
	}

	return l.prof.StrList(list, enableVal, enableRefVal)
}

func (l *Log) ProfValList(tabName string, valList []entity.ValList, enableVal, enableRefVal bool) entity.ValList {
	if !l.prof.Active() || !l.prof.TabName(tabName) {
		return nil
	}

	return l.prof.ValList(valList, enableVal, enableRefVal)
}

func (l *Log) ProfTraceDep(tabName string, rel entity.RelationInterface) {
	if !l.prof.Active() || !l.prof.TraceDep(tabName, rel.RefTab()) {
		return
	}

	l.Printf("Push dependence Table: %s RefTab: %s, RefCol: %s", tabName, rel.RefTab(), rel.RefCol())
}
