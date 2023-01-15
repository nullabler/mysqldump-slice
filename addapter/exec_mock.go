package addapter

type ExecMock struct {
	call string
	err  error
}

func NewExecMock() *ExecMock {
	return &ExecMock{}
}

func (e *ExecMock) Command(call string) error {
	e.call = call
	return e.err
}

func (e *ExecMock) Clear() {
	e.call = ""
	e.err = nil
}

func (e *ExecMock) Call() string {
	return e.call
}

func (e *ExecMock) SetErr(err error) {
	e.err = err
}
