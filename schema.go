package conflate

import (
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strings"

	"github.com/xeipuuv/gojsonreference"
	"github.com/xeipuuv/gojsonschema"
)

const (
	keySchema = "$schema"
	draft04   = "http://json-schema.org/draft-04/schema#"
	draft06   = "http://json-schema.org/draft-06/schema#"
	draft07   = "http://json-schema.org/draft-07/schema#"
)

// Schema contains a JSON v4 schema.
type Schema struct {
	s interface{}
}

// NewSchemaFile loads a JSON v4 schema from the given path.
func NewSchemaFile(path string) (*Schema, error) {
	u, err := toURL(nil, path)
	if err != nil {
		return nil, wrapError(err, "Failed to obtain url to schema file")
	}

	return NewSchemaURL(u)
}

// NewSchemaURL loads a JSON v4 schema from the given URL.
func NewSchemaURL(url url.URL) (*Schema, error) {
	data, err := loadURL(url)
	if err != nil {
		return nil, wrapError(err, "Failed to load schema url %v", url)
	}

	return NewSchemaData(data)
}

// NewSchemaData loads a JSON v4 schema from the given data.
func NewSchemaData(data []byte) (*Schema, error) {
	var s interface{}

	err := JSONUnmarshal(data, &s)
	if err != nil {
		return nil, wrapError(err, "Schema is not valid json")
	}

	return NewSchemaGo(s)
}

// NewSchemaGo creates a Schema instance from a schema represented as a golang object.
func NewSchemaGo(s interface{}) (*Schema, error) {
	// validate if the schema is properly constructed by its specified draft
	draft, err := validateSchema(s)
	if err != nil {
		return nil, wrapError(err, "The schema is not valid against the meta-schema "+draft)
	}

	return &Schema{s: s}, nil
}

// Validate checks the given golang data against the schema.
func (s *Schema) Validate(data interface{}) error {
	if s == nil {
		return makeError("Schema is not set")
	}

	return validate(data, s.s)
}

// ApplyDefaults adds default values defined in the schema to the data pointed to by pData.
func (s *Schema) ApplyDefaults(pData interface{}) error {
	if s == nil {
		return makeError("Schema is not set")
	}

	return applyDefaults(pData, s.s)
}

var metaSchema interface{}

func updateMetaSchema(s interface{}) (draft string, err error) {
	m, ok := s.(map[string]interface{})
	if !ok {
		return "Unknown", makeError("Invalid schema structure")
	}

	// use schema draft 04 if we don't have a key to specify it
	draft = draft04
	data := metaSchemaData[draft]

	if v, ok := m[keySchema]; ok {
		draft = fmt.Sprintf("%v", v)
		if d, ok := metaSchemaData[draft]; ok {
			data = d
		}
	}

	err = JSONUnmarshal(data, &metaSchema)
	if err != nil {
		return draft, wrapError(err, "Could not load json meta-schema")
	}

	return draft, nil
}

func validateSchema(schema interface{}) (string, error) {
	schemaLoader := gojsonschema.NewGoLoader(schema)
	sl := gojsonschema.NewSchemaLoader()
	sl.AutoDetect = true
	sl.Validate = true

	err := sl.AddSchemas(schemaLoader)
	if err != nil {
		draft := fmt.Sprintf("Draft0%v", sl.Draft)
		if sl.Draft == math.MaxInt32 {
			draft = "hybrid"
		}

		return draft, wrapError(err, "Schema validation failed")
	}

	draft, err := updateMetaSchema(schema)
	if err != nil {
		return draft, wrapError(err, "Cannot access the schema draft")
	}

	return draft, validate(schema, metaSchema)
}

func validate(data, schema interface{}) error {
	dataLoader := gojsonschema.NewGoLoader(data)
	schemaLoader := gojsonschema.NewGoLoader(schema)

	formatErrs.clear()

	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		return wrapError(err, "An error occurred during validation")
	}

	err = processResult(result)

	return wrapError(err, "Schema validation failed")
}

func processResult(result *gojsonschema.Result) error {
	if !result.Valid() {
		err := makeError("The document is not valid against the schema")

		for _, rerr := range result.Errors() {
			ctx := convertJSONContext(rerr.Context().String())
			ctxErr := makeContextError(ctx, rerr.Description())

			ferr := formatErrs.get(rerr.Details()["format"], rerr.Value())
			if ferr != nil {
				ctxErr = detailError(ctxErr, ferr.Error())
			}

			err = wrapError(ctxErr, err.Error())
		}

		return err
	}

	return nil
}

