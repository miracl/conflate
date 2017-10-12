package conflate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalAll(t *testing.T) {
	data, err := unmarshalAll(testMarshalJSON, testMarshalYAML, testMarshalTOML)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(data))
	assert.Equal(t, testMarshalData, data[0])
	assert.Equal(t, testMarshalData, data[1])
	assert.Equal(t, testMarshalData, data[2])
}

func TestUnmarshalAll_Error(t *testing.T) {
	data, err := unmarshalAll(testMarshalJSON, testMarshalInvalid, testMarshalTOML)
	assert.NotNil(t, err)
	assert.Nil(t, data)
}

// --------

func TestUnmarshal_Json(t *testing.T) {
	var out interface{}
	err := unmarshal(testMarshalJSON, &out)
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, out)
}

func TestUnmarshal_Yaml(t *testing.T) {
	var out interface{}
	err := unmarshal(testMarshalYAML, &out)
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, out)
}

func TestUnmarshal_Toml(t *testing.T) {
	var out interface{}
	err := unmarshal(testMarshalTOML, &out)
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, out)
}

func TestUnmarshal_Unsupported(t *testing.T) {
	var out interface{}
	err := unmarshal(testMarshalInvalid, &out)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Could not unmarshal data")
	assert.Contains(t, err.Error(), "could not be unmarshalled as json")
	assert.Contains(t, err.Error(), "could not be unmarshalled as yaml")
	assert.Contains(t, err.Error(), "could not be unmarshalled as toml")
}

// --------

func TestJsonMarshalUnmarshal(t *testing.T) {
	var out interface{}
	err := jsonMarshalUnmarshal(testMarshalData, &out)
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, out)
}

func TestJsonMarshalUnmarshal_MarshalError(t *testing.T) {
	var out interface{}
	err := jsonMarshalUnmarshal(testMarshalDataInvalid, out)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not be marshalled to json")
}

func TestJsonMarshalUnmarshal_UnmarshalError(t *testing.T) {
	err := jsonMarshalUnmarshal(testMarshalData, testMarshalDataInvalid)
	assert.NotNil(t, err)
	t.Log(err)
	assert.Contains(t, err.Error(), "could not be unmarshalled as json")
}

// --------

func TestJsonUnmarshal(t *testing.T) {
	var out interface{}
	err := jsonUnmarshal(testMarshalJSON, &out)
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, out)
}

func TestJsonUnmarshal_Error(t *testing.T) {
	var out interface{}
	err := jsonUnmarshal(testMarshalInvalid, &out)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not be unmarshalled as json")
}

func TestYamlUnmarshal(t *testing.T) {
	var out interface{}
	err := yamlUnmarshal(testMarshalYAML, &out)
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, out)
}

func TestYamlUnmarshal_Error(t *testing.T) {
	var out interface{}
	err := yamlUnmarshal(testMarshalInvalid, &out)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not be unmarshalled as yaml")
}

func TestTomlUnmarshal(t *testing.T) {
	var out interface{}
	err := tomlUnmarshal(testMarshalTOML, &out)
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, out)
}

func TestTomlUnmarshal_Error(t *testing.T) {
	var out interface{}
	err := tomlUnmarshal(testMarshalInvalid, &out)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not be unmarshalled as toml")
}

// --------

func TestJsonMarshal(t *testing.T) {
	out, err := jsonMarshal(testMarshalData)
	assert.Nil(t, err)
	assert.Equal(t, string(testMarshalJSON), string(out))
}

func TestJsonMarshal_Error(t *testing.T) {
	out, err := jsonMarshal(testMarshalDataInvalid)
	assert.NotNil(t, err)
	assert.Nil(t, out)

	assert.Contains(t, err.Error(), "marshalled to json")
}

func TestYamlMarshal(t *testing.T) {
	out, err := yamlMarshal(testMarshalData)
	assert.Nil(t, err)
	assert.Equal(t, string(testMarshalYAML), string(out))
}

func TestYamlMarshal_Error(t *testing.T) {
	out, err := yamlMarshal(testMarshalDataInvalid)
	assert.NotNil(t, err)
	assert.Nil(t, out)
	assert.Contains(t, err.Error(), "marshalled to yaml")
}

func TestTomlMarshal(t *testing.T) {
	out, err := tomlMarshal(testMarshalData)
	assert.Nil(t, err)
	assert.Equal(t, string(testMarshalTOML), string(out))
}

func TestTomlMarshal_PanicError(t *testing.T) {
	out, err := tomlMarshal(testMarshalDataInvalid)
	assert.NotNil(t, err)
	assert.Nil(t, out)
	assert.Contains(t, err.Error(), "marshalled to toml")
}

func TestTomlMarshal_Error(t *testing.T) {
	in := []interface{}{123, "123"}
	out, err := tomlMarshal(in)
	assert.NotNil(t, err)
	assert.Nil(t, out)
	assert.Contains(t, err.Error(), "marshalled to toml")
}

// --------

var (
	testMarshalData        = map[string]interface{}{"key": "value"}
	testMarshalDataInvalid = func() {}
	testMarshalJSON        = []byte(`{"key":"value"}`)
	testMarshalYAML        = []byte(`key: value
`)
	testMarshalTOML = []byte(`key = "value"
`)
	testMarshalInvalid = []byte(`{invalid`)
)
