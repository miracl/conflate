package conflate

import (
	"net/url"
)

// Includes is used to specify the top level key that holds the includes array
var Includes = "includes"

// Conflate contains a 'working' merged data set and optionally a JSON v4 schema
type Conflate struct {
	data   interface{}
	schema interface{}
	loader loader
}

// New constructs a new empty Conflate instance
func New() *Conflate {
	return &Conflate{
		loader: loader{
			newFiledata: newFiledata,
		},
	}
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

// Expand is an option to automatically expand environment variables in data files
func (c *Conflate) Expand(expand bool) {
	if expand {
		c.loader.newFiledata = newExpandedFiledata
	} else {
		c.loader.newFiledata = newFiledata
	}
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
	data, err := c.loader.loadURLsRecursive(nil, urls...)
	if err != nil {
		return err
	}
	return c.mergeData(data...)
}

// AddGo recursively merges the given (json-serializable) golang objects into the Conflate instance
func (c *Conflate) AddGo(objs ...interface{}) error {
	data, err := jsonMarshalAll(objs...)
	if err != nil {
		return err
	}
	return c.AddData(data...)
}

// AddData recursively merges the given data into the Conflate instance
func (c *Conflate) AddData(data ...[]byte) error {
	fdata, err := c.loader.wrapFiledatas(data...)
	if err != nil {
		return err
	}
	return c.addData(fdata...)
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
	err := JSONUnmarshal(data, &schema)
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

func (c *Conflate) addData(fdata ...filedata) error {
	fdata, err := c.loader.loadDataRecursive(nil, fdata...)
	if err != nil {
		return err
	}
	return c.mergeData(fdata...)
}

func (c *Conflate) mergeData(fdata ...filedata) error {
	doms := filedatas(fdata).objs()
	return mergeTo(&c.data, doms...)
}
