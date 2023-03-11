package entity

import (
	"testing"
)

func TestPush(t *testing.T) {
	tab := NewTab("simple")
	if len(tab.pool) != 0 {
		t.Error("Simple tab is not empty")
	}

	fl := &Fields{"id"}

	tab.Push(fl, "1")
	if len(tab.pool[fl]) != 1 {
		t.Error("Fail first push for simple tab")
	}

	tab.Push(fl, "1")
	if len(tab.pool[fl]) != 1 {
		t.Error("Dublicate simple tab")
	}

	tab.Push(fl, "2")
	if len(tab.pool[fl]) != 2 {
		t.Error("Fail push for simple tab")
	}
}
