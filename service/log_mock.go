package service

type LogMock struct {
}

func NewLogMock() *LogMock {
	return &LogMock{}
}

func (l *LogMock) Printf(string, ...any) {}
func (l *LogMock) Info(...string)        {}
func (l *LogMock) Infof(string, ...any)  {}
func (l *LogMock) Error(error)           {}
func (l *LogMock) Dump(...interface{})   {}
func (l *LogMock) State(string)          {}
