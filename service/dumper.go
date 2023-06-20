package service

import (
	"errors"
	"fmt"
	"mysqldump-slice/config"
	"mysqldump-slice/entity"
	"mysqldump-slice/module"
	"mysqldump-slice/repository"
	"strings"
	"time"
)

type Dumper struct {
	conf *config.Conf
	cli  repository.CliInterface
	db   repository.DbInterface
	log  module.LogInterface
}

func NewDumper(conf *config.Conf, cli repository.CliInterface, db repository.DbInterface, log module.LogInterface) *Dumper {
	return &Dumper{
		conf: conf,
		cli:  cli,
		db:   db,
		log:  log,
	}
}

func (d *Dumper) RmFile() error {
	filename, err := d.Filename()
	if err != nil {
		return err
	}

	return d.cli.RmFile(filename)
}

func (d *Dumper) Struct() error {
	if err := d.cli.InitHeaderToDump(); err != nil {
		return err
	}

	return d.cli.ExecDump(fmt.Sprintf("--single-transaction --no-data --skip-comments --routines %s", d.conf.DbName()))
}

func (d *Dumper) Full() error {
	return d.cli.ExecDump(fmt.Sprintf(
		"--single-transaction --skip-add-locks --skip-quick --skip-triggers --skip-comments --no-create-info %s %s",
		d.conf.DbName(),
		strings.Join(d.conf.Tables.Full, " "),
	))

}

func (d *Dumper) Slice(collect entity.CollectInterface, tabName string, rows []*entity.Row) (err error) {
	if d.conf.IsFull(tabName) {
		return
	}

	for _, where := range d.db.Sql().WhereSlice(rows, true) {
		if len(where) > 0 {
			err = d.cli.ExecDump(fmt.Sprintf(
				`--skip-routines --quick --skip-tz-utc --single-transaction --skip-add-locks --skip-add-drop-table --skip-disable-keys --skip-set-charset --skip-triggers --skip-comments --no-create-info %s %s --where="%s"`,
				d.conf.Database,
				tabName,
				where,
			))
			if err != nil {
				return
			}
		}
	}

	return nil
}

func (d *Dumper) Save() error {
	filename, err := d.Filename()
	if err != nil {
		return err
	}

	return d.cli.Save(filename)
}

func (d *Dumper) Filename() (string, error) {
	prefix := ""
	if len(d.conf.File.Prefix) > 0 {
		prefix = d.conf.File.Prefix + "_"
	}

	date := time.Now().Format(d.conf.DateFormat())

	tail := ""
	if d.conf.File.Gzip {
		tail += ".gz"
	}

	filename := fmt.Sprintf(
		"%s%s%s_%s.sql%s",
		d.conf.File.Path,
		prefix,
		date,
		d.conf.Database,
		tail,
	)

	if len(filename) == 0 {
		return "", errors.New("filename is empty")
	}

	return filename, nil
}
