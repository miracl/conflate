package conflate

import (
	"bytes"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/ghodss/yaml"
)

func unmarshalAll(data ...[]byte) ([]interface{}, error) {
	var outs []interface{}
	for i, datum := range data {
		var out interface{}
		err := unmarshal(datum, &out)
		if err != nil {
			return nil, wrapError(err, "Could not unmarshal data %v", i)
		}
		outs = append(outs, out)
	}
	return outs, nil
}

func unmarshal(data []byte, out interface{}) error {
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
	err = jsonUnmarshal(data, out)
	if err != nil {
		return err
	}
	return nil
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

func jsonMarshal(in interface{}) ([]byte, error) {
	data, err := json.Marshal(in)
	if err != nil {
		return nil, wrapError(err, "The data could not be marshalled to json")
	}
	return data, nil
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
