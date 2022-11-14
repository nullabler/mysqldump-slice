package relationship

import (
	"database/sql"
)

type Table struct {
	depInt map[Relation][]int
	depStr map[Relation][]string
}

func NewTable() *Table {
	return &Table{
		depInt: make(map[Relation][]int),
		depStr: make(map[Relation][]string),
	}
}

func (tab *Table) Parse(rel Relation, isInt bool, rows *sql.Rows) (err error) {
	var id int
	var uid string

	if isInt {
		err = rows.Scan(&id)
	} else {
		err = rows.Scan(&uid)
	}

	if err == nil {
		if isInt {
			tab.depInt[rel] = append(tab.depInt[rel], id)
		} else {
			tab.depStr[rel] = append(tab.depStr[rel], uid)
		}
	}
	return
}
