package repository

import "mysqldump-slice/entity"

type DbMock struct {
}

func NewDbMock() *DbMock {
	return &DbMock{}
}

func (db *DbMock) Close() {

}

func (db *DbMock) IsIntByCol(string, string) (bool, error) {
	return true, nil
}

func (db *DbMock) LoadRelations(entity.CollectInterface) error {
	return nil
}

func (db *DbMock) LoadTables(entity.CollectInterface) error {
	return nil
}

func (db *DbMock) PrimaryKeys(string) ([]string, error) {
	return nil, nil
}

func (db *DbMock) LoadIds(string, bool, Specs, []string, int) ([]*entity.Value, error) {
	return nil, nil
}

func (db *DbMock) LoadDeps(string, string, entity.RelationInterface) ([]string, error) {
	return nil, nil
}

func (db *DbMock) LoadPkByCol(string, string, []string, []string) ([]*entity.Value, error) {
	return nil, nil
}

func (db *DbMock) Sql() SqlInterface {
	return nil
}
