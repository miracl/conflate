package conflate

import (
	"fmt"
	"reflect"
)

func mergeTo(toData any, fromData ...any) error {
	for _, fromDatum := range fromData {
		err := merge(toData, fromDatum)
		if err != nil {
			return err
		}
	}

	return nil
}

func merge(pToData, fromData any) error {
	return mergeRecursive(rootContext(), pToData, fromData)
}

func mergeRecursive(ctx context, pToData, fromData any) error {
	if pToData == nil {
		return &contextError{
			context: ctx,
			msg:     "the destination variable must not be nil",
		}
	}

	pToVal := reflect.ValueOf(pToData)
	if pToVal.Kind() != reflect.Pointer {
		return &contextError{
			context: ctx,
			msg:     "the destination variable must be a pointer",
		}
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

	//nolint:exhaustive // to be refactored
	switch fromVal.Kind() {
	case reflect.Map:
		err = mergeMapRecursive(ctx, toData, fromData)
	case reflect.Slice:
		err = mergeSliceRecursive(ctx, toVal, toData, fromData)
	default:
		err = mergeDefaultRecursive(ctx, toVal, fromVal, toData, fromData)
	}

	return err
}

func mergeMapRecursive(ctx context, toData, fromData any) error {
	fromProps, ok := fromData.(map[string]any)
	if !ok {
		return &contextError{
			context: ctx,
			msg:     "the source value must be a map[string]interface{}",
		}
	}

	toProps, ok := toData.(map[string]any)
	if toProps == nil || !ok {
		return &contextError{
			context: ctx,
			msg:     "the destination value must be a map[string]interface{}",
		}
	}

	for name, fromProp := range fromProps {
		if val := toProps[name]; val == nil {
			toProps[name] = fromProp
		} else {
			err := merge(&val, fromProp)
			if err != nil {
				return &contextError{
					context: ctx.add(name),
					msg:     fmt.Sprintf("failed to merge object property : %v : %v", name, err.Error()),
				}
			}

			toProps[name] = val
		}
	}

	return nil
}

func mergeSliceRecursive(ctx context, toVal reflect.Value, toData, fromData any) error {
	fromItems, ok := fromData.([]any)
	if !ok {
		return &contextError{
			context: ctx,
			msg:     "the source value must be a []interface{}",
		}
	}

	toItems, ok := toData.([]any)
	if toItems == nil || !ok {
		return &contextError{
			context: ctx,
			msg:     "the destination value must be a []interface{}",
		}
	}

	toItems = append(toItems, fromItems...)
	toVal.Set(reflect.ValueOf(toItems))

	return nil
}

func mergeDefaultRecursive(ctx context, toVal, fromVal reflect.Value, toData, fromData any) error {
	if reflect.DeepEqual(toData, fromData) {
		return nil
	}

	fromType := fromVal.Type()
	toType := toVal.Type()

	if toType.Kind() == reflect.Interface {
		toType = toVal.Elem().Type()
	}

	if !fromType.AssignableTo(toType) {
		return &contextError{
			context: ctx,
			msg:     fmt.Sprintf("the destination type (%v) must be the same as the source type (%v)", toType, fromType),
		}
	}

	toVal.Set(fromVal)

	return nil
}
