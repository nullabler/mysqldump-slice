package service

import (
	"mysqldump-slice/config"
	"testing"
)

func TestFilename_WithDate(t *testing.T) {
	dumper := &Dumper{
		conf: &config.Conf{
			Database: "ubntshop",
			File: config.File{
				Path:       "./target/",
				Prefix:     "short",
				DateFormat: "2006",
				Gzip:       false,
			},
		},
	}

	filename, err := dumper.Filename()
	if err != nil {
		t.Fatalf("Filename err: %v", err)
	}

	if filename != "./target/short_2006_ubntshop.sql" {
		t.Fatalf("unexpected filename: %s", filename)
	}
}

func TestFilename_WithoutDate(t *testing.T) {
	dumper := &Dumper{
		conf: &config.Conf{
			Database: "ubntshop",
			File: config.File{
				Path:       "./target/",
				Prefix:     "short",
				DateFormat: "",
				Gzip:       false,
			},
		},
	}

	filename, err := dumper.Filename()
	if err != nil {
		t.Fatalf("Filename err: %v", err)
	}

	if filename != "./target/short_ubntshop.sql" {
		t.Fatalf("unexpected filename: %s", filename)
	}
}
