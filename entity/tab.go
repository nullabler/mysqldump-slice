package entity

type Tab struct {
	name string
	keys map[string][]string
}

func NewTab(tabName string) *Tab {
	return &Tab{
		name: tabName,
		keys: make(map[string][]string),
	}
}

func (tab *Tab) Push(col, val string) {
	for _, item := range tab.keys[col] {
		if item == val {
			return
		}
	}

	tab.keys[col] = append(tab.keys[col], val)
}


func (tab *Tab) Keys() map[string][]string {
	return tab.keys
}
