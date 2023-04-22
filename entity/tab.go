package entity

type TabInterface interface {
	isExist(valList []*Value) bool
	isUsed(valList []*Value) bool
	Rows() []*Row
	Push(valList []*Value)
	Pull() []*Row
}

type Tab struct {
	name string
	rows []*Row
}

func NewTab(tabName string) *Tab {
	return &Tab{
		name: tabName,
	}
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
