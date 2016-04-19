package util

import (
	"github.com/spf13/cast"
	"errors"
	"fmt"
)

func ToBool(s interface{}) bool {
	b, e := cast.ToBoolE(s)
	if e != nil {
		panic(errors.New(fmt.Sprintf("Could not parse bool '%s', reason: %v", s, e)))
	}
	return b
}

func ToBoolE(s interface{}) (bool, error) {
	return cast.ToBoolE(s)
}
