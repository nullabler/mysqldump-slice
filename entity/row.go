package entity

type Row struct {
	valList []*Value
	used    bool
}

func NewRow(valList []*Value) *Row {
	return &Row{
		valList: valList,
		used:    false,
	}
}

func (r *Row) ValList() []*Value {
	return r.valList
}

func (r *Row) IsUsed() bool {
	return r.used
}

func (r *Row) ApplyUsed() {
	r.used = true
}
