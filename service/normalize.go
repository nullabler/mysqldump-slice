package service

import (
	"mysqldump-slice/entity"
	"sort"
)

type Normalize struct {
	entity.TableList
}

func CallNormalize(collect entity.CollectInterface) {
	sort.Sort(Normalize{collect.Tables()})
}

func (n Normalize) Less(i, j int) bool {
	return n.TableList[i].Weight < n.TableList[j].Weight
}
