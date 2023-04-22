package entity

import "testing"

func TestSprint(t *testing.T) {
	v := NewValue("id", "2")

	exp := "`id` = '2'"
	got := v.Sprint(false)
	if exp != got {
		t.Errorf("Fail Sprint Exp: %s Got: %s", exp, got)
	}

	exp = "\\`id\\` = '2'"
	got = v.Sprint(true)
	if exp != got {
		t.Errorf("Fail Sprint Exp: %s Got: %s", exp, got)
	}
}

func TestContains(t *testing.T) {
	v := NewValue("uuid", "xxxx-cc-yyyyyyy")

	valList := []*Value{
		NewValue("uuid", "xxxx-cc-yyyyyyy"),
	}

	if !v.contains(valList) {
		t.Error("Fail is correct contains")
	}

	valList = append(valList, NewValue("id", "2"))
	if !v.contains(valList) {
		t.Error("Fail is correct contains")
	}

	valList = []*Value{
		NewValue("cat_id", "3"),
		NewValue("sel_id", "7"),
	}
	if v.contains(valList) {
		t.Error("Fail is not correct contains")
	}
}
