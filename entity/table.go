package entity

import (
	"database/sql"
	"mysqldump-slice/entity/types"
	"strings"
)

type Table[V, T types.IdInterface] struct {
	name string

	id    map[string][]V
	depId map[Relation][]T
}

func NewTable[V, T types.IdInterface](name string) *Table[V, T] {
	return &Table[V, T]{
		name: name,

		id:    make(map[string][]V),
		depId: make(map[Relation][]T),
	}
}

func (tab *Table[V, T]) Name() string {
	return tab.name
}

func (tab *Table[V, T]) Where() (query string, ok bool) {
	ok = false
	if len(tab.id) > 0 {
		query, _ = tab.WhereId()
	}

	if len(tab.depId) > 0 {
		ok = true
		if len(tab.id) > 0 {
			query += " OR "
		}
		query += tab.whereDepId()
	}
	return
}

func (tab *Table[V, T]) WhereId() (string, bool) {
	var query []string
	ok := false

	for colName, list := range tab.id {
		var idList []string
		ok = true
		for _, item := range list {
			idList = append(idList, item.String())
		}
		query = append(query, colName+" IN ("+strings.Join(idList, ", ")+")")
	}

	return strings.Join(query, " AND "), ok
}

func (tab *Table[V, T]) whereDepId() string {
	var query []string

	for rel, list := range tab.depId {
		var depFields []string
		for _, item := range list {
			depFields = append(depFields, item.String())
		}
		query = append(query, rel.refColumn+" IN ("+strings.Join(depFields, ", ")+")")
	}

	return strings.Join(query, " OR ")
}

func (tab *Table[V, T]) ParseId(colName string, rows *sql.Rows) (err error) {
	var id *V

	if err = rows.Scan(&id); err != nil {
		return
	}

	if id != nil {
		tab.id[colName] = append(tab.id[colName], *id)
	}

	return
}

func (tab *Table[V, T]) PushDep(rel Relation, rows *sql.Rows) (err error) {
	var depId *T

	if err = rows.Scan(&depId); err != nil {
		return
	}

	if depId != nil {
		tab.depId[rel] = append(tab.depId[rel], *depId)
	}

	return
}
