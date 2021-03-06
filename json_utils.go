package util

import (
	"encoding/json"
	"fmt"
	"errors"
	"github.com/imdario/mergo"
	"reflect"
	"strings"
)

func JsonParse(jsonStr string) (jsonObj interface{}) {
	err := json.Unmarshal([]byte(jsonStr), &jsonObj)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Cannot decode json '%v', reason %v", jsonStr, err)))
	}
	return
}

func JsonEncode(jsonObj interface{}) (jsonStr []byte) {
	jsonStr, err := json.Marshal(jsonObj)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Cannot encode json %v, reason %v", jsonObj, err)))
	}
	return
}

func JsonGet(jsonObj interface{}, keys ...string) interface{} {
	curr := jsonObj.(map[string]interface{})
	for i := 0; i < len(keys)-1; i++ {
		key := keys[i]
		curr = curr[key].(map[string]interface{})
	}
	lastKey := keys[len(keys)-1]
	return curr[lastKey]
}

func JsonGetWithDefault(jsonObj interface{}, defaultValue interface{}, keys ...string) (result interface{}) {
	defer func() {
		if err := recover(); err != nil {
			result = defaultValue
		}
	}()

	result = JsonGet(jsonObj, keys...)
	if result == nil {
		result = defaultValue
	}
	return
}

func JsonPut(jsonObj interface{}, val interface{}, keys ...string) {
	curr := jsonObj.(map[string]interface{})
	for i := 0; i < len(keys)-1; i++ {
		key := keys[i]
		curr = curr[key].(map[string]interface{})
	}
	lastKey := keys[len(keys)-1]
	curr[lastKey] = val
}

func JsonMerge(jsonObjDst map[string]interface{}, jsonObjSrc map[string]interface{}) {
	err := mergo.MergeWithOverwrite(&jsonObjDst, jsonObjSrc)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Could not merge json objs, reason: %v", err)))
	}
}


func parseJsonTag(tag string) (tagOptions string) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag
}

func MustGetJsonTagName(obj interface{}, fieldName string) (tagName string) {
	t := reflect.TypeOf(obj)
	field, found := t.FieldByName(fieldName); if !found {
		panic(fmt.Sprintf("Not found json tag in field %s", fieldName))
	}
	tagName = parseJsonTag(field.Tag.Get("json"))
	if tagName == "" {
		tagName = fieldName
	}
	return tagName
}
