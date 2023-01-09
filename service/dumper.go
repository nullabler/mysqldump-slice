package service

import (
	"fmt"
	"mysqldump-slice/entity"
	"mysqldump-slice/repository"
	"strings"
)

type Dumper struct {
	conf *repository.Conf
	cli  *repository.Cli
	db   *repository.Db
}

func NewDumper(conf *repository.Conf, cli *repository.Cli, db *repository.Db) *Dumper {
	return &Dumper{
		conf: conf,
		cli:  cli,
		db:   db,
	}
}

func (d *Dumper) RmFile() error {
	return d.cli.RmFile()
}

func (d *Dumper) Struct() error {
	return d.cli.ExecDump(fmt.Sprintf("--no-data --routines %s", d.conf.Database))
}

func (d *Dumper) Full() error {
	return d.cli.ExecDump(fmt.Sprintf(
		"--skip-triggers --no-create-info %s %s",
		d.conf.Database,
		strings.Join(d.conf.Tables.Full, " "),
	))

}

func (d *Dumper) Slice(collect *entity.Collect) error {
	for _, table := range collect.Tables() {
		if d.hasTabNameLikeFullData(table.Name) {
			continue
		}

		keys := collect.Tab(table.Name).Keys()
		for _, where := range d.db.Where(keys) {
			if len(where) > 0 {
				err := d.cli.ExecDump(fmt.Sprintf(
					"--skip-triggers --no-create-info %s %s --where=\"%s\"",
					d.conf.Database,
					table.Name,
					where,
				))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (d *Dumper) Save() error {
	return d.cli.Save()
}

func (d *Dumper) hasTabNameLikeFullData(val string) (ok bool) {
	for i := range d.conf.Tables.Full {
		if ok = d.conf.Tables.Full[i] == val; ok {
			return
		}
	}
	return
}
