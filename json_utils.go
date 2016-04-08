package util

import (
	"encoding/json"
	"net/http"
	"fmt"
)

// TODO: Rename to JsonParse
func ParseJson(jsonStr string) (jsonObj interface{}) {
	err := json.Unmarshal([]byte(jsonStr), &jsonObj)
	if err != nil {
		panic(CreateHttpError(http.StatusInternalServerError, fmt.Sprintf("Cannot decode json '%v', reason %v", jsonStr, err)))
	}
	return
}

func JsonEncode(jsonObj interface{}) (jsonStr []byte) {
	jsonStr, err := json.Marshal(jsonObj)
	if err != nil {
		panic(CreateHttpError(http.StatusInternalServerError, fmt.Sprintf("Cannot encode json %v, reason %v", jsonObj, err)))
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

func JsonPut(jsonObj interface{}, val interface{}, keys ...string) {
	curr := jsonObj.(map[string]interface{})
	for i := 0; i < len(keys)-1; i++ {
		key := keys[i]
		curr = curr[key].(map[string]interface{})
	}
	lastKey := keys[len(keys)-1]
	curr[lastKey] = val
}