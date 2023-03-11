package entity

type TabInterface interface {
	Pool() map[*Fields]Val
	Push(*Fields, string)
	Pull() map[*Fields][]string
}

type Tab struct {
	name string
	pool map[*Fields]Val
}

type Val map[string]*Spec

type Fields []string

type Spec struct {
	used bool
}

func NewTab(tabName string) *Tab {
	return &Tab{
		name: tabName,
		pool: make(map[*Fields]Val),
	}
}

func (tab *Tab) Push(fl *Fields, val string) {
	if tab.pool[fl][val] != nil {
		return
	}

	if tab.pool[fl] == nil {
		tab.pool[fl] = make(Val)
	}
	tab.pool[fl][val] = &Spec{
		used: false,
	}
}

func (tab *Tab) Pool() map[*Fields]Val {
	return tab.pool
}

func (tab *Tab) Pull() map[*Fields][]string {
	list := make(map[*Fields][]string)

	for f, pool := range tab.pool {
		for val, key := range pool {
			if !key.used {
				list[f] = append(list[f], val)
				key.used = true
			}
		}
	}

	return list
}
