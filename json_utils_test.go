package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestJsonEncode(t *testing.T) {
	parent := make(map[string]interface{})

	parent["integer"] = 10
	parent["float"] = 3.14
	parent["string"] = "some text"

	child := make(map[string]interface{})
	child["bool"] = true
	parent["child"] = child

	jsonStr := JsonEncode(parent)
	assert.Equal(t, "{\"child\":{\"bool\":true},\"float\":3.14,\"integer\":10,\"string\":\"some text\"}", string(jsonStr))
}

func TestJsonParseAndJsonGet(t *testing.T) {
	JsonStr := "{\"child\":{\"bool\":true},\"float\":3.14,\"integer\":10,\"string\":\"some text\"}"
	parent := ParseJson(JsonStr).(map[string]interface{})
	assert.Equal(t, 4, len(parent))
	assert.Equal(t, 10.0, JsonGet(parent, "integer").(float64)) // NOTE: It's not int after parse!
	assert.Equal(t, 3.14, JsonGet(parent, "float").(float64))
	assert.Equal(t, "some text", JsonGet(parent, "string").(string))

	child := JsonGet(parent, "child").(map[string]interface{})
	assert.Equal(t, 1, len(child))
	assert.Equal(t, true, JsonGet(child, "bool").(bool))
	assert.Equal(t, true, JsonGet(parent, "child", "bool").(bool)) // It's the same
}


func TestJsonPut(t *testing.T) {
	parent := make(map[string]interface{})
	child := make(map[string]interface{})
	child["existing_field"] = "old_val"
	parent["child"] = child

	JsonPut(parent, "new_val", "child", "existing_field")
	assert.Equal(t, "new_val", JsonGet(parent, "child", "existing_field"))

	JsonPut(parent, "new_val2", "child", "non_existing_field")
	assert.Equal(t, "new_val2", JsonGet(parent, "child", "non_existing_field"))
}

