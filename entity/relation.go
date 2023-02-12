package entity

import "database/sql"

type RelationInterface interface {
	Parse(*sql.Rows) error
	Load(string, string, string, string)
	Tab() string
	Col() string
	RefTab() string
	RefCol() string
}

type Relation struct {
	table     string
	column    string
	refTable  string
	refColumn string
}

func NewRelation() *Relation {
	return &Relation{}
}

func (rel *Relation) Parse(rows *sql.Rows) (err error) {
	return rows.Scan(&rel.table, &rel.refTable, &rel.column, &rel.refColumn)
}

func (rel *Relation) Load(tab, col, refTab, refCol string) {
	rel.table = tab
	rel.column = col
	rel.refTable = refTab
	rel.refColumn = refCol
}

func (rel *Relation) Tab() string {
	return rel.table
}

func (rel *Relation) Col() string {
	return rel.column
}

func (rel *Relation) RefTab() string {
	return rel.refTable
}

func (rel *Relation) RefCol() string {
	return rel.refColumn
}
