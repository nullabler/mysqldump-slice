package entity

import "database/sql"

type RelationInterface interface {
	Parse(*sql.Rows) error
	Load(string, string, string, string, int, bool)
	Tab() string
	Col() string
	RefTab() string
	RefCol() string
	Limit() int
	IsGreedy() bool
}

type Relation struct {
	table     string
	column    string
	refTable  string
	refColumn string
	limit     int
	isGreedy  bool
}

func NewRelation() *Relation {
	return &Relation{
		limit: 0,
	}
}

func (rel *Relation) Parse(rows *sql.Rows) (err error) {
	return rows.Scan(&rel.table, &rel.refTable, &rel.column, &rel.refColumn)
}

func (rel *Relation) Load(tab, col, refTab, refCol string, limit int, isGreedy bool) {
	rel.table = tab
	rel.column = col
	rel.refTable = refTab
	rel.refColumn = refCol
	rel.limit = limit
	rel.isGreedy = isGreedy
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

func (rel *Relation) Limit() int {
	return rel.limit
}

func (rel *Relation) IsGreedy() bool {
	return rel.isGreedy
}
