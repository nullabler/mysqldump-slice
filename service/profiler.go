package service

import (
	"fmt"
	"mysqldump-slice/entity"
	"mysqldump-slice/repository"
	"strings"
)

type Profiler struct {
	conf    *repository.Conf
	resList []string
}

func NewProfiler(conf *repository.Conf) *Profiler {
	return &Profiler{
		conf: conf,
	}
}

func (p *Profiler) Conf() repository.Profiler {
	return p.conf.Profiler
}

func (p *Profiler) push(format string, params ...any) {
	p.resList = append(p.resList, fmt.Sprintf(format, params...))
}

func (p *Profiler) String() string {
	return strings.Join(p.resList, "\r\t")
}

func (p *Profiler) Relation(rel entity.RelationInterface) {
	if !p.Conf().Active {
		return
	}

	if rel.Tab() == p.Conf().Table {
		p.push("Table: %s", rel.Tab())
	}

	if rel.Col() == p.Conf().Key {
		p.push("Column: %s", rel.Col())
	}

	if rel.RefTab() == p.Conf().Table {
		p.push("Referenced Table: %s", rel.Tab())
	}

	if rel.RefCol() == p.Conf().Key {
		p.push("Referenced Column: %s", rel.Col())
	}
}
