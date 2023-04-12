package module

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
func (l *LogMock) Prof(label, sql string) *Profiler { return l.p }
func (l *LogMock) State(string)                     {}
