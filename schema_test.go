package conflate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateSchema(t *testing.T) {
	metaSchema = nil
	data := `{"title": "test"}`
	var schema interface{}
	err := jsonUnmarshal([]byte(data), &schema)
	assert.Nil(t, err)
	err = validateSchema(schema)
	assert.Nil(t, err)
	assert.NotNil(t, metaSchema)
}

func TestValidateSchema_AnyOf(t *testing.T) {
	data := `{ "type": "object", "properties": { "test": { "anyOf": [ { "type": "integer" } ] } } }`
	var schema interface{}
	err := jsonUnmarshal([]byte(data), &schema)
	assert.Nil(t, err)
	err = validateSchema(schema)
	assert.Nil(t, err)
}

func TestValidateSchema_Error(t *testing.T) {
	metaSchema = nil
	oldMetaSchemaData := metaSchemaData
	defer func() {
		metaSchemaData = oldMetaSchemaData
		metaSchema = nil
	}()
	metaSchemaData = []byte("invalid json")
	err := validateSchema("test")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Could not load json meta-schema")
}

func TestValidate(t *testing.T) {
	var data interface{}
	var schema interface{}
	err := jsonUnmarshal(testSchemaData, &data)
	assert.Nil(t, err)
	err = jsonUnmarshal(testSchema, &schema)
	assert.Nil(t, err)
	err = validate(data, schema)
	assert.Nil(t, err)
}

func TestValidate_ValidateError(t *testing.T) {
	var data interface{}
	var schema interface{}
	err := jsonUnmarshal(testSchemaData, &data)
	assert.Nil(t, err)
	err = validate(data, schema)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "An error occurred during validation")
	assert.Contains(t, err.Error(), "Invalid JSON")
}

func TestValidate_NotValid(t *testing.T) {
	var data map[string]interface{}
	var schema map[string]interface{}
	err := jsonUnmarshal(testSchemaData, &data)
	assert.Nil(t, err)
	err = jsonUnmarshal(testSchema, &schema)
	assert.Nil(t, err)
	err = jsonUnmarshal(testSchema, &schema)
	assert.Nil(t, err)
	obj := data["obj"].(map[string]interface{})
	obj["str"] = 123
	err = validate(data, schema)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "The document is not valid against the schema")
	assert.Contains(t, err.Error(), "Invalid type. Expected: string, given: integer")
	assert.Contains(t, err.Error(), "(#/obj/str)")
}

func TestValidate_CustomFormatError(t *testing.T) {
	var data interface{}
	var schema map[string]interface{}
	err := jsonUnmarshal(testSchemaData, &data)
	assert.Nil(t, err)
	err = jsonUnmarshal(testSchema, &schema)
	assert.Nil(t, err)
	props := schema["properties"].(map[string]interface{})
	str := props["str"].(map[string]interface{})
	str["format"] = "xml-template"
	err = validate(data, schema)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "The document is not valid against the schema")
	assert.Contains(t, err.Error(), "Does not match format")
	assert.Contains(t, err.Error(), "(#/str)")
}

// -----------

func TestApplyDefaults_DataNil(t *testing.T) {
	schema := map[string]interface{}{}
	err := applyDefaults(nil, schema)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Destination value must not be nil")
}

func TestApplyDefaults_DataNotPtr(t *testing.T) {
	var data = 1
	schema := map[string]interface{}{}
	err := applyDefaults(data, schema)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Destination value must be a pointer")
}

func TestApplyDefaults_SchemaNotMap(t *testing.T) {
	data := 1
	schema := map[string]interface{}{"anyOf": nil}
	err := applyDefaults(&data, schema)
	assert.Nil(t, err)
	assert.Equal(t, 1, data)
}

func TestApplyDefaults_SchemaNoType(t *testing.T) {
	data := 1
	schema := map[string]interface{}{}
	err := applyDefaults(&data, schema)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Schema section does not have a valid 'type' attribute")
	assert.Equal(t, 1, data)
}

// --------

func TestApplyDefaults_NodeNotObject(t *testing.T) {
	data := 1
	schema := map[string]interface{}{"type": "object"}
	err := applyDefaults(&data, schema)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Node should be an 'object'")
}

func TestApplyDefaults_ObjectDefault(t *testing.T) {
	var data interface{}
	schema := map[string]interface{}{
		"type": "object",
		"default": map[string]interface{}{
			"val": 1,
		},
	}
	err := applyDefaults(&data, schema)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"val": 1}, data)
}

