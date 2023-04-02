package entity

import "fmt"

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
		return fmt.Sprintf("\\`%s\\` = '%s'", v.key, v.val)
	}

	return fmt.Sprintf("`%s` = '%s'", v.key, v.val)
}

func (v *Value) contains(valList []*Value) bool {
	for _, val := range valList {
		if v.key == val.key && v.val == val.val {
			return true
		}
	}

	return false
}
