package repository

import (
	"fmt"
	"mysqldump-slice/entity"
	"strings"
)

type SqlInterface interface {
	Where([]*entity.Row, bool) []string
	WhereSlice([]*entity.Row, bool) []string
}

type Sql struct {
	conf *Conf
}

func NewSql(conf *Conf) *Sql {
	return &Sql{
		conf: conf,
	}
}

func (s *Sql) Where(rowList []*entity.Row, isEscape bool) []string {
	where := []string{}
	for _, row := range rowList {

		itemWhere := []string{}
		for _, val := range row.ValList() {
			itemWhere = append(itemWhere, val.Sprint(isEscape))
		}

		where = append(where, fmt.Sprintf("(%s)", strings.Join(itemWhere, " AND ")))
	}

	return where
}

func (s *Sql) WhereSlice(rowList []*entity.Row, isEscaape bool) (chunks []string) {
	whereList := s.Where(rowList, isEscaape)
	lenWhereList := len(whereList)

	for i := 0; i < lenWhereList; i += s.conf.LimitCli {
		end := i + s.conf.LimitCli

		if end > lenWhereList {
			end = lenWhereList
		}

		chunks = append(chunks, strings.Join(whereList[i:end], " OR "))
	}

	return
}
