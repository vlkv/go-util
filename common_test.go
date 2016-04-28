package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFindIndexIntSlice(t *testing.T) {
	arr := []int{2, 4, 6, 8, 10, 12}

	assert.Equal(t, -1, FindIndex(len(arr), func(i int) bool { return arr[i] == 0 }))
	assert.Equal(t, 0, FindIndex(len(arr), func(i int) bool { return arr[i] == 2 }))
	assert.Equal(t, -1, FindIndex(len(arr), func(i int) bool { return arr[i] == 5 }))
	assert.Equal(t, 2, FindIndex(len(arr), func(i int) bool { return arr[i] == 6 }))
	assert.Equal(t, 5, FindIndex(len(arr), func(i int) bool { return arr[i] == 12 }))
	assert.Equal(t, -1, FindIndex(len(arr), func(i int) bool { return arr[i] == 13 }))
}

func TestFindIndexStringSlice(t *testing.T) {
	arr := []string{"Cindy", "Baz", "Karl", "Alla"}

	assert.Equal(t, -1, FindIndex(len(arr), func(i int) bool { return arr[i] == "Liza" }))
	assert.Equal(t, 0, FindIndex(len(arr), func(i int) bool { return arr[i] == "Cindy" }))
	assert.Equal(t, 3, FindIndex(len(arr), func(i int) bool { return arr[i] == "Alla" }))
}