func convertJSONContext(jsonCtx string) context {
	parts := strings.Split(jsonCtx, ".")

	return rootContext().add(parts[1:]...)
}

func applyDefaults(pData, schema interface{}) error {
	err := applyDefaultsRecursive(rootContext(), schema, pData, schema)

	return wrapError(err, "The defaults could not be applied")
}

func applyDefaultsRecursive(ctx context, rootSchema, pData, schema interface{}) error {
	if pData == nil {
		return makeContextError(ctx, "Destination value must not be nil")
	}

	pDataVal := reflect.ValueOf(pData)
	if pDataVal.Kind() != reflect.Ptr {
		return makeContextError(ctx, "Destination value must be a pointer")
	}

	dataVal := pDataVal.Elem()
	data := dataVal.Interface()

	schemaNode, ok := schema.(map[string]interface{})
	if !ok {
		return makeContextError(ctx, "Schema section is not a map")
	}

	val, ok := schemaNode["$ref"]
	if ok {
		ref, ok := val.(string)
		if !ok {
			return makeContextError(ctx, makeError("Reference is not a string '%v'", ref).Error())
		}

		jref, err := gojsonreference.NewJsonReference(ref)
		if err != nil {
			return makeContextError(ctx, wrapError(err, "Invalid reference '%v'", ref).Error())
		}

		subSchema, _, err := jref.GetPointer().Get(rootSchema)
		if subSchema == nil || err != nil {
			return makeContextError(ctx, wrapError(err, "Cannot find reference '%v'", ref).Error())
		}

		return applyDefaultsRecursive(ctx.add(ref), rootSchema, pData, subSchema)
	}

	schemaType, ok := schemaNode["type"]
	if !ok {
		if hasKey(schemaNode, "anyOf", "allOf", "oneOf", "not") {
			// the schema is valid, so it is not an error, but we do not support these types of schema yet
			return nil
		}

		return makeContextError(ctx, "Schema section does not have a valid 'type' attribute")
	}

	if value, ok := schemaNode["default"]; ok && data == nil {
		defaultVal := reflect.ValueOf(value)
		dataVal.Set(defaultVal)
		data = dataVal.Interface()
	}

	var err error

	switch schemaType {
	case "object":
		err = applyObjectDefaults(ctx, rootSchema, data, schemaNode)
	case "array":
		err = applyArrayDefaults(ctx, rootSchema, data, schemaNode)
	}

	return err
}

func hasKey(m map[string]interface{}, keys ...string) bool {
	for _, key := range keys {
		if _, ok := m[key]; ok {
			return true
		}
	}

	return false
}

func applyObjectDefaults(ctx context, rootSchema, data interface{}, schemaNode map[string]interface{}) error {
	if data == nil {
		return nil
	}

	dataProps, ok := data.(map[string]interface{})
	if !ok {
		return makeContextError(ctx, "Node should be an 'object'")
	}

	if dataProps == nil {
		return nil
	}

	var schemaProps map[string]interface{}
	if props, ok := schemaNode["properties"]; ok {
		schemaProps = props.(map[string]interface{})
		for name, schemaProp := range schemaProps {
			dataProp := dataProps[name]

			err := applyDefaultsRecursive(ctx.add(name), rootSchema, &dataProp, schemaProp)
			if err != nil {
				return wrapError(err, "Failed to apply defaults to object property")
			}

			if dataProp != nil {
				dataProps[name] = dataProp
			}
		}
	}

	//nolint:nestif // to be refactored
	if addProps, ok := schemaNode["additionalProperties"]; ok {
		if addProps, ok = addProps.(map[string]interface{}); ok {
			for name, dataProp := range dataProps {
				if schemaProps == nil || schemaProps[name] == nil {
					err := applyDefaultsRecursive(ctx.add(name), rootSchema, &dataProp, addProps) //nolint:gosec,scopelint // to be refactored carefully
					if err != nil {
						return wrapError(err, "Failed to apply defaults to additional object property")
					}

					if dataProp != nil {
						dataProps[name] = dataProp
					}
				}
			}
		}
	}

	return nil
}

