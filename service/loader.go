package service

import (
	"mysqldump-slice/config"
	"mysqldump-slice/entity"
	"mysqldump-slice/module"
	"mysqldump-slice/repository"
)

type Loader struct {
	conf *config.Conf
	db   repository.DbInterface
	cli  repository.CliInterface
	log  module.LogInterface
}

func NewLoader(conf *config.Conf, db repository.DbInterface, cli repository.CliInterface, log module.LogInterface) *Loader {
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
			rel.Load(spec.Name, fk.Col, fk.FkTab, fk.FkCol, spec.Limit)
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

		if l.conf.IsIgnore(table.Name) {
			continue
		}

		list, err := l.db.LoadIds(table.Name, &specs, prKeyList)
		if err != nil {
			return err
		}

		collect.PushValList(table.Name, list)

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

func (l *Loader) Dependences(collect entity.CollectInterface, rel entity.RelationInterface, tabName string, rows []*entity.Row) error {
	for _, where := range l.db.Sql().WhereSlice(rows, false) {
		list, err := l.db.LoadDeps(tabName, where, rel)

		if err != nil {
			return err
		}

		if len(list) == 0 {
			continue
		}

		if collect.IsPk(rel.RefTab(), rel.RefCol()) && len(collect.PkList(rel.RefTab())) == 1 {
			valList := [][]*entity.Value{}
			for _, v := range list {
				t := []*entity.Value{
					entity.NewValue(rel.RefCol(), v),
				}

				valList = append(valList, t)
			}

			collect.PushValList(rel.RefTab(), valList)

			continue
		}

		valList, err := l.db.LoadPkByCol(rel.RefTab(), rel.RefCol(), collect.PkList(rel.RefTab()), list)
		if err != nil {
			return err
		}

		collect.PushValList(rel.RefTab(), valList)
	}

	return nil
}
