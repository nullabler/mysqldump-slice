package repository

type CliMock struct {
	err error
}

func NewCliMock() *CliMock {
	return &CliMock{}
}

func (c *CliMock) ExecDump(string) error {
	return c.err
}

func (c *CliMock) RmFile(string) error {
	return c.err
}

func (c *CliMock) Save(string) error {
	return c.err
}
