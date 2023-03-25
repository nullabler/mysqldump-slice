package service

import (
	"mysqldump-slice/entity"
	"mysqldump-slice/repository"
	"testing"
)

func TestRelations(t *testing.T) {
	conf := &repository.Conf{}
	repository.FillSimpleSpecs(conf)

	dbMock := repository.NewDbMock(conf)
	l := getLoader(conf, dbMock)
	collect := entity.NewCollect()

	if len(collect.AllRelList()) != 0 {
		t.Error("Fail is not empty relList")
	}

	if err := l.Relations(collect); err != nil {
		t.Error("Fail is not return error")
	}

	if len(collect.AllRelList()) != 1 {
		t.Error("Fail is not first relList")
	}

	if len(collect.AllRelList()["test"]) != 3 {
		t.Error("Fail push relation from config Specs")
	}
}

func getLoader(conf *repository.Conf, dbMock *repository.DbMock) *Loader {
	return NewLoader(
		conf,
		repository.NewDbMockWrapper(dbMock),
		repository.NewCliMock(),
		NewLogMock(),
	)
}
