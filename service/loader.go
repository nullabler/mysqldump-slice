package service

import (
	"mysqldump-slice/entity"
	"mysqldump-slice/repository"
)

type Loader struct {
	conf *repository.Conf
	db   repository.DbInterface
	cli  repository.CliInterface
	log  LogInterface
}

func NewLoader(conf *repository.Conf, db repository.DbInterface, cli repository.CliInterface, log LogInterface) *Loader {
	return &Loader{
		conf: conf,
		db:   db,
		cli:  cli,
		log:  log,
	}
}

func (l *Loader) Relations(collect entity.CollectInterface) error {
	if err := l.db.LoadTables(collect); err != nil {
		return err
	}

	if err := l.db.LoadRelations(collect); err != nil {
		return err
	}

	for _, spec := range l.conf.Tables.Specs {
		for _, fk := range spec.Fk {
			rel := entity.NewRelation()
			rel.Load(spec.Name, fk.Col, fk.FkTab, fk.FkCol)
			collect.PushRelation(rel)
		}
	}

	return nil
}

func (l *Loader) Tables(collect entity.CollectInterface) error {
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

		collect.PushPkList(table.Name, prKeyList)

		limit := l.conf.Tables.Limit
		if l.conf.IsFull(table.Name) {
			limit = 0
		}

		list, err := l.db.LoadIds(table.Name, ok, specs, prKeyList, limit)
		if err != nil {
			return err
		}

		for key, valList := range list {
			collect.PushKeyList(table.Name, key, valList)
		}

		l.log.Infof("- %s......Done", table.Name)
	}
	return nil
}

func (l *Loader) Weight(collect entity.CollectInterface) error {
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

func (l *Loader) Dependences(collect entity.CollectInterface, rel entity.RelationInterface, tabName, where string) error {
	list, err := l.db.LoadDeps(tabName, where, rel)
	if err != nil {
		return err
	}

	if !collect.IsPk(rel.RefTab(), rel.RefCol()) {
		pkList, err := l.db.LoadPkByCol(rel.RefTab(), rel.RefCol(), collect.PkList(tabName), list)
		if err != nil {
			return err
		}
		for col, list := range pkList {
			collect.PushKeyList(rel.RefTab(), col, list)
		}
	} else {
		collect.PushKeyList(rel.RefTab(), rel.RefCol(), list)
	}

	return nil
}

func (l *Loader) WhereAllByKeys(keys map[string][]string) (where string, ok bool) {
	if len(keys) == 0 {
		return
	}

	where, ok = l.db.WhereAll(keys)

	return
}
