package conflate

import (
	"github.com/xeipuuv/gojsonschema"
	"reflect"
	"strings"
)

var metaSchema interface{}

func validateSchema(schema interface{}) error {
	if metaSchema == nil {
		err := JSONUnmarshal(metaSchemaData, &metaSchema)
		if err != nil {
			return wrapError(err, "Could not load json meta-schema")
		}
	}
	return validate(schema, metaSchema)
}

func validate(data interface{}, schema interface{}) error {
	dataLoader := gojsonschema.NewGoLoader(data)
	schemaLoader := gojsonschema.NewGoLoader(schema)
	formatErrs.clear()
	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		return wrapError(err, "An error occurred during validation")
	}
	return processResult(result)
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

func applyDefaults(pData interface{}, schema interface{}) error {
	return applyDefaultsRecursive(rootContext(), pData, schema)
}

func applyDefaultsRecursive(ctx context, pData interface{}, schema interface{}) error {
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

	for k := range schemaNode {
		switch k {
		case "anyOf":
			fallthrough
		case "allOf":
			fallthrough
		case "oneOf":
			fallthrough
		case "notOf":
			// we do not support setting defaults in these structures, so we return without error
			return nil
		}
	}

	schemaType, ok := schemaNode["type"]
	if !ok {
		return makeContextError(ctx, "Schema section does not have a valid 'type' attribute")
	}

	if value, ok := schemaNode["default"]; ok && data == nil {
		defaultVal := reflect.ValueOf(value)
		dataVal.Set(defaultVal)
		data = dataVal.Interface()
	}

	switch schemaType {
	case "object":
		if data == nil {
			break
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
				err := applyDefaultsRecursive(ctx.add(name), &dataProp, schemaProp)
				if err != nil {
					return wrapError(err, "Failed to apply defaults to object property")
				}
				if dataProp != nil {
					dataProps[name] = dataProp
				}
			}
		}
		if addProps, ok := schemaNode["additionalProperties"]; ok {
			if addProps, ok = addProps.(map[string]interface{}); ok {
				for name, dataProp := range dataProps {
					if schemaProps == nil || schemaProps[name] == nil {
						err := applyDefaultsRecursive(ctx.add(name), &dataProp, addProps)
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
	case "array":
		if data == nil {
			break
		}
		dataItems, ok := data.([]interface{})
		if !ok {
			return makeContextError(ctx, "Node should be an 'array'")
		}
		if items, ok := schemaNode["items"]; ok {
			schemaItem := items.(map[string]interface{})
			for i, dataItem := range dataItems {
				err := applyDefaultsRecursive(ctx.addInt(i), &dataItem, schemaItem)
				if err != nil {
					return wrapError(err, "Failed to apply defaults to array item")
				}
				if dataItem != nil {
					dataItems[i] = dataItem
				}
			}
		}
	}
	return nil
}

var metaSchemaData = []byte(`
{
    "id": "http://json-schema.org/draft-04/schema#",
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
}`)
