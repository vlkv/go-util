package util

import (
	"github.com/stretchr/testify/assert"
	"strings"
)

func EqualMultiline(t assert.TestingT, expected, actual string, msgAndArgs ...interface{}) bool {
	expectedArr := strings.Split(expected, "\n")
	actualArr := strings.Split(expected, "\n")
	return assert.Equal(t, expectedArr, actualArr)
}