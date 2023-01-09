package service

import (
	"mysqldump-slice/entity"
	"mysqldump-slice/repository"
)

type Loader struct {
	conf *repository.Conf
	db   *repository.Db
	cli  *repository.Cli
}

func NewLoader(conf *repository.Conf, db *repository.Db, cli *repository.Cli) *Loader {
	return &Loader{
		conf: conf,
		db:   db,
		cli:  cli,
	}
}

func (l *Loader) Relations(collect *entity.Collect) error {
	l.db.LoadTables(l.conf, collect)
	return l.db.LoadRelations(collect)
}

func (l *Loader) Tables(collect *entity.Collect) error {
	for _, tabName := range collect.Tables() {
		prKeyList := l.db.PrimaryKeys(tabName)

		collect.PushTab(tabName)

		specs, ok := l.conf.Specs(tabName)
		if len(prKeyList) == 0 {
			if !ok || len(specs.Pk) == 0 {
				continue
			}
			prKeyList = specs.Pk
		}

		l.db.LoadIds(tabName, collect, ok, specs, prKeyList, l.conf.Tables.Limit)
	}
	return nil
}

func (l *Loader) Dependences(collect *entity.Collect) error {
	for _, tabName := range collect.Tables() {
		for _, rel := range collect.RelList(tabName) {
			keys := collect.Tab(tabName).Keys()
			if len(keys) > 0 {
				continue
			}

			l.db.LoadDeps(tabName, collect, rel, keys)
		}

	}
	return nil
}
