package entity

import (
	"database/sql"
	"fmt"
)

type Collect struct {
	tables  []string
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
	c.tables = append(c.tables, tabName)
}

func (c *Collect) Tables() []string {
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

func (c *Collect) PushKey(tab, col string, isInt bool, rows *sql.Rows) {
	var id int
	var uid string
	if isInt {
		if err := rows.Scan(&id); err != nil {
			return
		}
		uid = fmt.Sprint(id)
	} else {
		if err := rows.Scan(&uid); err != nil {
			return
		}
		uid = fmt.Sprintf("'%s'", uid)
	}

	c.tabList[tab].Push(col, uid)
}

func (c *Collect) Tab(tab string) *Tab {
	return c.tabList[tab]
}
