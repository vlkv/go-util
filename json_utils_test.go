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
	parent := JsonParse(JsonStr).(map[string]interface{})
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

func TestJsonMerge(t *testing.T) {
	dst := map[string]interface{} {
		"name": "Bob",
		"salary": 1000000.0,
		"children": [...]string{"Mary", "Kevin", "Nancy"},
		"subObj": map[string]interface{} {
			"a": "dstA",
			"b": "dstB",
		},
	}

	src := map[string]interface{}{
		"name": "Alice",
		"children": [...]int{1, 2, 3},
		"subObj": map[string]interface{} {
			"b": "srcB",
			"c": "srcC",
		},
	}

	JsonMerge(dst, src)

	dstStr := JsonEncode(dst)
	assert.Equal(t, "{\"children\":[1,2,3],\"name\":\"Alice\",\"salary\":1e+06,\"subObj\":{\"b\":\"srcB\",\"c\":\"srcC\"}}", string(dstStr))
}

func TestJsonGetWithDefault(t *testing.T) {
	obj := map[string]interface{} {
		"name": "Alice",
		"children": [...]int{1, 2, 3},
		"subObj": map[string]interface{} {
			"b": "srcB",
			"c": "srcC",
			"reallyDeepObj": map[string]interface{} {
				"1": "one",
				"3": "three",
			},
		},
	}

	assert.Equal(t, "srcC", JsonGetWithDefault(obj, "defVal", "subObj", "c"))
	assert.Equal(t, "defVal", JsonGetWithDefault(obj, "defVal", "subObj", "a"))
	assert.Equal(t, "defVal", JsonGetWithDefault(obj, "defVal", "subObj", "nonExistentLevel", "nonExistentDeepObj"))
	assert.Equal(t, "one", JsonGetWithDefault(obj, "defVal", "subObj", "reallyDeepObj", "1"))
	assert.Equal(t, "defVal", JsonGetWithDefault(obj, "defVal", "subObj", "reallyDeepObj", "2"))
}

func TestMustGetJsonTagName(t *testing.T) {
	type TagCheck struct {
		Num int `json:"number"`
		Str string `json:"name,-"`
		NoTag string
	}
	tagCheck := TagCheck{1, "2", "3"}
	assert.Equal(t, "number", MustGetJsonTagName(tagCheck, "Num"))
	assert.Equal(t, "name", MustGetJsonTagName(tagCheck, "Str"))
	assert.Equal(t, "NoTag", MustGetJsonTagName(tagCheck, "NoTag"))
}
