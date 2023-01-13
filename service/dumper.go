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
	log  *Log
}

func NewDumper(conf *repository.Conf, cli *repository.Cli, db *repository.Db, log *Log) *Dumper {
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
		if d.conf.IsFull(table.Name) {
			continue
		}

		keys := collect.Tab(table.Name).Keys()

		point := repository.NewPoint(d.db.Where(keys, true))
		for _, where := range d.db.WhereSlice(point) {
			if len(where) > 0 {
				err := d.cli.ExecDump(fmt.Sprintf(
					`--skip-triggers --no-create-info %s %s --where="%s"`,
					d.conf.Database,
					table.Name,
					where,
				))
				if err != nil {
					return err
				}
			}
		}
		d.log.Infof("- %s......Done", table.Name)
	}

	return nil
}

func (d *Dumper) Save() (string, error) {
	return d.cli.Save()
}
