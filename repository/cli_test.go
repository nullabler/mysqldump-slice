package repository

import (
	"fmt"
	"mysqldump-slice/addapter"
	"testing"
	"time"
)

func TestRmFile_empty(t *testing.T) {
	conf := &Conf{}
	exec := addapter.NewExecMock()
	want := "rm -f _.sql 2> /dev/null"

	runRmFile(t, conf, exec, want)
}

func TestRmFile_filename(t *testing.T) {
	conf := &Conf{
		Database: "test",
		File: File{
			Path: "./build/",
			Prefix: "short",
			DateFormat: "2006-01-02",
			Gzip: true,
		},
	}

	date := time.Now().Format(conf.File.DateFormat)

	want := fmt.Sprintf("rm -f %s%s_%s_%s.sql.gz 2> /dev/null", 
		conf.File.Path,
		conf.File.Prefix,
		date,
		conf.Database,
	)

	runRmFile(t, conf, addapter.NewExecMock(), want)
}

func TestExecDump(t *testing.T) {
	exec := addapter.NewExecMock()

	conf := &Conf{}
	runExecDump_err(t, conf, exec, "", "not found tmp file")

	conf.Tmp = "/tmp/1234"
	runExecDump_err(t, conf, exec, "", "fail auth")


	conf.User = "root"
	conf.Password = "1234"
	conf.Host = "db_test"
	runExecDump(t, conf, exec, "test_3", "mysqldump -uroot -p1234 -h db_test --single-transaction test_3 >> /tmp/1234")

	conf.DefaultExtraFile = "./path/db.cnf"
	runExecDump(t, conf, exec, "test_2", "mysqldump --defaults-extra-file=./path/db.cnf --single-transaction test_2 >> /tmp/1234")

}

func TestSave(t *testing.T) {
	exec := addapter.NewExecMock()

	conf := &Conf{}
	runSave(t, conf, exec, "", "tmp file is empty")

	conf.Tmp = "/tmp/1234"
	runSave(t, conf, exec, "cp /tmp/1234 _.sql", "")

	conf = &Conf{
		Database: "test",
		Tmp: "/tmp/1234",
		File: File{
			Path: "./build/",
			Prefix: "short",
			DateFormat: "2006-01-02",
			Gzip: true,
		},
	}

	date := time.Now().Format(conf.File.DateFormat)
	runSave(t, conf, exec, "cat /tmp/1234 | gzip > ./build/short_"+date+"_test.sql.gz", "")
}

func newCliMock(t *testing.T, conf *Conf, exec addapter.ExecInterface) *Cli {
	cli, err := NewCli(conf, exec)
	if err != nil {
		t.Errorf("NewCli err: %s", err.Error())
	}

	return cli
}

func runSave(t *testing.T, conf *Conf, exec *addapter.ExecMock, want, wantErr string) {
	cli := newCliMock(t, conf, exec)

	filename, err := cli.Save()
	if err != nil {
		if err.Error() == wantErr {
			return
		}
		t.Errorf("This not expect err: %s", err.Error())
	}

	if filename != conf.Filename() {
		t.Errorf("Filename not equel got: %s, wanted: %s", filename, conf.Filename())
	}

	got := exec.Call()
	if got != want {
		t.Errorf("Got: %q, wanted: %q", got, want)
	}
}

func runExecDump(t *testing.T, conf *Conf, exec *addapter.ExecMock, call, want string) {
	cli := newCliMock(t, conf, exec)
	if err := cli.ExecDump(call); err != nil {
		t.Errorf("This not expect err: %s", err.Error())
	}

	got := exec.Call()
	if got != want {
		t.Errorf("Got: %q, wanted: %q", got, want)
	}
}

func runExecDump_err(t *testing.T, conf *Conf, exec *addapter.ExecMock, call, want string) {
	cli := newCliMock(t, conf, exec)
	err := cli.ExecDump(call)
	if err == nil {
		t.Error("Not found error")
	}

	got := err.Error()
	if got != want {
		t.Errorf("Got: %q, wanted: %q", got, want)
	}
}

func runRmFile(t *testing.T, conf *Conf, exec *addapter.ExecMock, want string) {
	cli := newCliMock(t, conf, exec)
	if err := cli.RmFile(); err != nil {
		t.Errorf("This not expect err: %s", err.Error())
	}

	got := exec.Call()
	if got != want {
		t.Errorf("Got: %q, wanted: %q", got, want)
	}
}
