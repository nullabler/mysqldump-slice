package repository

import (
	"fmt"
	"os"
	"os/exec"
)

type Cli struct {
	conf *Conf
}

func NewCli(conf *Conf) (*Cli, error) {
	return &Cli{
		conf: conf,
	}, nil
}

func (c *Cli) ExecDump(call string) error {
	return c.exec(fmt.Sprintf(
		"mysqldump %s --single-transaction %s >> %s",
		c.auth(),
		call,
		c.conf.Tmp,
	))
}

func (c *Cli) RmFile() error {
	return c.exec(fmt.Sprintf("rm -f %s 2> /dev/null", c.conf.Filename()))
}

func (c *Cli) Save() (string, error) {
	filename := c.conf.Filename()
	action := "cp %s %s"
	if c.conf.File.Gzip {
		filename += ".gz"
		action = "cat %s | gzip > %s"
	}

	return filename, c.exec(fmt.Sprintf(
		action,
		c.conf.Tmp,
		filename,
	))
}

func (c *Cli) exec(call string) error {
	cmd := exec.Command(c.conf.Shell(), "-c", call)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (c *Cli) auth() string {
	if len(c.conf.DefaultExtraFile) > 0 {
		return fmt.Sprintf("--defaults-extra-file=%s", c.conf.DefaultExtraFile)
	}

	return fmt.Sprintf("-u%s -p%s -h %s",
		c.conf.User,
		c.conf.Password,
		c.conf.Host,
	)
}