func TestApplyDefaults_ObjectDefaultWithDefaultProperty(t *testing.T) {
	var data interface{}
	schema := map[string]interface{}{
		"type":    "object",
		"default": map[string]interface{}{},
		"properties": map[string]interface{}{
			"val": map[string]interface{}{
				"type":    "integer",
				"default": 1,
			},
		},
	}
	err := applyDefaults(&data, schema)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"val": 1}, data)
}

func TestApplyDefaults_ObjectPropertyDefault(t *testing.T) {
	var data interface{}
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"val": map[string]interface{}{
				"type":    "integer",
				"default": 1,
			},
		},
	}
	err := applyDefaults(&data, schema)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"val": 1}, data)
}

func TestApplyDefaults_ObjectPropertyFailed(t *testing.T) {
	var data interface{}
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"val": nil,
		},
	}
	err := applyDefaults(&data, schema)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Failed to apply defaults to object property")
	assert.Contains(t, err.Error(), "Schema section is not a map (#/val)")
}

func TestApplyDefaults_ObjectAdditionalPropertyDefault(t *testing.T) {
	data := map[string]interface{}{"val": nil}
	schema := map[string]interface{}{
		"type": "object",
		"additionalProperties": map[string]interface{}{
			"type":    "integer",
			"default": 1,
		},
	}
	err := applyDefaults(&data, schema)
	assert.Nil(t, err)
	assert.Equal(t, 1, data["val"])
}

func TestApplyDefaults_ObjectAdditionalPropertyFailed(t *testing.T) {
	data := map[string]interface{}{"val": 1}
	schema := map[string]interface{}{
		"type":                 "object",
		"additionalProperties": map[string]interface{}{},
	}
	err := applyDefaults(&data, schema)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Failed to apply defaults to additional object property")
	assert.Contains(t, err.Error(), "Schema section does not have a valid 'type' attribute")
}

func TestApplyDefaults_ObjectAdditionalPropertyBool(t *testing.T) {
	data := map[string]interface{}{"val": 1}
	schema := map[string]interface{}{
		"type":                 "object",
		"additionalProperties": false,
	}
	err := applyDefaults(&data, schema)
	assert.Nil(t, err)
}

// --------

func TestApplyDefaults_NodeNotSlice(t *testing.T) {
	data := 1
	schema := map[string]interface{}{"type": "array"}
	err := applyDefaults(&data, schema)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Node should be an 'array'", err.Error())
}

func TestApplyDefaults_SliceDefault(t *testing.T) {
	var data interface{}
	schema := map[string]interface{}{
		"type": "array",
		"default": []interface{}{
			1,
		},
	}
	err := applyDefaults(&data, schema)
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{1}, data)
}

func TestApplyDefaults_SliceDefaultWithElementDefault(t *testing.T) {
	var data interface{}
	schema := map[string]interface{}{
		"type": "array",
		"default": []interface{}{
			nil,
		},
		"items": map[string]interface{}{
			"type":    "integer",
			"default": 1,
		},
	}
	err := applyDefaults(&data, schema)
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{1}, data)
}

func TestApplyDefaults_SliceElementDefault(t *testing.T) {
	data := []interface{}{nil}
	schema := map[string]interface{}{
		"type": "array",
		"items": map[string]interface{}{
			"type":    "integer",
			"default": 1,
		},
	}
	err := applyDefaults(&data, schema)
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{1}, data)
}

func TestApplyDefaults_SliceFailed(t *testing.T) {
	data := []interface{}{1}
	schema := map[string]interface{}{
		"type":  "array",
		"items": map[string]interface{}{},
	}
	err := applyDefaults(&data, schema)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Failed to apply defaults to array item")
	assert.Contains(t, err.Error(), "Schema section does not have a valid 'type' attribute (#[0])")
}

func TestApplyDefaults_Empty(t *testing.T) {
	var data interface{}
	var defaults interface{}
	var schema interface{}
	err := jsonUnmarshal(testSchemaDefaults, &defaults)
	assert.Nil(t, err)
	err = jsonUnmarshal(testSchema, &schema)
	assert.Nil(t, err)
	err = applyDefaults(&data, schema)
	assert.Nil(t, err)
	assert.Equal(t, defaults, data)
}

