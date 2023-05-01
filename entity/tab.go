package entity

type TabInterface interface {
	Name() string
	isExist(valList []*Value) bool
	isUsed(valList []*Value) bool
	Rows() []*Row
	Push(valList []*Value)
	Pull() []*Row
	Deep(rel RelationInterface) int
}

type Tab struct {
	name         string
	rows         []*Row
	countRelDeep map[RelationInterface]int
}

func NewTab(tabName string) *Tab {
	return &Tab{
		name:         tabName,
		countRelDeep: make(map[RelationInterface]int),
	}
}

func (tab *Tab) Name() string {
	return tab.name
}

func (tab *Tab) Rows() []*Row {
	return tab.rows
}

func (tab *Tab) Pull() (list []*Row) {
	for _, row := range tab.Rows() {
		if !row.IsUsed() {
			row.ApplyUsed()
			list = append(list, row)
		}
	}

	return
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

func (tab *Tab) Deep(rel RelationInterface) int {
	tab.countRelDeep[rel]++

	return tab.countRelDeep[rel]
}
