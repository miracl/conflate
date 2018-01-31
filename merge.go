package conflate

import (
	"reflect"
)

func mergeAll(fromData ...interface{}) (interface{}, error) {
	var toData interface{}
	err := mergeTo(&toData, fromData...)
	if err != nil {
		return nil, err
	}
	return toData, nil
}

func mergeTo(toData interface{}, fromData ...interface{}) error {
	for _, fromDatum := range fromData {
		err := merge(toData, fromDatum)
		if err != nil {
			return err
		}
	}
	return nil
}

func merge(pToData interface{}, fromData interface{}) error {
	return mergeRecursive(rootContext(), pToData, fromData)
}

func mergeRecursive(ctx context, pToData interface{}, fromData interface{}) error {
	if pToData == nil {
		return makeContextError(ctx, "The destination variable must not be nil")
	}
	pToVal := reflect.ValueOf(pToData)
	if pToVal.Kind() != reflect.Ptr {
		return makeContextError(ctx, "The destination variable must be a pointer")
	}

	if fromData == nil {
		return nil
	}

	toVal := pToVal.Elem()
	fromVal := reflect.ValueOf(fromData)

	toData := toVal.Interface()
	if toVal.Interface() == nil {
		toVal.Set(fromVal)
		return nil
	}

	switch fromVal.Kind() {
	case reflect.Map:
		fromProps, ok := fromData.(map[string]interface{})
		if !ok {
			return makeContextError(ctx, "The source value must be a map[string]interface{}")
		}
		toProps, ok := toData.(map[string]interface{})
		if !ok {
			return makeContextError(ctx, "The destination value must be a map[string]interface{}")
		}
		for name, fromProp := range fromProps {
			if val := toProps[name]; val == nil {
				toProps[name] = fromProp
			} else {
				err := merge(&val, fromProp)
				if err != nil {
					return makeContextError(ctx.add(name), "Failed to merge object property : %v : %v", name, err)
				}
				toProps[name] = val
			}
		}
	case reflect.Slice:
		fromItems, ok := fromData.([]interface{})
		if !ok {
			return makeContextError(ctx, "The source value must be a []interface{}")
		}
		toItems, ok := toData.([]interface{})
		if !ok {
			return makeContextError(ctx, "The destination value must be a []interface{}")
		}
		toItems = append(toItems, fromItems...)
		toVal.Set(reflect.ValueOf(toItems))

	default:
		if reflect.DeepEqual(toData, fromData) {
			return nil
		}
		toVal.Set(fromVal)
	}
	return nil
}
