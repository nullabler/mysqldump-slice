package module

import "mysqldump-slice/entity"

type LogMock struct {
	p *Profiler
}

func NewLogMock() *LogMock {
	return &LogMock{
		p: &Profiler{},
	}
}

func (l *LogMock) Printf(string, ...any)            {}
func (l *LogMock) Info(...string)                   {}
func (l *LogMock) Infof(string, ...any)             {}
func (l *LogMock) Error(error)                      {}
func (l *LogMock) Dump(data ...interface{})         {}
func (l *LogMock) Debug(label, sql string)          {}
func (l *LogMock) Prof(label, sql string) *Profiler { return l.p }
func (l *LogMock) State(string)                     {}
func (l *LogMock) ProfRowList(tabName string, rows []*entity.Row, enableUsed, enableKey, enableVal bool) bool {
	return true
}
func (l *LogMock) ProfRel(tabName string, rel entity.RelationInterface, enableRefTab, enableRefCol bool) bool {
	return true
}
func (l *LogMock) ProfStrList(list []string, enableVal, enableRefVal bool) bool { return true }
func (l *LogMock) ProfValList(tabName string, valList []entity.ValList, enableVal, enableRefVal bool) entity.ValList {
	return entity.ValList{}
}
func (l *LogMock) ProfTraceDep(tabName string, rel entity.RelationInterface) {}
