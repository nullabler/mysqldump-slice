package entity

type TableInterface interface {
	Up()
}

type Table struct {
	Name   string
	Weight int
}

func NewTable(name string) *Table {
	return &Table{
		Name:   name,
		Weight: 0,
	}
}

func (t *Table) Up() {
	t.Weight += 1
}

type TableList []*Table

func (s TableList) Len() int {
	return len(s)
}

func (s TableList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
