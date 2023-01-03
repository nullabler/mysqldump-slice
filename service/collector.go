package service

import (
	"database/sql"
	"fmt"
	"mysqldump-slice/entity"
	"mysqldump-slice/entity/types"
	"strings"
)

type Collector struct {
	relationList map[string][]entity.Relation

	tabIntInt map[string]*entity.Table[types.IdType, types.IdType]
	tabIntStr map[string]*entity.Table[types.IdType, types.UidType]
	tabStrInt map[string]*entity.Table[types.UidType, types.IdType]
	tabStrStr map[string]*entity.Table[types.UidType, types.UidType]
}


func NewCollector() *Collector {
	return &Collector {
		relationList: make(map[string][]entity.Relation),

		tabIntInt: make(map[string]*entity.Table[types.IdType, types.IdType]),
		tabIntStr: make(map[string]*entity.Table[types.IdType, types.UidType]),
		tabStrInt: make(map[string]*entity.Table[types.UidType, types.IdType]),
		tabStrStr: make(map[string]*entity.Table[types.UidType, types.UidType]),
	}
}

func (ctl *Collector) PushRelation(relation entity.Relation) {
	ctl.relationList[relation.Tab()] = append(ctl.relationList[relation.Tab()], relation)
}

func (ctl *Collector) RelationList(tabName string) []entity.Relation {
	return ctl.relationList[tabName]
}

func (ctl *Collector) PushTable(tabName string) {
		ctl.tabIntInt[tabName] = entity.NewTable[types.IdType, types.IdType](tabName)
		ctl.tabIntStr[tabName] = entity.NewTable[types.IdType, types.UidType](tabName)
		ctl.tabStrInt[tabName] = entity.NewTable[types.UidType, types.IdType](tabName)
		ctl.tabStrStr[tabName] = entity.NewTable[types.UidType, types.UidType](tabName)
}

func (ctl *Collector) Where(tabName string) string {
	var list []string

	list = append(list, ctl.tabIntInt[tabName].Where())
	list = append(list, ctl.tabIntStr[tabName].Where())
	list = append(list, ctl.tabStrInt[tabName].Where())
	list = append(list, ctl.tabStrStr[tabName].Where())

	return strings.Join(list, " OR ")
}

func (ctl *Collector) WhereId(tabName string) string {
	var list []string

	list = append(list, ctl.tabIntInt[tabName].WhereId())
	list = append(list, ctl.tabStrStr[tabName].WhereId())

	return strings.Join(list, " AND ")
}

func (ctl *Collector) ParseId(tabName string, colName string, isInt bool, rows *sql.Rows) {
	if isInt {
		ctl.tabIntInt[tabName].ParseId(colName, rows)
		ctl.tabIntStr[tabName].ParseId(colName, rows)
	} else {
		ctl.tabStrInt[tabName].ParseId(colName, rows)
		ctl.tabStrStr[tabName].ParseId(colName, rows)
	}
}

func (ctl *Collector) PushDep(tabName string, isIntDep bool, rel entity.Relation, rows *sql.Rows) {
	fmt.Println("DEBUGING PushDep:", tabName, isIntDep, ctl.tabIntInt[tabName], ctl.tabStrInt[tabName])
	if ctl.tabIntInt[tabName] != nil {
		if isIntDep {
			ctl.tabIntInt[tabName].PushDep(rel, rows)
		} else {
			ctl.tabIntStr[tabName].PushDep(rel, rows)
		}
	} else {
		if isIntDep {
			ctl.tabStrInt[tabName].PushDep(rel, rows)
		} else {
			ctl.tabStrStr[tabName].PushDep(rel, rows)
		}
	}
}
