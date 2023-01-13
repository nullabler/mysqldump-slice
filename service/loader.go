package service

import (
	"mysqldump-slice/entity"
	"mysqldump-slice/repository"
)

type Loader struct {
	conf *repository.Conf
	db   *repository.Db
	cli  *repository.Cli
	log  *Log
}

func NewLoader(conf *repository.Conf, db *repository.Db, cli *repository.Cli, log *Log) *Loader {
	return &Loader{
		conf: conf,
		db:   db,
		cli:  cli,
		log:  log,
	}
}

func (l *Loader) Relations(collect *entity.Collect) error {
	if err := l.db.LoadTables(collect); err != nil {
		return err
	}

	return l.db.LoadRelations(collect)
}

func (l *Loader) Tables(collect *entity.Collect) error {
	for _, table := range collect.Tables() {
		prKeyList, err := l.db.PrimaryKeys(table.Name)
		if err != nil {
			return err
		}

		collect.PushTab(table.Name)

		specs, ok := l.conf.Specs(table.Name)
		if len(prKeyList) == 0 {
			if !ok || len(specs.Pk) == 0 {
				continue
			}
			prKeyList = specs.Pk
		}

		if err := l.db.LoadIds(table.Name, collect, ok, specs, prKeyList, l.conf.Tables.Limit); err != nil {
			return err
		}

		l.log.Infof("- %s......Done", table.Name)
	}
	return nil
}

func (l *Loader) Weight(collect *entity.Collect) error {
	for _, table := range collect.Tables() {
		for _, rel := range collect.RelList(table.Name) {
			for _, refTab := range collect.Tables() {
				if refTab.Name == rel.RefTab() {
					refTab.Up()
				}
			}
		}

	}
	return nil
}

func (l *Loader) Dependences(collect *entity.Collect) error {
	for _, table := range collect.Tables() {
		for _, rel := range collect.RelList(table.Name) {
			keys := collect.Tab(table.Name).Keys()
			if len(keys) == 0 {
				continue
			}

			if err := l.db.LoadDeps(table.Name, collect, rel, keys); err != nil {
				return err
			}
		}

	}
	return nil
}
