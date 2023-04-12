package module

import (
	"fmt"
	"mysqldump-slice/config"
	"mysqldump-slice/entity"
	"strings"
)

type Profiler struct {
	conf     *config.Conf
	headList []string
	resList  []string
}

func NewProfiler(conf *config.Conf) *Profiler {
	return &Profiler{
		conf: conf,
	}
}

func (p *Profiler) Conf() config.Profiler {
	return p.conf.Profiler
}

func (p *Profiler) String() string {
	return strings.Join(p.resList, "\r\t")
}

func (p *Profiler) PushHead(h string) *Profiler {
	if !p.Conf().Active {
		return p
	}

	p.headList = append(p.headList, h)

	return p
}

func (p *Profiler) Table(tabName string) *Profiler {
	if !p.Conf().Active {
		return p
	}

	p.tab("Main", tabName)

	return p
}

func (p *Profiler) KeyList(list []string) *Profiler {
	if !p.Conf().Active {
		return p
	}

	for _, k := range list {
		p.col("KeyList", k)
	}

	return p
}

func (p *Profiler) ValList(list []string) *Profiler {
	if !p.Conf().Active {
		return p
	}

	for _, v := range list {
		p.col("ValList", v)
	}

	return p
}

func (p *Profiler) ValueList(list [][]*entity.Value) *Profiler {
	if !p.Conf().Active {
		return p
	}

	for _, vl := range list {
		for _, v := range vl {
			p.col("ValueList", v.Key())
			p.val("ValueList", v.Val(false))
		}
	}

	return p
}

func (p *Profiler) Relation(rel entity.RelationInterface) *Profiler {
	if !p.Conf().Active {
		return p
	}

	p.tab("Relation::Main", rel.Tab())
	p.tab("Relation::Referenced", rel.RefTab())

	p.col("Relation::Main", rel.Col())
	p.col("Relation::Referenced", rel.RefCol())

	return p
}

func (p *Profiler) clear() {
	p.headList = []string{}
}

func (p *Profiler) push(format string, params ...any) {
	if len(p.headList) > 0 {
		for _, i := range p.headList {
			p.resList = append(p.resList, i)
		}
		p.clear()
	}

	p.resList = append(p.resList, fmt.Sprintf(format, params...))
}

func (p *Profiler) tab(label, tabName string) bool {
	return p.compare("Table["+label+"]", tabName, p.Conf().Table)
}

func (p *Profiler) col(label, col string) bool {
	return p.compare("Column["+label+"]", col, p.Conf().Key)
}

func (p *Profiler) val(label, val string) bool {
	return p.compare("Value["+label+"]", val, p.Conf().Val)
}

func (p *Profiler) compare(label, got, wont string) bool {
	if got == wont {
		p.push("%s: %s", label, got)

		return true
	}

	return false
}
