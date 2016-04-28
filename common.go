package util

import (
	"strings"
	"fmt"
	"strconv"
	"errors"
	"bytes"
	"text/template"
)

func ParseVersion(pineVersionStr string) (int64, int64) {
	versions := strings.Split(pineVersionStr, ".")
	if len(versions) == 1 {
		maj := versions[0]
		return mustParseInt(maj), 0

	} else if len(versions) == 2 {
		maj := versions[0]
		min := versions[1]
		return mustParseInt(maj), mustParseInt(min)

	} else {
		panic(errors.New(fmt.Sprintf("Bad version, cannot parse %v", pineVersionStr)))
	}
}

func mustParseInt(s string) int64 {
	var pineVersion, err = strconv.ParseInt(s, 10, 0)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Cannot parse int (%s), reason %v", s, err)))
	}
	return pineVersion
}

func MustExecuteTemplate(t *template.Template, data interface{}) string {
	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

/*
Returns element index of a slice (of any type!) for which given predicate is true. Returns -1 if not found.
Taken from http://stackoverflow.com/a/18203895
xs := []int{2, 4, 6, 8}
ys := []string{"C", "B", "K", "A"}

fmt.Println(
    FindIndex(len(xs), func(i int) bool { return xs[i] == 5 }),
    FindIndex(len(xs), func(i int) bool { return xs[i] == 6 }),
    FindIndex(len(ys), func(i int) bool { return ys[i] == "Z" }),
    FindIndex(len(ys), func(i int) bool { return ys[i] == "A" }))
*/
func FindIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}
