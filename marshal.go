package conflate

import (
	"bytes"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/ghodss/yaml"
)

func marshalAll(fMarshal func(interface{}) ([]byte, error), data ...interface{}) ([][]byte, error) {
	var outs [][]byte
	for i, datum := range data {
		out, err := fMarshal(datum)
		if err != nil {
			return nil, wrapError(err, "Could not marshal data %v", i)
		}
		outs = append(outs, out)
	}
	return outs, nil
}

func unmarshalAll(fUnmarshal func([]byte, interface{}) error, data ...[]byte) ([]interface{}, error) {
	var outs []interface{}
	for i, datum := range data {
		var out interface{}
		err := fUnmarshal(datum, &out)
		if err != nil {
			return nil, wrapError(err, "Could not unmarshal data %v", i)
		}
		outs = append(outs, out)
	}
	return outs, nil
}

func unmarshalAny(data []byte, out interface{}) error {
	errs := makeError("Could not unmarshal data")

	err := jsonUnmarshal(data, out)
	if err == nil {
		return nil
	}
	errs = wrapError(err, errs.Error())

	err = tomlUnmarshal(data, out)
	if err == nil {
		return nil
	}
	errs = wrapError(err, errs.Error())

	err = yamlUnmarshal(data, out)
	if err == nil {
		return nil
	}
	errs = wrapError(err, errs.Error())

	return errs
}

func jsonMarshalUnmarshal(in interface{}, out interface{}) error {
	data, err := jsonMarshal(in)
	if err != nil {
		return err
	}
	return jsonUnmarshal(data, out)
}

func jsonUnmarshal(data []byte, out interface{}) error {
	err := json.Unmarshal(data, out)
	if err != nil {
		return wrapError(err, "The data could not be unmarshalled as json")
	}
	return nil
}

func yamlUnmarshal(data []byte, out interface{}) error {
	err := yaml.Unmarshal(data, out)
	if err != nil {
		return wrapError(err, "The data could not be unmarshalled as yaml")
	}
	return nil
}

func tomlUnmarshal(data []byte, out interface{}) error {
	err := toml.Unmarshal(data, out)
	if err != nil {
		return wrapError(err, "The data could not be unmarshalled as toml")
	}
	return nil
}

func jsonMarshal(data interface{}) ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(data)
	if err != nil {
		return nil, wrapError(err, "The data could not be marshalled to json")
	}
	return buffer.Bytes(), nil
}

func yamlMarshal(in interface{}) ([]byte, error) {
	data, err := yaml.Marshal(in)
	if err != nil {
		return nil, wrapError(err, "The data could not be marshalled to yaml")
	}
	return data, nil
}

func tomlMarshal(in interface{}) (out []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(makeError("%v", r), "The data could not be marshalled to toml")
		}
	}()
	buf := bytes.Buffer{}
	enc := toml.NewEncoder(&buf)
	err = enc.Encode(in)
	if err != nil {
		return nil, wrapError(err, "The data could not be marshalled to toml")
	}
	out = buf.Bytes()
	return out, nil
}