func applyArrayDefaults(ctx context, rootSchema, data interface{}, schemaNode map[string]interface{}) error {
	if data == nil {
		return nil
	}

	dataItems, ok := data.([]interface{})
	if !ok {
		return makeContextError(ctx, "Node should be an 'array'")
	}

	if items, ok := schemaNode["items"]; ok {
		schemaItem := items.(map[string]interface{})

		for i, dataItem := range dataItems {
			err := applyDefaultsRecursive(ctx.addInt(i), rootSchema, &dataItem, schemaItem) //nolint:gosec,scopelint // to be refactored carefully
			if err != nil {
				return wrapError(err, "Failed to apply defaults to array item")
			}

			if dataItem != nil {
				dataItems[i] = dataItem
			}
		}
	}

	return nil
}

var metaSchemaData = map[string][]byte{
	draft04: []byte(`
{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "description": "Core schema meta-schema",
    "definitions": {
        "schemaArray": {
            "type": "array",
            "minItems": 1,
            "items": { "$ref": "#" }
        },
        "positiveInteger": {
            "type": "integer",
            "minimum": 0
        },
        "positiveIntegerDefault0": {
            "allOf": [ { "$ref": "#/definitions/positiveInteger" }, { "default": 0 } ]
        },
        "simpleTypes": {
            "enum": [ "array", "boolean", "integer", "null", "number", "object", "string" ]
        },
        "stringArray": {
            "type": "array",
            "items": { "type": "string" },
            "minItems": 1,
            "uniqueItems": true
        }
    },
    "type": "object",
    "properties": {
        "id": {
            "type": "string",
            "format": "uri"
        },
        "$schema": {
            "type": "string",
            "format": "uri"
        },
        "title": {
            "type": "string"
        },
        "description": {
            "type": "string"
        },
        "default": {},
        "multipleOf": {
            "type": "number",
            "minimum": 0,
            "exclusiveMinimum": true
        },
        "maximum": {
            "type": "number"
        },
        "exclusiveMaximum": {
            "type": "boolean",
            "default": false
        },
        "minimum": {
            "type": "number"
        },
        "exclusiveMinimum": {
            "type": "boolean",
            "default": false
        },
        "maxLength": { "$ref": "#/definitions/positiveInteger" },
        "minLength": { "$ref": "#/definitions/positiveIntegerDefault0" },
        "pattern": {
            "type": "string",
            "format": "regex"
        },
        "additionalItems": {
            "anyOf": [
                { "type": "boolean" },
                { "$ref": "#" }
            ],
            "default": {}
        },
        "items": {
            "anyOf": [
                { "$ref": "#" },
                { "$ref": "#/definitions/schemaArray" }
            ],
            "default": {}
        },
        "maxItems": { "$ref": "#/definitions/positiveInteger" },
        "minItems": { "$ref": "#/definitions/positiveIntegerDefault0" },
        "uniqueItems": {
            "type": "boolean",
            "default": false
        },
        "maxProperties": { "$ref": "#/definitions/positiveInteger" },
        "minProperties": { "$ref": "#/definitions/positiveIntegerDefault0" },
        "required": { "$ref": "#/definitions/stringArray" },
        "additionalProperties": {
            "anyOf": [
                { "type": "boolean" },
                { "$ref": "#" }
            ],
            "default": {}
        },
        "definitions": {
            "type": "object",
            "additionalProperties": { "$ref": "#" },
            "default": {}
        },
        "properties": {
            "type": "object",
            "additionalProperties": { "$ref": "#" },
            "default": {}
        },
        "patternProperties": {
            "type": "object",
            "additionalProperties": { "$ref": "#" },
            "default": {}
        },
        "dependencies": {
            "type": "object",
            "additionalProperties": {
                "anyOf": [
                    { "$ref": "#" },
                    { "$ref": "#/definitions/stringArray" }
                ]
            }
        },
        "enum": {
            "type": "array",
            "minItems": 1,
            "uniqueItems": true
        },
        "type": {
            "anyOf": [
                { "$ref": "#/definitions/simpleTypes" },
                {
                    "type": "array",
                    "items": { "$ref": "#/definitions/simpleTypes" },
                    "minItems": 1,
                    "uniqueItems": true
                }
            ]
        },
        "allOf": { "$ref": "#/definitions/schemaArray" },
        "anyOf": { "$ref": "#/definitions/schemaArray" },
        "oneOf": { "$ref": "#/definitions/schemaArray" },
        "not": { "$ref": "#" }
    },
    "dependencies": {
        "exclusiveMaximum": [ "maximum" ],
        "exclusiveMinimum": [ "minimum" ]
    },
    "default": {}
}`),
	draft06: []byte(`{
		"$schema": "http://json-schema.org/draft-06/schema#",
		"$id": "http://json-schema.org/draft-06/schema#",
		"title": "Core schema meta-schema",
		"definitions": {
			"schemaArray": {
				"type": "array",
				"minItems": 1,
				"items": { "$ref": "#" }
			},
			"nonNegativeInteger": {
				"type": "integer",
				"minimum": 0
			},
			"nonNegativeIntegerDefault0": {
				"allOf": [
					{ "$ref": "#/definitions/nonNegativeInteger" },
					{ "default": 0 }
				]
			},
			"simpleTypes": {
				"enum": [
					"array",
					"boolean",
					"integer",
					"null",
					"number",
					"object",
					"string"
				]
			},
			"stringArray": {
				"type": "array",
				"items": { "type": "string" },
				"uniqueItems": true,
				"default": []
			}
		},
		"type": ["object", "boolean"],
		"properties": {
			"$id": {
				"type": "string",
				"format": "uri-reference"
			},
			"$schema": {
				"type": "string",
				"format": "uri"
			},
			"$ref": {
				"type": "string",
				"format": "uri-reference"
			},
			"title": {
				"type": "string"
			},
			"description": {
				"type": "string"
			},
			"default": {},
			"examples": {
				"type": "array",
				"items": {}
			},
			"multipleOf": {
				"type": "number",
				"exclusiveMinimum": 0
			},
			"maximum": {
				"type": "number"
			},
			"exclusiveMaximum": {
				"type": "number"
			},
			"minimum": {
				"type": "number"
			},
			"exclusiveMinimum": {
				"type": "number"
			},
			"maxLength": { "$ref": "#/definitions/nonNegativeInteger" },
			"minLength": { "$ref": "#/definitions/nonNegativeIntegerDefault0" },
			"pattern": {
				"type": "string",
				"format": "regex"
			},
			"additionalItems": { "$ref": "#" },
			"items": {
				"anyOf": [
					{ "$ref": "#" },
					{ "$ref": "#/definitions/schemaArray" }
				],
				"default": {}
			},
			"maxItems": { "$ref": "#/definitions/nonNegativeInteger" },
			"minItems": { "$ref": "#/definitions/nonNegativeIntegerDefault0" },
			"uniqueItems": {
				"type": "boolean",
				"default": false
			},
			"contains": { "$ref": "#" },
			"maxProperties": { "$ref": "#/definitions/nonNegativeInteger" },
			"minProperties": { "$ref": "#/definitions/nonNegativeIntegerDefault0" },
			"required": { "$ref": "#/definitions/stringArray" },
			"additionalProperties": { "$ref": "#" },
			"definitions": {
				"type": "object",
				"additionalProperties": { "$ref": "#" },
				"default": {}
			},
			"properties": {
				"type": "object",
				"additionalProperties": { "$ref": "#" },
				"default": {}
			},
			"patternProperties": {
				"type": "object",
				"additionalProperties": { "$ref": "#" },
				"default": {}
			},
			"dependencies": {
				"type": "object",
				"additionalProperties": {
					"anyOf": [
						{ "$ref": "#" },
						{ "$ref": "#/definitions/stringArray" }
					]
				}
			},
			"propertyNames": { "$ref": "#" },
			"const": {},
			"enum": {
				"type": "array",
				"minItems": 1,
				"uniqueItems": true
			},
			"type": {
				"anyOf": [
					{ "$ref": "#/definitions/simpleTypes" },
					{
						"type": "array",
						"items": { "$ref": "#/definitions/simpleTypes" },
						"minItems": 1,
						"uniqueItems": true
					}
				]
			},
			"format": { "type": "string" },
			"allOf": { "$ref": "#/definitions/schemaArray" },
			"anyOf": { "$ref": "#/definitions/schemaArray" },
			"oneOf": { "$ref": "#/definitions/schemaArray" },
			"not": { "$ref": "#" }
		},
		"default": {}
	}`),
	draft07: []byte(`{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"$id": "http://json-schema.org/draft-07/schema#",
		"title": "Core schema meta-schema",
		"definitions": {
			"schemaArray": {
				"type": "array",
				"minItems": 1,
				"items": { "$ref": "#" }
			},
			"nonNegativeInteger": {
				"type": "integer",
				"minimum": 0
			},
			"nonNegativeIntegerDefault0": {
				"allOf": [
					{ "$ref": "#/definitions/nonNegativeInteger" },
					{ "default": 0 }
				]
			},
			"simpleTypes": {
				"enum": [
					"array",
					"boolean",
					"integer",
					"null",
					"number",
					"object",
					"string"
				]
			},
			"stringArray": {
				"type": "array",
				"items": { "type": "string" },
				"uniqueItems": true,
				"default": []
			}
		},
		"type": ["object", "boolean"],
		"properties": {
			"$id": {
				"type": "string",
				"format": "uri-reference"
			},
			"$schema": {
				"type": "string",
				"format": "uri"
			},
			"$ref": {
				"type": "string",
				"format": "uri-reference"
			},
			"$comment": {
				"type": "string"
			},
			"title": {
				"type": "string"
			},
			"description": {
				"type": "string"
			},
			"default": true,
			"readOnly": {
				"type": "boolean",
				"default": false
			},
			"examples": {
				"type": "array",
				"items": true
			},
			"multipleOf": {
				"type": "number",
				"exclusiveMinimum": 0
			},
			"maximum": {
				"type": "number"
			},
			"exclusiveMaximum": {
				"type": "number"
			},
			"minimum": {
				"type": "number"
			},
			"exclusiveMinimum": {
				"type": "number"
			},
			"maxLength": { "$ref": "#/definitions/nonNegativeInteger" },
			"minLength": { "$ref": "#/definitions/nonNegativeIntegerDefault0" },
			"pattern": {
				"type": "string",
				"format": "regex"
			},
			"additionalItems": { "$ref": "#" },
			"items": {
				"anyOf": [
					{ "$ref": "#" },
					{ "$ref": "#/definitions/schemaArray" }
				],
				"default": true
			},
			"maxItems": { "$ref": "#/definitions/nonNegativeInteger" },
			"minItems": { "$ref": "#/definitions/nonNegativeIntegerDefault0" },
			"uniqueItems": {
				"type": "boolean",
				"default": false
			},
			"contains": { "$ref": "#" },
			"maxProperties": { "$ref": "#/definitions/nonNegativeInteger" },
			"minProperties": { "$ref": "#/definitions/nonNegativeIntegerDefault0" },
			"required": { "$ref": "#/definitions/stringArray" },
			"additionalProperties": { "$ref": "#" },
			"definitions": {
				"type": "object",
				"additionalProperties": { "$ref": "#" },
				"default": {}
			},
			"properties": {
				"type": "object",
				"additionalProperties": { "$ref": "#" },
				"default": {}
			},
			"patternProperties": {
				"type": "object",
				"additionalProperties": { "$ref": "#" },
				"propertyNames": { "format": "regex" },
				"default": {}
			},
			"dependencies": {
				"type": "object",
				"additionalProperties": {
					"anyOf": [
						{ "$ref": "#" },
						{ "$ref": "#/definitions/stringArray" }
					]
				}
			},
			"propertyNames": { "$ref": "#" },
			"const": true,
			"enum": {
				"type": "array",
				"items": true,
				"minItems": 1,
				"uniqueItems": true
			},
			"type": {
				"anyOf": [
					{ "$ref": "#/definitions/simpleTypes" },
					{
						"type": "array",
						"items": { "$ref": "#/definitions/simpleTypes" },
						"minItems": 1,
						"uniqueItems": true
					}
				]
			},
			"format": { "type": "string" },
			"contentMediaType": { "type": "string" },
			"contentEncoding": { "type": "string" },
			"if": { "$ref": "#" },
			"then": { "$ref": "#" },
			"else": { "$ref": "#" },
			"allOf": { "$ref": "#/definitions/schemaArray" },
			"anyOf": { "$ref": "#/definitions/schemaArray" },
			"oneOf": { "$ref": "#/definitions/schemaArray" },
			"not": { "$ref": "#" }
		},
		"default": true
	}`),
}
