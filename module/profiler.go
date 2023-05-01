package module

import (
	"mysqldump-slice/config"
	"mysqldump-slice/entity"
	"mysqldump-slice/helper"
)

type Profiler struct {
	conf              *config.Conf
	tabNameParentList []string
}

func NewProfiler(conf *config.Conf) *Profiler {
	return &Profiler{
		conf: conf,
	}
}

func (p *Profiler) Active() bool {
	return p.conf.Profiler.Active
}

func (p *Profiler) TabName(tabName string) bool {
	if len(p.conf.Profiler.Table) == 0 {
		return false
	}

	return tabName == p.conf.Profiler.Table
}

func (p *Profiler) TraceDep(tabName, refTab string) bool {
	if len(p.conf.Profiler.TraceDep) == 0 {
		return true
	}

	if refTab != p.conf.Profiler.TraceDep {
		return false
	}

	if !helper.SliceIsExist(tabName, p.tabNameParentList) {
		p.tabNameParentList = append(p.tabNameParentList, tabName)

		return true
	}

	return false
}

func (p *Profiler) RowList(rows []*entity.Row, enableUsed, enableKey, enableVal bool) bool {
	for _, row := range rows {
		if enableUsed && row.IsUsed() {
			continue
		}

		for _, val := range row.ValList() {
			isKey := false
			if !enableKey {
				isKey = true
			} else if len(p.conf.Profiler.Key) == 0 {
				isKey = true
			}

			if !isKey && val.Key() == p.conf.Profiler.Key {
				isKey = true
			}

			isVal := false
			if !enableVal {
				isVal = true
			} else if len(p.conf.Profiler.Val) == 0 {
				isVal = true
			}

			if !isVal && val.Val(false) == p.conf.Profiler.Val {
				isVal = true
			}

			if isKey && isVal {
				p.conf.Profiler.Trace = true
				return true
			}
		}
	}

	return false
}

func (p *Profiler) ValList(data []entity.ValList, enableVal, enableRefVal bool) (valList entity.ValList) {
	for _, list := range data {
		for _, v := range list {
			if (enableVal && v.Val(false) == p.conf.Profiler.Val) ||
				(enableRefVal && v.Val(false) == p.conf.Profiler.RefVal) {
				valList = append(valList, v)
			}
		}
	}

	return
}

func (p *Profiler) Rel(rel entity.RelationInterface, enableRefTab, enableRefCol bool) bool {
	isRefTab := false
	if !enableRefTab {
		isRefTab = true
	} else if rel.RefTab() == p.conf.Profiler.RefTab {
		isRefTab = true
	}

	isRefKey := false
	if !enableRefCol {
		isRefKey = true
	} else if rel.RefCol() == p.conf.Profiler.RefKey {
		isRefKey = true
	}

	return isRefTab && isRefKey
}

func (p *Profiler) StrList(list []string, enableVal, enableRefVal bool) bool {
	isVal := false
	if !enableVal {
		isVal = true
	}

	isRefVal := false
	if !enableRefVal {
		isRefVal = true
	}
	for _, v := range list {
		if !isVal && v == p.conf.Profiler.Val {
			isVal = true
		}
		if !isRefVal && v == p.conf.Profiler.RefVal {
			isRefVal = true
		}
	}

	return isVal && isRefVal
}
