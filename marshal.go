package conflate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/ghodss/yaml"
)

var errToml = errors.New("the data could not be marshalled to toml")

func jsonMarshalAll(data ...interface{}) ([][]byte, error) {
	var outs [][]byte

	for i, datum := range data {
		out, err := jsonMarshal(datum)
		if err != nil {
			return nil, fmt.Errorf("could not marshal data %v: %w", i, err)
		}

		outs = append(outs, out)
	}

	return outs, nil
}

func jsonMarshalUnmarshal(in, out interface{}) error {
	data, err := jsonMarshal(in)
	if err != nil {
		return err
	}

	return JSONUnmarshal(data, out)
}

// JSONUnmarshal unmarshals the data as JSON.
func JSONUnmarshal(data []byte, out interface{}) error {
	err := json.Unmarshal(data, out)
	if err != nil {
		return fmt.Errorf("the data could not be unmarshalled as json: %w", err)
	}

	return nil
}

// YAMLUnmarshal unmarshals the data as YAML.
func YAMLUnmarshal(data []byte, out interface{}) error {
	err := yaml.Unmarshal(data, out)
	if err != nil {
		return fmt.Errorf("the data could not be unmarshalled as yaml: %w", err)
	}

	return nil
}

// TOMLUnmarshal unmarshals the data as TOML.
func TOMLUnmarshal(data []byte, out interface{}) error {
	err := toml.Unmarshal(data, out)
	if err != nil {
		return fmt.Errorf("the data could not be unmarshalled as toml: %w", err)
	}

	return nil
}

func jsonMarshal(data interface{}) ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	err := encoder.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("the data could not be marshalled to json: %w", err)
	}

	return buffer.Bytes(), nil
}

func yamlMarshal(in interface{}) ([]byte, error) {
	data, err := yaml.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("the data could not be marshalled to yaml: %w", err)
	}

	return data, nil
}

func tomlMarshal(in interface{}) (out []byte, err error) {
	defer func() {
		if isPanicking := recover(); isPanicking != nil {
			err = fmt.Errorf("%w : %v", errToml, isPanicking)
		}
	}()

	buf := bytes.Buffer{}
	enc := toml.NewEncoder(&buf)

	err = enc.Encode(in)
	if err != nil {
		return nil, fmt.Errorf("the data could not be marshalled to toml: %w", err)
	}

	out = buf.Bytes()

	return out, nil
}
