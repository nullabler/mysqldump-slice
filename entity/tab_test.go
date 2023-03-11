package entity

import (
	"testing"
)

func TestIsExistForSimple(t *testing.T) {
	tab := NewTab("test")

	for _, i := range []string{"1", "2", "3"} {
		valList := []*Value{}
		valList = append(valList, NewValue("id", i))
		tab.rows = append(tab.rows, NewRow(valList))
	}

	valListOk := []*Value{}
	valListOk = append(valListOk, NewValue("id", "1"))
	if !tab.isExist(valListOk) {
		t.Error("Fail isExist where is not exist simple")
	}

	valListFail := []*Value{}
	valListFail = append(valListFail, NewValue("id", "5"))
	if tab.isExist(valListFail) {
		t.Error("Fail isExist where is exist simple")
	}
}

func TestIsExistForComposition(t *testing.T) {
	tab := NewTab("test")

	for _, i := range []string{"1", "2", "3"} {
		valList := []*Value{}
		valList = append(valList, NewValue("user_id", i))
		valList = append(valList, NewValue("category_id", i+"5"))
		tab.rows = append(tab.rows, NewRow(valList))
	}

	valListOk := []*Value{}
	valListOk = append(valListOk, NewValue("user_id", "1"))
	valListOk = append(valListOk, NewValue("category_id", "15"))
	if !tab.isExist(valListOk) {
		t.Error("Fail isExist where is not exist composition")
	}

	valListFail_1 := []*Value{}
	valListFail_1 = append(valListFail_1, NewValue("user_id", "9"))
	valListFail_1 = append(valListFail_1, NewValue("category_id", "15"))
	if tab.isExist(valListFail_1) {
		t.Error("Fail isExist where is exist composition (1)")
	}

	valListFail_2 := []*Value{}
	valListFail_2 = append(valListFail_2, NewValue("user_id", "1"))
	valListFail_2 = append(valListFail_2, NewValue("category_id", "19"))
	if tab.isExist(valListFail_2) {
		t.Error("Fail isExist where is exist composition (2)")
	}

	valListFail_3 := []*Value{}
	valListFail_3 = append(valListFail_3, NewValue("user_id", "0"))
	valListFail_3 = append(valListFail_3, NewValue("category_id", "65"))
	if tab.isExist(valListFail_3) {
		t.Error("Fail isExist where is exist composition (3)")
	}
}

func TestIsUsedForSimple(t *testing.T) {
	tab := NewTab("test")

	for _, i := range []string{"1", "2", "3"} {
		valList := []*Value{}
		valList = append(valList, NewValue("id", i))
		tab.rows = append(tab.rows, NewRow(valList))
	}
	tab.rows[1].used = true

	vlOk := []*Value{}
	vlOk = append(vlOk, NewValue("id", "2"))
	if !tab.isUsed(vlOk) {
		t.Error("Fail isUsed where is used simple")
	}

	vlFail := []*Value{}
	vlFail = append(vlFail, NewValue("id", "3"))
	if tab.isUsed(vlFail) {
		t.Error("Fail isUsed where is not used simple")
	}
}

func TestIsUsedForComposition(t *testing.T) {
	tab := NewTab("test")

	for _, i := range []string{"1", "2", "3"} {
		valList := []*Value{}
		valList = append(valList, NewValue("user_id", i))
		valList = append(valList, NewValue("category_id", i+"5"))
		tab.rows = append(tab.rows, NewRow(valList))
	}
	tab.rows[2].used = true

	vlOk := []*Value{}
	vlOk = append(vlOk, NewValue("user_id", "3"))
	vlOk = append(vlOk, NewValue("category_id", "35"))
	if !tab.isUsed(vlOk) {
		t.Error("Fail isUsed where is used composition")
	}

	vlFail := []*Value{}
	vlFail = append(vlFail, NewValue("user_id", "2"))
	vlFail = append(vlFail, NewValue("category_id", "25"))
	if tab.isUsed(vlFail) {
		t.Error("Fail isUsed where is not used composition")
	}
}

func TestPush(t *testing.T) {
	tab := NewTab("test")

	for _, i := range []string{"1", "2", "3"} {
		valList := []*Value{}
		valList = append(valList, NewValue("id", i))
		tab.rows = append(tab.rows, NewRow(valList))
	}

	if len(tab.Rows()) != 3 {
		t.Error("Fail count after push simple")
	}

	for _, i := range []string{"1", "2", "3", "8"} {
		valList := []*Value{}
		valList = append(valList, NewValue("user_id", i))
		valList = append(valList, NewValue("category_id", i+"5"))
		tab.rows = append(tab.rows, NewRow(valList))
	}

	if len(tab.Rows()) != 7 {
		t.Error("Fail count after push mixed")
	}
}
