package relationship

import "database/sql"

type Table struct {
	depIds  []int
	depUids []string
}

func (tab *Table) Parse(isInt bool, rows *sql.Rows) (err error) {
	var id int
	var uid string

	if isInt {
		err = rows.Scan(&id)
	} else {
		err = rows.Scan(&uid)
	}

	if err == nil {
		if isInt {
			tab.depIds = append(tab.depIds, id)
		} else {
			tab.depUids = append(tab.depUids, uid)
		}
	}
	return
}
