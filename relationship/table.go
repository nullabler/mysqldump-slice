package relationship

import (
	"database/sql"
	"strconv"
	"strings"
)

type IdInterface interface {
	string | int
	String() string
}

type IdType int

func (id IdType) String() string {
	return strconv.Itoa(int(id))
}

type UidType string

func (str UidType) String() string {
	return string(str)
}

type Table[V IdInterface] struct {
	depInt map[Relation][]*int
	depStr map[Relation][]*string
	queryList []string
	groupByList []string

	name string
	relationList []Relation

	id []V
	uid []*string
	isIntId *bool
	refInt map[string][]*int
	refStr map[string][]*string
}

func NewTable[V IdInterface](name string) *Table[V] {
	return &Table[V]{
		depInt: make(map[Relation][]*int),
		depStr: make(map[Relation][]*string),

		name: name,
		refInt: make(map[string][]*int),
		refStr: make(map[string][]*string),
	}
}

func (tab *Table[V]) Name() string {
	return tab.name
}

func (tab *Table[V]) SetIsIntId(isIntId bool) {
	tab.isIntId = &isIntId
}

func (tab *Table[V]) IsIntId() *bool {
	return tab.isIntId
}

//func (tab *Table[V]) WhereId() (string, bool) {
	//if tab.isIntId == nil {
		//return "", false
	//}

	//if *tab.isIntId {
		//return tab.whereIdLikeInt(), true 
	//}

	//return tab.whereIdLikeStr(), true 
//}


func (tab *Table[V]) WhereId(list []V) string {
	var res []string

	for _, item := range list {
		res = append(res, item.String())
	}

	return strings.Join(res, ", ")
}

//func (tab *Table[V]) whereIdLikeStr() string {
	//var list []string

	//for _, uid := range tab.id {
		//list = append(list, "'" + *uid + "'")
	//}

	//return strings.Join(list, ", ")
//}

func (tab *Table[V]) ParseOld(rel Relation, isInt bool, rows *sql.Rows) (err error) {
	var id *int
	var uid *string

	if isInt {
		err = rows.Scan(&id)
	} else {
		err = rows.Scan(&uid)
	}

	if err == nil {
		if isInt {
			if id != nil {
				tab.depInt[rel] = append(tab.depInt[rel], id)
			}
		} else {
			tab.depStr[rel] = append(tab.depStr[rel], uid)
		}
	}
	return
}

func (tab *Table[V]) ParseId(rows *sql.Rows) (err error) {
	var id *V

	if err = rows.Scan(&id); err != nil {
		return
	}

	tab.id = append(tab.id, *id)

	return
}

func (tab *Table[V]) ParseRef(isInt bool, refCol string, rows *sql.Rows) (err error) {
	var refId *int
	var refUid *string

	if *tab.isIntId {
		err = rows.Scan(&refId)
	} else {
		err = rows.Scan(&refUid)
	}

	if err == nil {
		if isInt {
			if refId != nil {
				tab.refInt[refCol] = append(tab.refInt[refCol], refId)
			}
		} else {
			tab.refStr[refCol] = append(tab.refStr[refCol], refUid)
		}
	}

	return
}

func ParseRelation[V *int|*string](res V, rows *sql.Rows) error {
	if err := rows.Scan(&res); err != nil {
		return err
	}

	return nil
}

func (tab *Table[V]) PushRelation(relation Relation) {
	tab.relationList = append(tab.relationList, relation)
}

func (tab *Table[V]) RelationList() []Relation {
	return tab.relationList
}

func (tab *Table[V]) PushDep(depCol string, isInt bool, depId int, depUid string) {

}

//func (tab *Table[V]) Where() string {
	//tab.whereByStr()
	//tab.whereByInt()

	//if len(tab.groupByList) == 0 {
		//return ""
	//}

	//return strings.Join(tab.queryList, " OR ") + " GROUP BY " + strings.Join(tab.groupByList, ", ")
//}

//func (tab *Table[V]) whereByStr() {
	//for rel, valList := range tab.depStr {
		//var list []string
		//for _, val := range valList {
			//if val != nil {
				//list = append(list, "'" + *val + "'")
			//}
		//}
		//if len(list) > 0 {
			//tab.queryList = append(tab.queryList, fmt.Sprintf("%s in (%s)", rel.Col(), strings.Join(list, ", ")))
			//tab.groupByList = append(tab.groupByList, rel.Col())
		//}
	//}
//}

//func (tab *Table[V]) whereByInt() {
	//for rel, valList := range tab.depInt {
		//var list []string
		//for _, val := range valList {
			//if val != nil {
				//list = append(list, strconv.Itoa(*val))
			//}
		//}
		//if len(list) > 0 {
			//tab.queryList = append(tab.queryList, fmt.Sprintf("%s in (%s)", rel.Col(), strings.Join(list, ", ")))
			//tab.groupByList = append(tab.groupByList, rel.Col())
		//}
	//}
//}

//func (tab *Table[V]) whereRefByInt() (where []string) {
	//for colName, valList := range tab.refInt {
		//var list []string
		//for _, val := range valList {
			//if val != nil {
				//list = append(list, strconv.Itoa(*val))
			//}
		//}

		//if len(list) > 0 {
			//where = append(where, fmt.Sprintf("%s IN (%s)"), colName, strings.Join(list, ", "))
		//}
	//}

	//return
//}

//func (tab *Table[V]) whereRefByStr() (where []string) {
	//for colName, valList := range tab.refStr {
		//var list []string
		//for _, val := range valList {
			//if val != nil {
				//list = append(list, "'" + *val + "'")
			//}
		//}

		//if len(list) > 0 {
			//where = append(where, fmt.Sprintf("%s IN (%s)", colName, strings.Join(list, ", ")))
		//}
	//}

	//return
//}
