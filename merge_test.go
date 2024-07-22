package conflate

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeTo(t *testing.T) {
	var toData interface{}

	data1 := 1
	data2 := 2
	data3 := 3
	err := mergeTo(&toData, data1, data2, data3)
	assert.Nil(t, err)
	assert.NotNil(t, toData)
	assert.Equal(t, toData, data3)
}

func TestMergeTo_MergeError(t *testing.T) {
	var toData interface{}

	fromData := make(map[string]interface{})
	err := mergeTo(toData, fromData)
	assert.NotNil(t, err)
}

func TestMerge(t *testing.T) {
	toData := testMergeGetData(t, testMergeData1)
	fromData := testMergeGetData(t, testMergeData2)
	merged := testMergeGetData(t, testMergeData1)
	err := merge(&merged, fromData)
	assert.Nil(t, err)
	testMergeCheck(t, merged, toData, fromData)
}

func TestMergeReversed(t *testing.T) {
	toData := testMergeGetData(t, testMergeData2)
	fromData := testMergeGetData(t, testMergeData1)
	merged := testMergeGetData(t, testMergeData2)
	err := merge(&merged, fromData)
	assert.Nil(t, err)
	testMergeCheck(t, merged, toData, fromData)
}

func TestMerge_SimpleString(t *testing.T) {
	toData := "x"
	fromData := "y"
	err := merge(&toData, fromData)
	assert.Nil(t, err)
	assert.Equal(t, "y", toData)
}

func TestMerge_SimpleInt(t *testing.T) {
	toData := 1
	fromData := 2
	err := merge(&toData, fromData)
	assert.Nil(t, err)
	assert.Equal(t, 2, toData)
}

func TestMerge_SimpleFloat(t *testing.T) {
	toData := 1.0
	fromData := 2.0
	err := merge(&toData, fromData)
	assert.Nil(t, err)
	assert.Equal(t, 2.0, toData)
}

func TestMerge_SimpleMap(t *testing.T) {
	toData := map[string]interface{}{"x": 1}
	fromData := map[string]interface{}{"x": 2, "y": 2}
	err := merge(&toData, fromData)
	assert.Nil(t, err)
	assert.Equal(t, 2, toData["x"])
	assert.Equal(t, 2, toData["y"])
}

func TestMerge_SimpleSlice(t *testing.T) {
	toData := []interface{}{1, 2, 3}
	fromData := []interface{}{4, 5, 6}
	err := merge(&toData, fromData)
	assert.Nil(t, err)
	assert.Equal(t, toData, []interface{}{1, 2, 3, 4, 5, 6})
}

func TestMerge_ToNil(t *testing.T) {
	fromData := make(map[string]interface{})
	err := merge(nil, fromData)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "must not be nil")
}

func TestMerge_ToNotPtr(t *testing.T) {
	fromData := make(map[string]interface{})
	toData := make(map[string]interface{})
	err := merge(toData, fromData)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "must be a pointer")
}

func TestMerge_FromNil(t *testing.T) {
	data := make(map[string]interface{})
	err := merge(&data, nil)
	assert.Nil(t, err)
	assert.Equal(t, data, data)
}

func TestMerge_ToValNil(t *testing.T) {
	fromData := make(map[string]interface{})

	var toData interface{}

	err := merge(&toData, fromData)
	assert.Nil(t, err)
	assert.Equal(t, toData, fromData)
}

func TestMerge_FromMapInvalid(t *testing.T) {
	fromData := make(map[int]int)
	toData := make(map[string]interface{})
	err := merge(&toData, fromData)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "source value must be a map[string]interface{}")
}

func TestMerge_ToMapInvalid(t *testing.T) {
	fromData := make(map[string]interface{})
	toData := make(map[int]int)
	err := merge(&toData, fromData)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "destination value must be a map[string]interface{}")
}

func TestMerge_FromSliceInvalid(t *testing.T) {
	fromData := make([]int, 0)
	toData := make([]interface{}, 0)
	err := merge(&toData, fromData)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "source value must be a []interface{}")
}

func TestMerge_ToSliceInvalid(t *testing.T) {
	fromData := make([]interface{}, 0)
	toData := make([]int, 0)
	err := merge(&toData, fromData)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "destination value must be a []interface{}")
}

func TestMerge_IntToSliceInvalid(t *testing.T) {
	fromData := 123
	toData := make([]int, 0)
	err := merge(&toData, fromData)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "the destination type ([]int) must be the same as the source type (int)")
}

func TestMerge_IntToMapInvalid(t *testing.T) {
	fromData := 123
	toData := make(map[string]int)
	err := merge(&toData, fromData)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "the destination type (map[string]int) must be the same as the source type (int)")
}

func TestMerge_BadPropertyMerge(t *testing.T) {
	toData := map[string]interface{}{"x": 1}
	fromData := map[string]interface{}{"x": map[string]string{}}
	err := merge(&toData, fromData)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to merge object property")
}

func TestMerge_Equal(t *testing.T) {
	toData := map[string]interface{}{"x": 1}
	fromData := map[string]interface{}{"x": 1}
	err := merge(&toData, fromData)
	assert.Nil(t, err)
	assert.Equal(t, toData, fromData)
}

