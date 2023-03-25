package repository

import (
	"errors"
	"mysqldump-slice/entity"
)

type DbMock struct {
	point  string
	iter   int
	impact int

	flag    bool
	err     error
	valList []*entity.Value
	strList []string
	sql     SqlInterface
}

func NewDbMock(conf *Conf) *DbMock {
	return &DbMock{
		iter: 1,
		sql:  NewSql(conf),
	}
}

func (d *DbMock) Point(point string, iter int) {
	d.point = point
	d.iter = iter
	d.impact = 0
}

func (d *DbMock) Impact(point string) bool {
	if point != d.point {
		return false
	}

	d.impact++
	if d.impact == d.iter {
		d.impact = 0
		return true
	}

	return false
}

func (d *DbMock) Flag(flag bool) {
	d.flag = flag
}

func (d *DbMock) Err(msg string) {
	d.err = errors.New(msg)
}

func (d *DbMock) ValList(valList []*entity.Value) {
	d.valList = valList
}

func (d *DbMock) StrList(strList []string) {
	d.strList = strList
}

type DbMockWrapper struct {
	mock *DbMock
}

func NewDbMockWrapper(mock *DbMock) *DbMockWrapper {
	return &DbMockWrapper{
		mock: mock,
	}
}

func (db *DbMockWrapper) Close() {}

func (db *DbMockWrapper) IsIntByCol(string, string) (bool, error) {
	if !db.mock.Impact("IsIntByCol") {
		return false, nil
	}

	return db.mock.flag, db.mock.err
}

func (db *DbMockWrapper) LoadRelations(entity.CollectInterface) error {
	if !db.mock.Impact("LoadRelations") {
		return nil
	}

	return db.mock.err
}

func (db *DbMockWrapper) LoadTables(entity.CollectInterface) error {
	if !db.mock.Impact("LoadTables") {
		return nil
	}

	return db.mock.err
}

func (db *DbMockWrapper) PrimaryKeys(string) ([]string, error) {
	if !db.mock.Impact("PrimaryKeys") {
		return nil, nil
	}

	return db.mock.strList, db.mock.err
}

func (db *DbMockWrapper) LoadIds(string, bool, Specs, []string, int) ([]*entity.Value, error) {
	if !db.mock.Impact("LoadIds") {
		return nil, nil
	}

	return db.mock.valList, db.mock.err
}

func (db *DbMockWrapper) LoadDeps(string, string, entity.RelationInterface) ([]string, error) {
	if !db.mock.Impact("LoadDeps") {
		return nil, nil
	}

	return db.mock.strList, db.mock.err
}

func (db *DbMockWrapper) LoadPkByCol(string, string, []string, []string) ([]*entity.Value, error) {
	if !db.mock.Impact("LoadPkByCol") {
		return nil, nil
	}

	return db.mock.valList, db.mock.err
}

func (db *DbMockWrapper) Sql() SqlInterface {
	return db.mock.sql
}
