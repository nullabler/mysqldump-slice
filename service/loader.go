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
			rel.Load(spec.Name, fk.Col, fk.FkTab, fk.FkCol, fk.Limit, fk.IsGreedy)
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
		l.profForTableValList(table.Name, list)

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

func (l *Loader) Dependences(
	collect entity.CollectInterface,
	rel entity.RelationInterface,
	tabName string,
	rows []*entity.Row,
) error {
	isProf := l.profRel(tabName, rel)

	if l.conf.IsFull(rel.RefTab()) {
		if isProf && rel.RefTab() == l.conf.Profiler.RefTab {
			l.log.Printf("Skip load Dependence table: %s because was full load", rel.RefTab())
		}
		return nil
	}

	for _, where := range l.db.Sql().WhereSlice(rows, false) {
		list, err := l.db.LoadDeps(tabName, where, rel)
		if isProf {
			if l.log.ProfStrList(list, true, false) {
				l.log.Printf("Find in StrList Dependence table: %s Val: %s Where: %s", tabName, l.conf.Profiler.Val, where)
			}

			if l.log.ProfStrList(list, false, true) {
				l.log.Printf("Find in StrList Dependence table: %s RefVal: %s Where: %s", tabName, l.conf.Profiler.RefVal, where)
			}
		}

		if err != nil {
			return err
		}

		if len(list) == 0 {
			continue
		}

		if collect.IsPk(rel.RefTab(), rel.RefCol()) && len(collect.PkList(rel.RefTab())) == 1 {
			valList := []entity.ValList{}
			for _, v := range list {
				t := []*entity.Value{
					entity.NewValue(rel.RefCol(), v),
				}

				valList = append(valList, t)
			}

			collect.PushValList(rel.RefTab(), valList)

			l.log.ProfTraceDep(tabName, rel)
			if isProf {
				l.profForDependenceValList("IsPk:TRUE", tabName, where, rel, valList)
			}

			continue
		}

		valList, err := l.db.LoadPkByCol(
			rel.RefTab(),
			rel.RefCol(),
			collect.PkList(rel.RefTab()),
			list,
			rel.IsGreedy(),
		)

		if err != nil {
			return err
		}

		collect.PushValList(rel.RefTab(), valList)

		l.log.ProfTraceDep(tabName, rel)
		if isProf {
			l.profForDependenceValList("IsPk:FALSE", tabName, where, rel, valList)
		}
	}

	return nil
}

func (l *Loader) profForDependenceValList(label, tabName, where string, rel entity.RelationInterface, valList []entity.ValList) {
	if checkList := l.log.ProfValList(rel.RefTab(), valList, true, false); len(checkList) > 0 {
		l.log.Printf("Find Dependence<%s> Table: %s RefTab: %s Val: %s Where: %s", label, tabName, rel.RefTab(), l.conf.Profiler.Val, where)
		for _, v := range checkList {
			l.log.Printf("+ Key: %s Val: %s", v.Key(), v.Val(false))
		}
	}

	if checkList := l.log.ProfValList(rel.RefTab(), valList, false, true); len(checkList) > 0 {
		l.log.Printf("Find Dependence<%s> Table: %s RefTab: %s RefVal: %s Where %s", label, tabName, rel.RefTab(), l.conf.Profiler.RefVal, where)
		for _, v := range checkList {
			l.log.Printf("+ Key: %s Val: %s", v.Key(), v.Val(false))
		}
	}
}

func (l *Loader) profForTableValList(tabName string, valList []entity.ValList) {
	if checkList := l.log.ProfValList(tabName, valList, true, false); len(checkList) > 0 {
		l.log.Printf("Find ValList for Table: %s Val: %s", tabName, l.conf.Profiler.Val)
		for _, v := range checkList {
			l.log.Printf("+ Key: %s Val: %s", v.Key(), v.Val(false))
		}
	}

	if checkList := l.log.ProfValList(tabName, valList, false, true); len(checkList) > 0 {
		l.log.Printf("Find ValList for Table: %s RefVal: %s", tabName, l.conf.Profiler.RefVal)
		for _, v := range checkList {
			l.log.Printf("+ Key: %s Val: %s", v.Key(), v.Val(false))
		}
	}
}

func (l *Loader) profRel(tabName string, rel entity.RelationInterface) bool {
	isTab := false
	if l.log.ProfRel(tabName, rel, true, false) {
		isTab = true
		l.log.Printf("Find Table: %s RelTab: %s", tabName, rel.RefTab())
	}

	isCol := false
	if l.log.ProfRel(tabName, rel, false, true) {
		isCol = true
		l.log.Printf("Find Table: %s RelCol: %s", tabName, rel.RefCol())
	}

	return isTab || isCol
}
