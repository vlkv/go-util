package util

import (
	"strings"
	"regexp"
)

func Abbrev(s string) (res string) {
	var maxLen = 128
	if len(s) <= maxLen {
		res = s
	} else {
		res = s[:maxLen / 2] + "<...>" + s[len(s) - maxLen / 2:]
	}
	res = strings.Replace(res, "\n", "\\n", -1)
	res = strings.Replace(res, "\r", "\\r", -1)
	return
}

/**
	This function helps with formatting in multiline string by
	deleting whitespaces followed by '|' in the beginning of the line
	Usage:
		stringWithCode := StripMargin(`
			|fun(x, y) =>
			|	res = x + y
			|	res
			|`)
 */
func StripMargin(s string) string {
	r := regexp.MustCompile(`(?m:^[ \t]*\|)`)
	return r.ReplaceAllLiteralString(s, "")
}
