package util

import (
	"net/http"
	"strings"
	"fmt"
	"strconv"
	"time"
	"math/rand"
)

func ParseVersion(pineVersionStr string) (uint64, uint64) {
	versions := strings.Split(pineVersionStr, ".")
	if len(versions) == 1 {
		maj := versions[0]
		return mustParseUint(maj), 0

	} else if len(versions) == 2 {
		maj := versions[0]
		min := versions[1]
		return mustParseUint(maj), mustParseUint(min)

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

func mustParseUint(s string) uint64 {
	var pineVersion, err = strconv.ParseUint(s, 10, 0)
	if err != nil {
		panic(CreateHttpError(http.StatusBadRequest, fmt.Sprintf("Cannot convert version %s to uint, reason %v", s, err)))
	}
	return pineVersion
}

func TimeNowUnixMillis() int64 {
	return time.Now().UnixNano() / 1000000
}


func GenerateId(n int) string {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
	}
	return string(b)
}