package entity

import (
	"fmt"
	"mysqldump-slice/helper"
)

type Value struct {
	key string
	val string
}

func NewValue(key, val string) *Value {
	return &Value{
		key: key,
		val: val,
	}
}

func (v *Value) Sprint(isEscape bool) string {
	if isEscape {
		return fmt.Sprintf("\\`%s\\` = %s", v.key, v.Val(true))
	}

	return fmt.Sprintf("`%s` = %s", v.key, v.Val(true))
}

func (v *Value) contains(valList []*Value) bool {
	for _, val := range valList {
		if v.key == val.key && v.val == val.val {
			return true
		}
	}

	return false
}

func (v *Value) Key() string {
	return v.key
}

func (v *Value) Val(isWrap bool) string {
	if isWrap {
		return fmt.Sprintf("'%s'", helper.EscapeVal(v.val))
	}

	return v.val
}