func TestApplyDefaults_NoDefaults(t *testing.T) {
	var data interface{}
	var dataExpected interface{}
	var schema interface{}
	err := jsonUnmarshal(testSchemaData, &data)
	assert.Nil(t, err)
	err = jsonUnmarshal(testSchemaData, &dataExpected)
	assert.Nil(t, err)
	err = jsonUnmarshal(testSchema, &schema)
	assert.Nil(t, err)
	err = applyDefaults(&data, schema)
	assert.Nil(t, err)
	assert.Equal(t, dataExpected, data)
}

func TestApplyDefaults_MissingIntFields(t *testing.T) {
	var data map[string]interface{}
	var schema interface{}
	err := jsonUnmarshal(testSchemaData, &data)
	assert.Nil(t, err)
	err = jsonUnmarshal(testSchema, &schema)
	assert.Nil(t, err)

	delete(data, "int")
	delete(data, "array_of_int")
	obj := data["obj"].(map[string]interface{})
	delete(obj, "int")
	arr := data["array_of_obj"].([]interface{})
	arrObj := arr[0].(map[string]interface{})
	delete(arrObj, "int")

	err = applyDefaults(&data, schema)
	assert.Nil(t, err)

	assert.Equal(t, 1.0, data["int"])
	assert.Equal(t, []interface{}{1.0}, data["array_of_int"])
	assert.Equal(t, 1.0, obj["int"])
	assert.Equal(t, 1.0, arrObj["int"])
}

// -----------

var testSchemaData = []byte(`
{
  "int": 0,
  "str": "str",
  "bool": false,
  "obj": {
    "int": 0,
    "str": "str",
    "bool": false,
		"array_of_int": [ 0 ],
		"array_of_str": [ "str" ],
		"array_of_bool": [ false ]
  },
  "array_of_int": [ 0 ],
  "array_of_str": [ "str" ],
  "array_of_bool": [ false ],
  "array_of_obj": [
		{
			"int": 0,
			"str": "str",
			"bool": false
  	}
	]
}`)

var testSchemaDefaults = []byte(`
{
  "int": 1,
  "str": "test",
  "bool": true,
  "obj": {
		"int": 1,
		"str": "test",
		"bool": true,
		"array_of_int": [ 1 ],
		"array_of_str": [ "test" ],
		"array_of_bool": [ true ]
  },
  "array_of_int": [ 1 ],
  "array_of_str": [ "test" ],
  "array_of_bool": [ true ],
  "array_of_obj": [
		{
		"int": 1,
		"str": "test",
		"bool": true
  	}
	]
}`)

var testSchema = []byte(`
{
  "title": "test",
  "type": "object",
  "properties": {
    "int": { "type": "integer", "default": 1 },
    "str": { "type": "string", "default": "test" },
    "bool": { "type": "boolean", "default": true },
    "obj": {
      "type": "object",
      "properties": {
        "int": { "type": "integer", "default": 1 },
        "str": { "type": "string", "default": "test" },
        "bool": { "type": "boolean", "default": true },
				"array_of_int": {
					"type": "array",
					"items": { "type": "integer" },
					"default": [ 1 ]
				},
				"array_of_str": {
					"type": "array",
					"items": { "type": "string" },
					"default": [ "test" ]
				},
				"array_of_bool": {
					"type": "array",
					"items": { "type": "boolean" },
					"default": [ true ]
				}
			}
    }, 
    "array_of_int": {
			"type": "array",
      "items": { "type": "integer" },
      "default": [ 1 ]
    },
    "array_of_str": {
			"type": "array",
      "items": { "type": "string" },
      "default": [ "test" ]
    },
    "array_of_bool": {
			"type": "array",
      "items": { "type": "boolean" },
      "default": [ true ]
    },
    "array_of_obj": {
			"type": "array",
      "items": {
        "type": "object",
        "properties": {
          "int": { "type": "integer", "default": 1 },
          "str": { "type": "string", "default": "test" },
          "bool": { "type": "boolean", "default": true }
        }
      },
			"default" : [
				{
				"int": 1,
				"str": "test",
				"bool": true
				}
			]
    }
  }
}
`)

var testSchemaData1 = []byte(`
{
  "int": 0,
  "str": "str",
  "bool": false,
	"obj": {
		"str": "str"
	}
}`)

var testSchema1 = []byte(`
{
  "title": "test",
  "type": "object",
  "properties": {
    "int": { "type": "integer", "default": 1 },
    "str": { "type": "string", "default": "test" },
    "bool": { "type": "boolean", "default": true },
    "obj": {
			"type": "object",
			"properties": {
				"str": { "type": "string", "default": "test" }
			}
		}
  }
}
`)
