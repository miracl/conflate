package conflate

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// --------

func TestMarshalAll(t *testing.T) {
	mockMarshal := func(obj interface{}) ([]byte, error) {
		return []byte(obj.(string)), nil
	}
	data, err := marshalAll(mockMarshal, "a", "b", "c")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(data))
	assert.Equal(t, []byte("a"), data[0])
	assert.Equal(t, []byte("b"), data[1])
	assert.Equal(t, []byte("c"), data[2])
}

func TestMarshalAll_Error(t *testing.T) {
	mockMarshal := func(obj interface{}) ([]byte, error) {
		return nil, errors.New("my error")
	}
	data, err := marshalAll(mockMarshal, "a")
	assert.NotNil(t, err)
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "my error")
}

// --------

func TestUnmarshalAll(t *testing.T) {
	mockUnmarshal := func(data []byte, obj interface{}) error {
		reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(string(data)))
		return nil
	}
	data, err := unmarshalAll(mockUnmarshal, []byte("1"), []byte("2"), []byte("3"))
	assert.Nil(t, err)
	assert.Equal(t, 3, len(data))
	assert.Equal(t, "1", data[0])
	assert.Equal(t, "2", data[1])
	assert.Equal(t, "3", data[2])
}

func TestUnmarshalAll_Error(t *testing.T) {
	mockUnmarshal := func(data []byte, obj interface{}) error {
		return errors.New("my error")
	}
	data, err := unmarshalAll(mockUnmarshal, []byte("1"))
	assert.NotNil(t, err)
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "my error")
}

// --------

func TestUnmarshal_Json(t *testing.T) {
	var out interface{}
	err := unmarshalAny(testMarshalJSON, &out)
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, out)
}

func TestUnmarshal_Yaml(t *testing.T) {
	var out interface{}
	err := unmarshalAny(testMarshalYAML, &out)
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, out)
}

func TestUnmarshal_Toml(t *testing.T) {
	var out interface{}
	err := unmarshalAny(testMarshalTOML, &out)
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, out)
}

func TestUnmarshal_Unsupported(t *testing.T) {
	var out interface{}
	err := unmarshalAny(testMarshalInvalid, &out)
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
