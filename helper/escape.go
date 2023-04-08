package helper

import "strings"

func EscapeVal(val string) string {
	replace := map[string]string{"\\": "\\\\", "'": `\'`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z"}

	for b, a := range replace {
		val = strings.ReplaceAll(val, b, a)
	}

	return val
}
