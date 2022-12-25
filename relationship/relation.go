package relationship

import "database/sql"

type Relation struct {
	table string
	column     string
	refTable string
	refColumn    string
}

func (rel *Relation) Parse(rows *sql.Rows) (err error) {
	return rows.Scan(&rel.table, &rel.refTable, &rel.column, &rel.refColumn)
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
