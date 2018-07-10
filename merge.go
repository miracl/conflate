package conflate

import (
	"reflect"
)

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

	var err error
	switch fromVal.Kind() {
	case reflect.Map:
		err = mergeMapRecursive(ctx, toVal, fromVal, toData, fromData)
	case reflect.Slice:
		err = mergeSliceRecursive(ctx, toVal, fromVal, toData, fromData)
	default:
		err = mergeDefaultRecursive(ctx, toVal, fromVal, toData, fromData)
	}
	return err
}

func mergeMapRecursive(ctx context, toVal reflect.Value, fromVal reflect.Value,
	toData interface{}, fromData interface{}) error {

	fromProps, ok := fromData.(map[string]interface{})
	if !ok {
		return makeContextError(ctx, "The source value must be a map[string]interface{}")
	}
	toProps, _ := toData.(map[string]interface{})
	if toProps == nil {
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
	return nil
}

func mergeSliceRecursive(ctx context, toVal reflect.Value, fromVal reflect.Value,
	toData interface{}, fromData interface{}) error {

	fromItems, ok := fromData.([]interface{})
	if !ok {
		return makeContextError(ctx, "The source value must be a []interface{}")
	}
	toItems, _ := toData.([]interface{})
	if toItems == nil {
		return makeContextError(ctx, "The destination value must be a []interface{}")
	}
	toItems = append(toItems, fromItems...)
	toVal.Set(reflect.ValueOf(toItems))
	return nil
}

func mergeDefaultRecursive(ctx context, toVal reflect.Value, fromVal reflect.Value,
	toData interface{}, fromData interface{}) error {

	if reflect.DeepEqual(toData, fromData) {
		return nil
	}
	fromType := fromVal.Type()
	toType := toVal.Type()
	if toType.Kind() == reflect.Interface {
		toType = toVal.Elem().Type()
	}
	if !fromType.AssignableTo(toType) {
		return makeContextError(ctx, "The destination type (%v) must be the same as the source type (%v)", toType, fromType)
	}
	toVal.Set(fromVal)
	return nil
}
