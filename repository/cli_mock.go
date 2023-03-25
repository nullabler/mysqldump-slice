package repository

type CliMock struct {
	str string
	err error
}

func NewCliMock() *CliMock {
	return &CliMock{}
}

func (c *CliMock) ExecDump(string) error {
	return c.err
}

func (c *CliMock) RmFile() error {
	return c.err
}

func (c *CliMock) Save() (string, error) {
	return c.str, c.err
}
