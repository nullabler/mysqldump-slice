package repository

import (
	"errors"
	"mysqldump-slice/config"
	"mysqldump-slice/entity"
	"strconv"
)

type DbMock struct {
	flag    bool
	err     error
	valList [][]*entity.Value
	strList []string
}

func NewDbMock() *DbMock {
	return &DbMock{}
}

func (d *DbMock) Flag(flag bool) {
	d.flag = flag
}

func (d *DbMock) Err(msg string) {
	d.err = errors.New(msg)
}

func (d *DbMock) ValList(valList []*entity.Value) {
	d.valList = append(d.valList, valList)
}

func (d *DbMock) StrList(strList []string) {
	d.strList = strList
}

type ConfDbMock map[string]*DbMock

func KeyForMock(name string, iter int) string {
	return name + ":" + strconv.Itoa(iter)
}

type DbMockWrapper struct {
	mock ConfDbMock
	iter map[string]int

	sql SqlInterface
}

func NewDbMockWrapper(conf *config.Conf, confMock ConfDbMock) *DbMockWrapper {
	return &DbMockWrapper{
		mock: confMock,
		iter: make(map[string]int),
		sql:  NewSql(conf),
	}
}

func (db *DbMockWrapper) Impact(name string) *DbMock {
	_, ok := db.iter[name]
	if !ok {
		db.iter[name] = 0
	}

	db.iter[name]++

	return db.mock[KeyForMock(name, db.iter[name])]
}

func (db *DbMockWrapper) Close() {}

func (db *DbMockWrapper) IsIntByCol(string, string) (bool, error) {
	mock := db.Impact("IsIntByCol")
	if mock == nil {
		return false, nil
	}

	return mock.flag, mock.err
}

func (db *DbMockWrapper) LoadRelations(entity.CollectInterface) error {
	mock := db.Impact("LoadRelations")
	if mock == nil {
		return nil
	}

	return mock.err
}

func (db *DbMockWrapper) LoadTables(entity.CollectInterface) error {
	mock := db.Impact("LoadTables")
	if mock == nil {
		return nil
	}

	return mock.err
}

func (db *DbMockWrapper) PrimaryKeys(key string) ([]string, error) {
	mock := db.Impact("PrimaryKeys_" + key)
	if mock == nil {
		return nil, nil
	}

	return mock.strList, mock.err
}

func (db *DbMockWrapper) LoadIds(key string, s *config.Specs, l []string) ([][]*entity.Value, error) {
	mock := db.Impact("LoadIds_" + key)
	if mock == nil {
		return nil, nil
	}

	return mock.valList, mock.err
}

func (db *DbMockWrapper) LoadDeps(string, string, entity.RelationInterface) ([]string, error) {
	mock := db.Impact("LoadDeps")
	if mock == nil {
		return nil, nil
	}

	return mock.strList, mock.err
}

func (db *DbMockWrapper) LoadPkByCol(string, string, []string, []string) ([][]*entity.Value, error) {
	mock := db.Impact("LoadPkByCol")
	if mock == nil {
		return nil, nil
	}

	return mock.valList, mock.err
}

func (db *DbMockWrapper) Sql() SqlInterface {
	return db.sql
}
