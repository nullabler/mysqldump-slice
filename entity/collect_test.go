package entity

import (
	"testing"
)

func TestPushValListSuccess(t *testing.T) {
	c := NewCollect()
	c.PushPkList("test", []string{"uuid", "ulid"})

	valList := []*Value{
		NewValue("uuid", "4"),
		NewValue("ulid", "9"),
	}

	if c.Tab("test") != nil {
		t.Error("TabList should be nil")
	}

	c.PushTab("test")
	if c.Tab("test") == nil {
		t.Error("TabList is nil after PushTab")
	}

	if er := c.PushValList("test", valList); er != nil {
		t.Error("Fail push valList to collect")
	}
}

func TestPushValListFail(t *testing.T) {
	c := NewCollect()
	c.PushPkList("test", []string{"uuid", "ulid"})

	valList := []*Value{
		NewValue("uuid", "4"),
	}

	if c.Tab("test") != nil {
		t.Error("TabList should be nil")
	}

	c.PushTab("test")
	if c.Tab("test") == nil {
		t.Error("TabList is nil after PushTab")
	}

	if er := c.PushValList("test", valList); er == nil {
		t.Error("Not full pk and not returned error")
	}
}
