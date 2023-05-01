package entity

type Row struct {
	valList ValList
	used    bool
}

func NewRow(valList ValList) *Row {
	return &Row{
		valList: valList,
		used:    false,
	}
}

func (r *Row) ValList() ValList {
	return r.valList
}

func (r *Row) IsUsed() bool {
	return r.used
}

func (r *Row) ApplyUsed() {
	r.used = true
}