func testMergeCheck(t *testing.T, merged, data1, data2 interface{}) {
	t.Helper()

	mergedVal := reflect.ValueOf(merged)

	//nolint:exhaustive // test caseload
	switch mergedVal.Kind() {
	case reflect.Map:
		testMergeCheckMap(t, merged, data1, data2)
	case reflect.Slice:
		testMergeCheckSlice(t, merged, data1, data2)
	default:
		switch {
		case data1 == nil && data2 != nil:
			assert.Equal(t, data2, merged)
		case data1 != nil && data2 == nil:
			assert.Equal(t, data1, merged)
		case data1 != nil && data2 != nil:
			assert.Equal(t, data2, merged)
		default:
			assert.Nil(t, merged)
		}
	}
}

func testMergeCheckMap(t *testing.T, merged, data1, data2 interface{}) {
	t.Helper()

	var (
		mergedMap, data1Map, data2Map map[string]interface{}
		ok                            bool
	)

	if merged != nil {
		mergedMap, ok = merged.(map[string]interface{})
		assert.True(t, ok)
	}

	if data1 != nil {
		data1Map, ok = data1.(map[string]interface{})
		assert.True(t, ok)
	}

	if data2 != nil {
		data2Map, ok = data2.(map[string]interface{})
		assert.True(t, ok)
	}

	for name, mergedItem := range mergedMap {
		var data1Item, data2Item interface{}

		if data1Map != nil {
			data1Item = data1Map[name]
		}

		if data2Map != nil {
			data2Item = data2Map[name]
		}

		testMergeCheck(t, mergedItem, data1Item, data2Item)
	}
}

func testMergeCheckSlice(t *testing.T, merged, data1, data2 interface{}) {
	t.Helper()

	var (
		mergedArr, data1Arr, data2Arr []interface{}
		ok                            bool
	)

	if merged != nil {
		mergedArr, ok = merged.([]interface{})
		assert.True(t, ok)
	}

	if data1 != nil {
		data1Arr, ok = data1.([]interface{})
		assert.True(t, ok)
	}

	if data2 != nil {
		data2Arr, ok = data2.([]interface{})
		assert.True(t, ok)
	}

	assert.Equal(t, len(data1Arr)+len(data2Arr), len(mergedArr))

	data1Pad := make([]interface{}, len(data2Arr))
	data2Pad := make([]interface{}, len(data1Arr))
	data1Arr = append(data1Arr, data1Pad...)
	data2Arr = append(data2Pad, data2Arr...)

	assert.Equal(t, len(data1Arr), len(mergedArr))
	assert.Equal(t, len(data2Arr), len(mergedArr))

	for i, mergedItem := range mergedArr {
		data1Item := data1Arr[i]
		data2Item := data2Arr[i]
		testMergeCheck(t, mergedItem, data1Item, data2Item)
	}
}

// ----------

func testMergeGetData(t *testing.T, data []byte) interface{} {
	t.Helper()

	var out interface{}

	err := json.Unmarshal(data, &out)
	assert.Nil(t, err)
	assert.NotNil(t, out)

	return out
}

var testMergeData1 = []byte(`
{
  "int_to_only": 1,
  "str_to_only": "str_to",
  "bool_to_only": true,

  "int_both": 1,
  "str_both": "str_to",
  "bool_both": true,

  "map_to_only" : {
    "int_to_only": 1,
    "str_to_only": "str_to",
    "bool_to_only": true,

    "int_both": 1,
    "str_both": "str_to",
    "bool_both": true,

    "array_to_only": [
      "str1_to",
      "str2_to",
      "str3_to"
    ],

    "array_both": [
      "str1_to",
      "str2_to",
      "str3_to"
    ]
  },

  "map_both" : {
    "int_to_only": 1,
    "str_to_only": "str_to",
    "bool_to_only": true,

    "int_both": 1,
    "str_both": "str_to",
    "bool_both": true
  },

  "array_to_only": [
    "str1_to",
    "str2_to",
    "str3_to"
  ],

  "array_both": [
    "str1_to",
    "str2_to",
    "str3_to"
  ]
}
`)

var testMergeData2 = []byte(`
{
  "int_from_only": 2,
  "str_from_only": "str_from",
  "bool_from_only": false,

  "int_both": 2,
  "str_both": "str_from",
  "bool_both": false,

  "map_from_only" : {
    "int_from_only": 2,
    "str_from_only": "str_from",
    "bool_from_only": false,

    "int_both": 2,
    "str_both": "str_from",
    "bool_both": false,

    "array_from_only": [
      "str1_from",
      "str2_from",
      "str3_from"
    ],

    "array_both": [
      "str1_from",
      "str2_from",
      "str3_from"
    ]
  },

  "map_both" : {
    "int_from_only": 2,
    "str_from_only": "str_from",
    "bool_from_only": false,

    "int_both": 2,
    "str_both": "str_from",
    "bool_both": false
  },

  "array_from_only": [
    "str1_from",
    "str2_from",
    "str3_from"
  ],

  "array_both": [
    "str1_from",
    "str2_from",
    "str3_from"
  ]
}
`)
