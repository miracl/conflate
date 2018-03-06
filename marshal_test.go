package conflate

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

// --------

func TestJSONMarshalAll(t *testing.T) {
	data, err := jsonMarshalAll("a", "b", "c")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(data))
	assert.Equal(t, []byte("\"a\"\n"), data[0])
	assert.Equal(t, []byte("\"b\"\n"), data[1])
	assert.Equal(t, []byte("\"c\"\n"), data[2])
}

func TestJSONMarshalAll_Error(t *testing.T) {
	mockMarshal := func(obj interface{}) ([]byte, error) {
		return nil, errors.New("my error")
	}
	data, err := jsonMarshalAll(mockMarshal, "a")
	assert.NotNil(t, err)
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "The data could not be marshalled")
}

// --------

func TestJSONMarshalUnmarshal(t *testing.T) {
	var out interface{}
	err := jsonMarshalUnmarshal(testMarshalData, &out)
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, out)
}

func TestJSONMarshalUnmarshal_MarshalError(t *testing.T) {
	var out interface{}
	err := jsonMarshalUnmarshal(testMarshalDataInvalid, out)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not be marshalled to json")
}

func TestJSONMarshalUnmarshal_UnmarshalError(t *testing.T) {
	err := jsonMarshalUnmarshal(testMarshalData, testMarshalDataInvalid)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not be unmarshalled as json")
}

// --------

func TestJSONUnmarshal(t *testing.T) {
	var out interface{}
	err := JSONUnmarshal(testMarshalJSON, &out)
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, out)
}

func TestJSONUnmarshal_Error(t *testing.T) {
	var out interface{}
	err := JSONUnmarshal(testMarshalInvalid, &out)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not be unmarshalled as json")
}

func TestYAMLUnmarshal(t *testing.T) {
	var out interface{}
	err := YAMLUnmarshal(testMarshalYAML, &out)
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, out)
}

func TestYAMLUnmarshal_Error(t *testing.T) {
	var out interface{}
	err := YAMLUnmarshal(testMarshalInvalid, &out)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not be unmarshalled as yaml")
}

func TestTOMLUnmarshal(t *testing.T) {
	var out interface{}
	err := TOMLUnmarshal(testMarshalTOML, &out)
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, out)
}

func TestTOMLUnmarshal_Error(t *testing.T) {
	var out interface{}
	err := TOMLUnmarshal(testMarshalInvalid, &out)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not be unmarshalled as toml")
}

// --------

func TestJSONMarshal(t *testing.T) {
	out, err := jsonMarshal(testMarshalData)
	assert.Nil(t, err)
	assert.Equal(t, string(testMarshalJSON), string(out))
}

func TestJSONMarshal_Error(t *testing.T) {
	out, err := jsonMarshal(testMarshalDataInvalid)
	assert.NotNil(t, err)
	assert.Nil(t, out)

	assert.Contains(t, err.Error(), "marshalled to json")
}

func TestYAMLMarshal(t *testing.T) {
	out, err := yamlMarshal(testMarshalData)
	assert.Nil(t, err)
	assert.Equal(t, string(testMarshalYAML), string(out))
}

func TestYAMLMarshal_Error(t *testing.T) {
	out, err := yamlMarshal(testMarshalDataInvalid)
	assert.NotNil(t, err)
	assert.Nil(t, out)
	assert.Contains(t, err.Error(), "marshalled to yaml")
}

func TestTOMLMarshal(t *testing.T) {
	out, err := tomlMarshal(testMarshalData)
	assert.Nil(t, err)
	assert.Equal(t, string(testMarshalTOML), string(out))
}

func TestTOMLMarshal_PanicError(t *testing.T) {
	out, err := tomlMarshal(testMarshalDataInvalid)
	assert.NotNil(t, err)
	assert.Nil(t, out)
	assert.Contains(t, err.Error(), "marshalled to toml")
}

func TestTOMLMarshal_Error(t *testing.T) {
	in := []interface{}{123, "123"}
	out, err := tomlMarshal(in)
	assert.NotNil(t, err)
	assert.Nil(t, out)
	assert.Contains(t, err.Error(), "marshalled to toml")
}

// --------

var (
	testValue              = `value!£$%^&*()_+-={}[]:@~;'#<>?,./|`
	testMarshalData        = map[string]interface{}{"key": testValue}
	testMarshalDataInvalid = func() {}
	testMarshalJSON        = []byte(`{
  "key": "` + testValue + `"
}
`)
	testMarshalYAML = []byte(`key: ` + testValue + `
`)
	testMarshalTOML = []byte(`key = "` + testValue + `"
`)
	testMarshalInvalid = []byte(`{invalid`)
)
