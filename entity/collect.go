package entity

import "errors"

type CollectInterface interface {
	PushTable(string)
	Tables() TableList
	PushRelation(RelationInterface)
	AllRelList() map[string][]RelationInterface
	RelList(string) []RelationInterface
	PushTab(string)
	PushValList(string, []*Value) error
	Tab(string) TabInterface
	PushPkList(string, []string)
	PkList(string) []string
	IsPk(string, string) bool
	isValid(tabName string, valList []*Value) bool
}

type Collect struct {
	tables  TableList
	relList map[string][]RelationInterface
	tabList map[string]TabInterface
	pkList  map[string][]string
}

func NewCollect() *Collect {
	return &Collect{
		relList: make(map[string][]RelationInterface),
		tabList: make(map[string]TabInterface),
		pkList:  make(map[string][]string),
	}
}

func (c *Collect) PushTable(tabName string) {
	c.tables = append(c.tables, NewTable(tabName))
}

func (c *Collect) Tables() TableList {
	return c.tables
}

func (c *Collect) PushRelation(rel RelationInterface) {
	c.relList[rel.Tab()] = append(c.relList[rel.Tab()], rel)
}

func (c *Collect) AllRelList() map[string][]RelationInterface {
	return c.relList
}

func (c *Collect) RelList(tabName string) []RelationInterface {
	return c.relList[tabName]
}

func (c *Collect) PushTab(tabName string) {
	c.tabList[tabName] = NewTab(tabName)
}

func (c *Collect) PushValList(tabName string, valList []*Value) error {
	if !c.isValid(tabName, valList) {
		return errors.New("ValList is not valid")
	}

	for _, val := range valList {
		if !c.IsPk(tabName, val.key) {
			return errors.New("Key is not primary key; Where KEY: " + val.key + " TabName: " + tabName)
		}
	}

	c.Tab(tabName).Push(valList)

	return nil
}

func (c *Collect) Tab(tab string) TabInterface {
	return c.tabList[tab]
}

func (c *Collect) PushPkList(tabName string, pkList []string) {
	c.pkList[tabName] = pkList
}

func (c *Collect) PkList(tabName string) []string {
	return c.pkList[tabName]
}

func (c *Collect) isValid(tabName string, valList []*Value) bool {
	if len(c.PkList(tabName)) != len(valList) {
		return false
	}

	for _, pk := range c.PkList(tabName) {
		flag := false
		for _, val := range valList {
			if val.key == pk {
				flag = true
			}
		}

		if !flag {
			return false
		}
	}

	return true
}

func (c *Collect) IsPk(tabName, tabCol string) bool {
	for _, pk := range c.PkList(tabName) {
		if pk == tabCol {
			return true
		}
	}

	return false
}
