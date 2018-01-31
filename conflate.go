package conflate

import (
	"net/url"
)

// Conflate contains a 'working' merged data set and optionally a JSON v4 schema
type Conflate struct {
	data   interface{}
	schema interface{}
}

// New constructs a new empty Conflate instance
func New() *Conflate {
	return &Conflate{}
}

// FromFiles constructs a new Conflate instance populated with the data from the given files
func FromFiles(paths ...string) (*Conflate, error) {
	c := New()
	err := c.AddFiles(paths...)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// FromURLs constructs a new Conflate instance populated with the data from the given URLs
func FromURLs(urls ...url.URL) (*Conflate, error) {
	c := New()
	err := c.AddURLs(urls...)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// FromData constructs a new Conflate instance populated with the given data
func FromData(data ...[]byte) (*Conflate, error) {
	c := New()
	err := c.AddData(data...)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// FromGo constructs a new Conflate instance populated with the given golang objects
func FromGo(data ...interface{}) (*Conflate, error) {
	c := New()
	err := c.AddGo(data...)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// AddFiles recursively merges the data from the given files into the Conflate instance
func (c *Conflate) AddFiles(paths ...string) error {
	urls, err := toURLs(nil, paths...)
	if err != nil {
		return err
	}
	return c.AddURLs(urls...)
}

// AddURLs recursively merges the data from the given urls into the Conflate instance
func (c *Conflate) AddURLs(urls ...url.URL) error {
	data, err := loadURLsRecursive(nil, urls...)
	if err != nil {
		return err
	}
	return c.MergeData(data...)
}

// AddGo recursively merges the given (json-serializable) golang objects into the Conflate instance
func (c *Conflate) AddGo(objs ...interface{}) error {
	data, err := marshalAll(jsonMarshal, objs...)
	if err != nil {
		return err
	}
	return c.AddData(data...)
}

// AddData recursively merges the given data into the Conflate instance
func (c *Conflate) AddData(data ...[]byte) error {
	data, err := loadDataRecursive(nil, data...)
	if err != nil {
		return err
	}
	return c.MergeData(data...)
}

// MergeData merges the given data into the Conflate instance
func (c *Conflate) MergeData(data ...[]byte) error {
	doms, err := unmarshalAll(unmarshalAny, data...)
	if err != nil {
		return err
	}
	err = mergeTo(&c.data, doms...)
	if err != nil {
		return err
	}
	c.removeIncludes()
	return nil
}

func (c *Conflate) removeIncludes() {
	m, ok := c.data.(map[string]interface{})
	if ok {
		delete(m, "includes")
	}
}

// SetSchemaFile loads a JSON v4 schema from the given path
func (c *Conflate) SetSchemaFile(path string) error {
	url, err := toURL(nil, path)
	if err != nil {
		return wrapError(err, "Failed to obtain url to schema file")
	}
	return c.SetSchemaURL(url)
}

// SetSchemaURL loads a JSON v4 schema from the given URL
func (c *Conflate) SetSchemaURL(url url.URL) error {
	data, err := loadURL(url)
	if err != nil {
		return wrapError(err, "Failed to load schema file")
	}
	return c.SetSchemaData(data)
}

// SetSchemaData loads a JSON v4 schema from the given data
func (c *Conflate) SetSchemaData(data []byte) error {
	var schema interface{}
	err := jsonUnmarshal(data, &schema)
	if err != nil {
		return wrapError(err, "Schema is not valid json")
	}
	err = validateSchema(schema)
	if err != nil {
		return wrapError(err, "The schema is not valid against the meta-schema http://json-schema.org/draft-04/schema")
	}
	c.schema = schema
	return nil
}

// ApplyDefaults sets any nil or missing values in the data, to the default values defined in the JSON v4 schema
func (c *Conflate) ApplyDefaults() error {
	if c.schema == nil {
		return makeError("Schema is not set")
	}
	err := applyDefaults(&c.data, c.schema)
	return wrapError(err, "The defaults could not be applied")
}

// Validate checks the data against the JSON v4 schema
func (c *Conflate) Validate() error {
	if c.schema == nil {
		return makeError("Schema is not set")
	}
	err := validate(&c.data, c.schema)
	return wrapError(err, "Schema validation failed")
}

// Unmarshal extracts the data as a Golang object
func (c *Conflate) Unmarshal(out interface{}) error {
	return jsonMarshalUnmarshal(c.data, out)
}

// MarshalJSON exports the data as JSON
func (c *Conflate) MarshalJSON() ([]byte, error) {
	return jsonMarshal(c.data)
}

// MarshalYAML exports the data as YAML
func (c *Conflate) MarshalYAML() ([]byte, error) {
	return yamlMarshal(c.data)
}

// MarshalTOML exports the data as TOML
func (c *Conflate) MarshalTOML() ([]byte, error) {
	return tomlMarshal(c.data)
}

// MarshalSchema exports the schema as JSON
func (c *Conflate) MarshalSchema() ([]byte, error) {
	return jsonMarshal(c.schema)
}
