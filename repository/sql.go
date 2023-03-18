package repository

import (
	"fmt"
	"mysqldump-slice/entity"
	"strings"
)

type SqlInterface interface {
	Where([]*entity.Row, bool) []string
	WhereSlice([]*entity.Row, bool) []string

	QueryRelations() string
	QueryFullTables(string) string
	QueryPrimaryKeys() string
	QueryIsIntByCol() string
	QueryLoadIds(string, string, bool, Specs, []string, int) string
	QueryLoadDeps(string, string, string, string) string
	QueryLoadPkByCol(string, string, string, []string) string
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

func (s *Sql) QueryRelations() string {
	return `SELECT fks.table_name as foreign_table,
			fks.referenced_table_name as primary_table,
			kcu.column_name as fk_column,
			kcu.referenced_column_name as ref_column
		FROM information_schema.referential_constraints fks
		JOIN information_schema.key_column_usage kcu
			ON fks.constraint_schema = kcu.table_schema
			AND fks.table_name = kcu.table_name
			AND fks.constraint_name = kcu.constraint_name
		WHERE fks.constraint_schema = ?
		GROUP BY fks.constraint_schema,
			fks.table_name,
			fks.unique_constraint_schema,
			fks.referenced_table_name,
			fks.constraint_name
		ORDER BY fks.constraint_schema, fks.table_name`
}

func (s *Sql) QueryFullTables(dbName string) string {
	return fmt.Sprintf("SHOW FULL TABLES FROM `%s`", dbName)
}

func (s *Sql) QueryPrimaryKeys() string {
	return `SELECT COLUMN_NAME 
		FROM information_schema.KEY_COLUMN_USAGE 
		WHERE TABLE_NAME = ? 
		AND CONSTRAINT_NAME = ?`
}

func (s *Sql) QueryIsIntByCol() string {
	return `SELECT data_type  
		FROM information_schema.columns 
		WHERE table_schema = ? 
		AND table_name = ? 
		AND column_name = ?;`
}

func (s *Sql) QueryLoadIds(key, tabName string, okSpecs bool, specs Specs, prKeyList []string, confLimit int) string {
	var sort string
	if okSpecs && len(specs.Sort) > 0 {
		sort = strings.Join(s.wrapKeys(specs.Sort, "`"), ", ")
	} else {
		sort = strings.Join(s.wrapKeys(prKeyList, "`"), ", ")
	}

	var condition, limit string
	if okSpecs && len(specs.Condition) > 0 {
		condition = "WHERE " + specs.Condition
		if specs.Limit > 0 {
			limit = fmt.Sprintf("LIMIT %d", specs.Limit)
		}
	}

	if len(condition) == 0 && confLimit > 0 {
		limit = fmt.Sprintf("LIMIT %d", confLimit)
		if okSpecs && specs.Limit > 0 {
			limit = fmt.Sprintf("LIMIT %d", specs.Limit)
		}
	}

	return fmt.Sprintf("SELECT `%s` FROM `%s` %s ORDER BY %s DESC %s",
		key, tabName, condition, sort, limit)
}

func (s *Sql) QueryLoadDeps(col, tabName, where, limit string) string {
	return fmt.Sprintf("SELECT `%s` FROM `%s` WHERE %s %s",
		col,
		tabName,
		where,
		limit,
	)
}

func (s *Sql) QueryLoadPkByCol(key, tabName, tabCol string, valList []string) string {
	return fmt.Sprintf("SELECT `%s` FROM `%s` WHERE `%s` IN (%s)",
		key, tabName, tabCol, strings.Join(valList, ", "))
}

func (s *Sql) wrapKeys(keys []string, wrapSym string) (list []string) {
	for _, key := range keys {
		list = append(list, wrapSym+key+wrapSym)
	}

	return
}
