package service

import (
	"fmt"
	"mysqldump-slice/entity"
	"mysqldump-slice/repository"
	"strings"
)

type Dumper struct {
	conf *repository.Conf
	cli  repository.CliInterface
	db   repository.DbInterface
	log  LogInterface
}

func NewDumper(conf *repository.Conf, cli repository.CliInterface, db repository.DbInterface, log LogInterface) *Dumper {
	return &Dumper{
		conf: conf,
		cli:  cli,
		db:   db,
		log:  log,
	}
}

func (d *Dumper) RmFile() error {
	return d.cli.RmFile()
}

func (d *Dumper) Struct() error {
	return d.cli.ExecDump(fmt.Sprintf("--no-data --skip-comments --routines %s", d.conf.DbName()))
}

func (d *Dumper) Full() error {
	return d.cli.ExecDump(fmt.Sprintf(
		"--skip-triggers --no-create-info %s %s",
		d.conf.DbName(),
		strings.Join(d.conf.Tables.Full, " "),
	))

}

func (d *Dumper) Slice(collect entity.CollectInterface, tabName string, keys map[string][]string) (err error) {
	if d.conf.IsFull(tabName) {
		return
	}

	point := repository.NewPoint(d.db.Where(keys, true))
	for _, where := range d.db.WhereSlice(point) {
		if len(where) > 0 {
			err = d.cli.ExecDump(fmt.Sprintf(
				`--skip-triggers --no-create-info %s %s --where="%s"`,
				d.conf.Database,
				tabName,
				where,
			))
			if err != nil {
				return
			}
		}
	}
	d.log.Infof("- %s......Done", tabName)

	return nil
}

func (d *Dumper) Save() (string, error) {
	return d.cli.Save()
}
