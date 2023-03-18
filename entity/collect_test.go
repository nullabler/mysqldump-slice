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

func TestIsValid(t *testing.T) {
	c := NewCollect()
	c.PushPkList("test", []string{"uuid", "ulid"})

	valList := []*Value{}
	if c.isValid("test", valList) {
		t.Error("Fail is not correct PK for ValList")
	}

	valList = append(valList, NewValue("uuid", "4"))
	if c.isValid("test", valList) {
		t.Error("Fail is not correct PK for ValList")
	}

	valList = append(valList, NewValue("ulid", "2"))
	if !c.isValid("test", valList) {
		t.Error("Fail is correct PK for ValList")
	}

	valList = append(valList, NewValue("id", "9"))
	if !c.isValid("test", valList) {
		t.Error("Fail is correct PK for ValList")
	}
}

func TestIsPk(t *testing.T) {
	keys := []string{"uuid", "ulid"}
	c := NewCollect()
	c.PushPkList("test", keys)

	if c.IsPk("test", "id") {
		t.Error("Fail is not PK for 'id'")
	}

	for _, key := range keys {
		if !c.IsPk("test", key) {
			t.Errorf("Fail is PK for '%s'", key)
		}
	}
}
