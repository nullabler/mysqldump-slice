package service

import (
	"fmt"
	"mysqldump-slice/config"
	"mysqldump-slice/entity"
	"mysqldump-slice/module"
	"mysqldump-slice/repository"
	"testing"
)

func TestRelations(t *testing.T) {
	conf := &config.Conf{}
	config.FillSimpleSpecs(conf)

	confDbMock := repository.ConfDbMock{}
	l := getLoader(conf, confDbMock)
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

func TestTables(t *testing.T) {
	conf := &config.Conf{}

	confDbMock := repository.ConfDbMock{}
	repository.FillPrimaryKeys(confDbMock)
	repository.FillValList(confDbMock)

	l := getLoader(conf, confDbMock)

	collect := entity.NewCollect()
	entity.FillTables(collect)

	if err := l.Tables(collect); err != nil {
		t.Error("Call tables returned error: ", err.Error())
	}

	if collect.Tab("user") == nil ||
		collect.Tab("product") == nil ||
		collect.Tab("category") == nil {
		t.Error("Fail add tab struct")
	}

	if collect.Tab("fail") != nil {
		t.Error("Fail tab struct without to call add")
	}

	if len(collect.PkList("user")) != 1 ||
		len(collect.PkList("product")) != 1 ||
		len(collect.PkList("category")) != 1 {
		t.Error("Fail contains pkList for User tab")
	}

	if len(collect.PkList("fail")) != 0 {
		t.Error("Fail pkList for broken tab")
	}

	if len(collect.Tab("user").Rows()) != 2 ||
		len(collect.Tab("product").Rows()) != 4 ||
		len(collect.Tab("category").Rows()) != 1 {
		t.Error("Fail incorrect load valList")
	}
}

func TestDependences(t *testing.T) {
	conf := &config.Conf{}

	confDbMock := repository.ConfDbMock{}
	repository.FillPrimaryKeys(confDbMock)
	repository.FillValList(confDbMock)

	l := getLoader(conf, confDbMock)

	collect := entity.NewCollect()
	entity.FillTables(collect)
	entity.FillTabList(collect)

	for _, table := range collect.Tables() {
		fmt.Println(table.Name, collect.Tab(table.Name).Pull())
		rows := collect.Tab(table.Name).Pull()
		for _, rel := range collect.RelList(table.Name) {
			if err := l.Dependences(collect, rel, table.Name, rows); err != nil {
				t.Error(err)
			}
		}
	}
}

func getLoader(conf *config.Conf, confDbMock repository.ConfDbMock) *Loader {
	return NewLoader(
		conf,
		repository.NewDbMockWrapper(conf, confDbMock),
		repository.NewCliMock(),
		module.NewLogMock(),
	)
}
