package entity

type TabInterface interface {
	isExist([]*Value) bool
	isUsed([]*Value) bool
	Rows() []*Row
	Push([]*Value)
}

type Tab struct {
	name string
	rows []*Row
}

type Row struct {
	valList []*Value
	used    bool
}

type Value struct {
	key string
	val string
}

func NewValue(key, val string) *Value {
	return &Value{
		key: key,
		val: val,
	}
}

func (v *Value) contains(valList []*Value) bool {
	for _, val := range valList {
		if v.key == val.key && v.val == val.val {
			return true
		}
	}

	return false
}

func NewRow(valList []*Value) *Row {
	return &Row{
		valList: valList,
		used:    false,
	}
}

func NewTab(tabName string) *Tab {
	return &Tab{
		name: tabName,
	}
}

func (tab *Tab) Rows() []*Row {
	return tab.rows
}

func (tab *Tab) isExist(valList []*Value) bool {
	for _, row := range tab.rows {
		flag := true
		for _, val := range row.valList {
			if flag && !val.contains(valList) {
				flag = false
			}
		}

		if flag {
			return true
		}
	}

	return false
}

func (tab *Tab) isUsed(valList []*Value) bool {
	for _, row := range tab.rows {
		flag := true
		for _, val := range row.valList {
			if flag && !val.contains(valList) {
				flag = false
			}
		}

		if flag {
			return row.used
		}
	}

	return false
}

func (tab *Tab) Push(valList []*Value) {
	if tab.isExist(valList) {
		return
	}

	tab.rows = append(tab.rows, NewRow(valList))
}
