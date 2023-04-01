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
	QueryLoadIds(string, *Specs, []string) (string, error)
	QueryLoadDeps(string, string, string, string) string
	QueryLoadPkByCol([]string, string, string, []string) string
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

func (s *Sql) WhereSlice(rowList []*entity.Row, isEscape bool) (chunks []string) {
	whereList := s.Where(rowList, isEscape)
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
		AND CONSTRAINT_NAME = ? 
		AND CONSTRAINT_SCHEMA = ?`
}

func (s *Sql) QueryIsIntByCol() string {
	return `SELECT data_type  
		FROM information_schema.columns 
		WHERE table_schema = ? 
		AND table_name = ? 
		AND column_name = ?;`
}

func (s *Sql) QueryLoadIds(tabName string, specs *Specs, pkList []string) (string, error) {
	if len(pkList) == 0 {
		return "", fmt.Errorf("Empty PrimaryKeyList for TabName: %s", tabName)
	}

	sort := s.sort(pkList, specs)

	return fmt.Sprintf("SELECT %s FROM `%s`%s ORDER BY %s DESC %s",
		s.wrapAndJoin(pkList),
		tabName,
		s.condition(specs),
		sort,
		s.limit(tabName, specs),
	), nil
}

func (s *Sql) sort(fields []string, specs *Specs) string {
	if specs != nil && len(specs.Sort) > 0 {
		fields = specs.Sort
	}

	return s.wrapAndJoin(fields)
}

func (s *Sql) wrapAndJoin(fields []string) string {
	return strings.Join(s.wrapKeys(fields, "`"), ", ")
}

func (s *Sql) condition(specs *Specs) string {
	if specs != nil && len(specs.Condition) > 0 {
		return fmt.Sprintf(" WHERE %s", specs.Condition)
	}

	return ""
}

func (s *Sql) limit(tabName string, specs *Specs) string {
	if s.conf.IsFull(tabName) {
		return ""
	}

	limit := s.conf.Tables.Limit
	if specs != nil && specs.Limit > 0 {
		limit = specs.Limit
	}

	return fmt.Sprintf("LIMIT %d", limit)
}

func (s *Sql) QueryLoadDeps(col, tabName, where, limit string) string {
	return fmt.Sprintf("SELECT `%s` FROM `%s` WHERE %s %s",
		col,
		tabName,
		where,
		limit,
	)
}

func (s *Sql) QueryLoadPkByCol(keyList []string, tabName, tabCol string, valList []string) string {
	return fmt.Sprintf("SELECT %s FROM `%s` WHERE `%s` IN (%s)",
		s.wrapAndJoin(keyList), tabName, tabCol, strings.Join(valList, ", "))
}

func (s *Sql) wrapKeys(keys []string, wrapSym string) (list []string) {
	for _, key := range keys {
		list = append(list, wrapSym+key+wrapSym)
	}

	return
}
