package application

import (
	"mysqldump-slice/config"
	"mysqldump-slice/entity"
	"mysqldump-slice/module"
	"mysqldump-slice/repository"
	"mysqldump-slice/service"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	conf   *config.Conf
	loader *service.Loader
	dumper *service.Dumper
	log    module.LogInterface

	index []string
	pool  map[string]Callback
}

type Callback func(*entity.Collect) error

func NewApp(conf *config.Conf, log module.LogInterface, db repository.DbInterface, cli repository.CliInterface) *App {
	return &App{
		conf:   conf,
		loader: service.NewLoader(conf, db, cli, log), dumper: service.NewDumper(conf, cli, db, log), log: log,
		pool: make(map[string]Callback),
	}
}

func (app *App) Run() {
	app.initPool()

	collect := entity.NewCollect()
	app.call(collect)

	app.flush()
}

func (app *App) initPool() {
	app.loadRelationsPool()
	app.sortPool()
	app.loadTablesPool()
	app.removeDumpFilePool()
	app.dumpStructPool()
	app.dumpFullPool()
	app.loadDepAndDumpSlicePool()
}

func (app *App) call(c *entity.Collect) {
	for _, label := range app.index {
		app.log.Infof("%s......Start", label)
		if err := app.pool[label](c); err != nil {
			app.log.Error(err)
		}
		app.log.Infof("%s......Done", label)
	}
}

func (app *App) addPool(label string, fn Callback) {
	app.pool[label] = fn
	app.index = append(app.index, label)
}

func (app *App) flush() {
	if err := app.dumper.Save(); err != nil {
		app.log.Error(err)
	}

	filename, err := app.dumper.Filename()
	if err != nil {
		app.log.Error(err)
		return
	}

	app.log.State(filename)
}

func (app *App) loadRelationsPool() {
	app.addPool(
		"Load relations",
		func(collect *entity.Collect) error {
			return app.loader.Relations(collect)
		},
	)
}

func (app *App) sortPool() {
	app.addPool(
		"Sort",
		func(collect *entity.Collect) error {
			err := app.loader.Weight(collect)
			if err == nil {
				service.CallNormalize(collect)
			}

			return err
		},
	)
}

func (app *App) loadTablesPool() {
	app.addPool(
		"Load tables",
		func(collect *entity.Collect) error {
			return app.loader.Tables(collect)
		},
	)
}

func (app *App) removeDumpFilePool() {
	app.addPool(
		"Remove dump file",
		func(collect *entity.Collect) error {
			return app.dumper.RmFile()
		},
	)
}

func (app *App) dumpStructPool() {
	app.addPool(
		"Dump struct",
		func(collect *entity.Collect) error {
			return app.dumper.Struct()
		},
	)
}

func (app *App) dumpFullPool() {
	app.addPool(
		"Dump full data",
		func(collect *entity.Collect) error {
			return app.dumper.Full()
		},
	)
}

func (app *App) loadDepAndDumpSlicePool() {
	app.addPool(
		"Load dependences and dump slice data",
		func(collect *entity.Collect) error {
			return app.loadDepAndDumpSlice(collect)
		},
	)
}

func (app *App) loadDepAndDumpSlice(collect *entity.Collect) error {
	isLoop := true
	for {
		if isLoop {
			isLoop = false
		} else {
			break
		}

		for _, table := range collect.Tables() {
			tab := collect.Tab(table.Name)
			rows := tab.Pull()

			if app.log.ProfRowList(tab.Name(), rows, false, false, true) {
				app.log.Printf("Find in Application Table: %s RowList: %s", tab.Name(), app.conf.Profiler.Val)
			}

			lenRows := len(rows)
			if lenRows == 0 {
				continue
			} else {
				isLoop = true
			}

			app.log.Infof("- %s (%d)......Loading", tab.Name(), lenRows)

			for _, rel := range collect.RelList(tab.Name()) {
				if !app.deep(tab, rel) {
					app.log.Infof("^ Skip deep RefTab: %s RefCol: %s", rel.RefTab(), rel.RefCol())
					continue
				}

				if err := app.loader.Dependences(collect, rel, tab.Name(), rows); err != nil {
					return err
				}
			}

			if err := app.dumper.Slice(collect, tab.Name(), rows); err != nil {
				return err
			}

			app.log.Infof("- %s......Done", tab.Name())
		}
	}

	return nil
}

func (app *App) deep(tab entity.TabInterface, rel entity.RelationInterface) bool {
	spec, ok := app.conf.Specs(tab.Name())
	if !ok {
		return true
	}

	for _, fk := range spec.Fk {
		if fk.Deep > 0 && fk.FkTab == rel.RefTab() {
			return tab.Deep(rel) <= fk.Deep
		}
	}

	return true
}
