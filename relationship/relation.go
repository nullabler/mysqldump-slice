package relationship

import "database/sql"

type Relation struct {
	foreignTable string
	primaryTable string
	fkColumn     string
	refColumn    string
}

func (rel *Relation) Parse(rows *sql.Rows) (err error) {
	return rows.Scan(&rel.foreignTable, &rel.primaryTable, &rel.fkColumn, &rel.refColumn)
}

func (rel *Relation) FrTab() string {
	return rel.foreignTable
}

func (rel *Relation) PrTab() string {
	return rel.primaryTable
}

func (rel *Relation) FkCol() string {
	return rel.fkColumn
}

func (rel *Relation) RefCol() string {
	return rel.refColumn
}
