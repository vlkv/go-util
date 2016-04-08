package util

import (
	"net/http"
	"strings"
	"fmt"
	"strconv"
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
		// TODO: do not use http
		panic(CreateHttpError(http.StatusBadRequest, fmt.Sprintf("Bad version, cannot parse %v", pineVersionStr)))
	}
}

func mustParseInt(s string) int64 {
	var pineVersion, err = strconv.ParseInt(s, 10, 0)
	if err != nil {
		// TODO: do not use http
		panic(CreateHttpError(http.StatusBadRequest, fmt.Sprintf("Cannot parse int (%s), reason %v", s, err)))
	}
	return pineVersion
}

