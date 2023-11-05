package repository

import (
	"errors"
	"fmt"
	"mysqldump-slice/addapter"
	"mysqldump-slice/config"
	"time"
)

type CliInterface interface {
	InitHeaderToDump() error
	ExecDump(call string) error
	RmFile(filename string) error
	Save(filename string) error
}

type Cli struct {
	conf *config.Conf
	exec addapter.ExecInterface
}

func NewCli(conf *config.Conf, exec addapter.ExecInterface) (*Cli, error) {
	return &Cli{
		conf: conf,
		exec: exec,
	}, nil
}

func (c *Cli) InitHeaderToDump() error {
	return c.exec.Command(fmt.Sprintf(
		"echo '-- Slicer version: %s \n-- DateTime: %s\n' >> %s",
		c.conf.Version(),
		time.Now().Format("2006.01.02 15:04:05"),
		c.conf.Tmp,
	))
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

func (c *Cli) RmFile(filename string) error {
	if err := c.isValidFilename(filename); err != nil {
		return err
	}

	return c.exec.Command(fmt.Sprintf("rm -f %s 2> /dev/null", filename))
}

func (c *Cli) Save(filename string) error {
	if err := c.isValidFilename(filename); err != nil {
		return err
	}

	if len(c.conf.Tmp) == 0 {
		return errors.New("not found tmp file")
	}

	action := "cp %s %s"
	if c.conf.File.Gzip {
		action = "cat %s | gzip > %s"
	}

	return c.exec.Command(fmt.Sprintf(
		action,
		c.conf.Tmp,
		filename,
	))
}

func (c *Cli) isValidFilename(filename string) error {
	if len(filename) == 0 {
		return errors.New("Filename is empty")
	}

	return nil
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
