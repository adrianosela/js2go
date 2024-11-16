//go:build js

package js2go

import (
	"errors"
	"fmt"
	"reflect"
	"syscall/js"
)

// Encode takes an input struct and converts it into a js.Value.
func Encode(input any) (js.Value, error) {
	output := js.Global().Get("Object").New()

	// ensure pointer is a struct or pointer to a struct
	inputValue := reflect.ValueOf(input)
	if inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
	}
	if inputValue.Kind() != reflect.Struct {
		return output, errors.New("input must be a struct or pointer to a struct")
	}

	inputValueType := inputValue.Type()
	for i := 0; i < inputValueType.NumField(); i++ {
		field := inputValue.Field(i)
		tag := inputValueType.Field(i).Tag.Get("js")
		if tag == "" {
			continue
		}
		jsValue, err := getJSValue(field)
		if err != nil {
			return js.Value{}, fmt.Errorf("error encoding field '%s' at index %d: %w", inputValueType.Field(i).Name, i, err)
		}
		output.Set(tag, jsValue)
	}

	return output, nil
}

// getJSValue converts a reflect.Value to js.Value based on its type.
func getJSValue(field reflect.Value) (js.Value, error) {
	switch field.Kind() {
	case reflect.String:
		return js.ValueOf(field.String()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return js.ValueOf(field.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return js.ValueOf(field.Uint()), nil
	case reflect.Bool:
		return js.ValueOf(field.Bool()), nil
	case reflect.Float32, reflect.Float64:
		return js.ValueOf(field.Float()), nil
	case reflect.Slice, reflect.Array:
		jsArray := js.Global().Get("Array").New(field.Len())
		for i := 0; i < field.Len(); i++ {
			jsval, err := getJSValue(field.Index(i))
			if err != nil {
				return js.Value{}, fmt.Errorf("error processing array element at index %d: %w", i, err)
			}
			jsArray.SetIndex(i, jsval)
		}
		return jsArray, nil
	case reflect.Struct:
		jsval, err := Encode(field.Interface())
		if err != nil {
			return js.Value{}, fmt.Errorf("error processing nested struct: %w", err)
		}
		return jsval, nil
	case reflect.Map:
		jsMap := js.Global().Get("Object").New()
		iter := field.MapRange()
		for iter.Next() {
			keyStr := fmt.Sprint(iter.Key())
			val, err := getJSValue(iter.Value())
			if err != nil {
				return js.Value{}, fmt.Errorf("error processing map value for key '%s': %w", keyStr, err)
			}
			jsMap.Set(keyStr, val)
		}
		return jsMap, nil
	default:
		return js.Undefined(), fmt.Errorf("unsupported field type %s", field.Type())
	}
}
