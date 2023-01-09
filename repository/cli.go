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
		"mysqldump -u%s -p%s -h %s %s >> %s",
		c.conf.User,
		c.conf.Password,
		c.conf.Host,
		call,
		c.conf.Tmp,
	))
}

func (c *Cli) RmFile() error {
	return c.exec(fmt.Sprintf("rm -f %s 2> /dev/null", c.conf.Filename()))
}

func (c *Cli) Save() error {
	action := "cp %s %s"
	if c.conf.File.Gzip {
		action = "cat %s | gzip > %s.gz"
	}

	return c.exec(fmt.Sprintf(
		action,
		c.conf.Tmp,
		c.conf.Filename(),
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
