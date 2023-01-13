package entity

import (
	"database/sql"
	"fmt"
)

type Collect struct {
	tables  TableList
	relList map[string][]Relation
	tabList map[string]*Tab
}

func NewCollect() *Collect {
	return &Collect{
		relList: make(map[string][]Relation),
		tabList: make(map[string]*Tab),
	}
}

func (c *Collect) PushTable(tabName string) {
	c.tables = append(c.tables, NewTable(tabName))
}

func (c *Collect) Tables() TableList {
	return c.tables
}

func (c *Collect) PushRelation(rel Relation) {
	c.relList[rel.Tab()] = append(c.relList[rel.Tab()], rel)
}

func (c *Collect) RelList(tabName string) []Relation {
	return c.relList[tabName]
}

func (c *Collect) PushTab(tabName string) {
	c.tabList[tabName] = NewTab(tabName)
}

func (c *Collect) PushKey(tab, col string, isInt bool, rows *sql.Rows) error {
	var id *int
	var uid *string
	var key string

	if isInt {
		if err := rows.Scan(&id); err != nil {
			return err
		}

		if id != nil {
			key = fmt.Sprint(*id)
		}
	} else {
		if err := rows.Scan(&uid); err != nil {
			return err
		}

		if uid != nil {
			key = fmt.Sprintf("'%s'", *uid)
		}
	}

	if len(key) > 0 {
		c.tabList[tab].Push(col, key)
	}

	return nil
}

func (c *Collect) Tab(tab string) *Tab {
	return c.tabList[tab]
}
