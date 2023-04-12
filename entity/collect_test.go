package entity

import (
	"testing"
)

func TestPushValListSuccess(t *testing.T) {
	c := NewCollect()
	c.PushPkList("test", []string{"uuid", "ulid"})

	valList := [][]*Value{}
	valList = append(valList, []*Value{
		NewValue("uuid", "xxxx-yyyy-0001"),
		NewValue("ulid", "dddd-cccc-0001"),
	})
	valList = append(valList, []*Value{
		NewValue("uuid", "xxxx-yyyy-0004"),
		NewValue("ulid", "dddd-cccc-0008"),
	})

	if c.Tab("test") != nil {
		t.Error("TabList should be nil")
	}

	c.PushTab("test")
	if c.Tab("test") == nil {
		t.Error("TabList is nil after PushTab")
	}

	c.PushValList("test", valList)
}

func TestPushValListFail(t *testing.T) {
	c := NewCollect()
	c.PushPkList("test", []string{"uuid", "ulid"})

	valList := [][]*Value{}
	valList = append(valList, []*Value{NewValue("uuid", "4")})

	if c.Tab("test") != nil {
		t.Error("TabList should be nil")
	}

	c.PushTab("test")
	if c.Tab("test") == nil {
		t.Error("TabList is nil after PushTab")
	}

	c.PushValList("test", valList)
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
