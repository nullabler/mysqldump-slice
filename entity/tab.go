package entity

type TabInterface interface {
	Pool() map[string]Key
	Push(string, string)
	Pull() map[string][]string
}

type Tab struct {
	name string
	exist map[string]Exist
	pool  map[string]Key
}

type Exist map[string]bool

type Key map[string]*Val

type Val struct {
	used bool
}

func NewTab(tabName string) *Tab {
	return &Tab{
		name: tabName,
		exist: make(map[string]Exist),
		pool:  make(map[string]Key),
	}
}

func (tab *Tab) Push(col, val string) {
	if tab.exist[col][val] {
		return
	}

	if tab.exist[col] == nil {
		tab.exist[col] = make(Exist)
	}
	tab.exist[col][val] = true

	if tab.pool[col] == nil {
		tab.pool[col] = make(Key)
	}
	tab.pool[col][val] = &Val{
		used: false,
	}
}

func (tab *Tab) Pool() map[string]Key {
	return tab.pool
}

func (tab *Tab) Pull() map[string][]string {
	list := make(map[string][]string)

	for col, pool := range tab.pool {
		for val, key := range pool {
			if !key.used {
				list[col] = append(list[col], val)
				key.used = true
			}
		}
	}

	return list
}
