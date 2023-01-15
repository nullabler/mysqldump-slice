package repository

import (
	"errors"
	"fmt"
	"mysqldump-slice/addapter"
)

type CliInterface interface {
	ExecDump(string) error
	RmFile() error
	Save() (string, error)
}

type Cli struct {
	conf *Conf
	exec addapter.ExecInterface 
}

func NewCli(conf *Conf, exec addapter.ExecInterface) (*Cli, error) {
	return &Cli{
		conf: conf,
		exec: exec,
	}, nil
}

func (c *Cli) ExecDump(call string) error {
	if len(c.conf.Tmp) == 0 {
		return errors.New("not found tmp file")
	}

	auth, err := c.auth()
	if err != nil {
		return err
	}

	return c.exec.Command(fmt.Sprintf(
		"mysqldump %s --single-transaction %s >> %s",
		auth,
		call,
		c.conf.Tmp,
	))
}

func (c *Cli) RmFile() error {
	return c.exec.Command(fmt.Sprintf("rm -f %s 2> /dev/null", c.conf.Filename()))
}

func (c *Cli) Save() (string, error) {
	filename := c.conf.Filename()
	if len(filename) == 0 {
		return "", errors.New("filename is empty")
	}
	
	if len(c.conf.Tmp) == 0 {
		return "", errors.New("tmp file is empty")
	}

	action := "cp %s %s"
	if c.conf.File.Gzip {
		action = "cat %s | gzip > %s"
	}

	return filename, c.exec.Command(fmt.Sprintf(
		action,
		c.conf.Tmp,
		filename,
	))
}

func (c *Cli) auth() (string, error) {
	if len(c.conf.DefaultExtraFile) > 0 {
		return fmt.Sprintf("--defaults-extra-file=%s", c.conf.DefaultExtraFile), nil
	}

	if len(c.conf.User) == 0 || len(c.conf.Host) == 0 {
		return "", errors.New("fail auth")
	}

	return fmt.Sprintf("-u%s -p%s -h %s",
		c.conf.User,
		c.conf.Password,
		c.conf.Host,
	), nil
}
