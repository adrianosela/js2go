//go:build js

package js2go

import (
	"fmt"
	"reflect"
	"syscall/js"
)

// Decode takes a js.Value and decodes it into a struct.
func Decode(input js.Value, result any) error {
	val := reflect.ValueOf(result)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("result argument must be a pointer to a struct")
	}

	val = val.Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}
		tag := typ.Field(i).Tag.Get("js")
		if tag == "" || input.Get(tag).IsUndefined() {
			continue
		}
		if err := setFieldValue(field, input.Get(tag)); err != nil {
			return fmt.Errorf("error setting field '%s' at index %d: %w", typ.Field(i).Name, i, err)
		}
	}

	return nil
}

// setFieldValue assigns js.Value to the corresponding field of a Go struct.
func setFieldValue(field reflect.Value, jsVal js.Value) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(jsVal.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if jsVal.Type() != js.TypeNumber {
			return fmt.Errorf("expected number for integer field, got %s", jsVal.Type().String())
		}
		field.SetInt(int64(jsVal.Float()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if jsVal.Type() != js.TypeNumber {
			return fmt.Errorf("expected number for unsigned integer field, got %s", jsVal.Type().String())
		}
		field.SetUint(uint64(jsVal.Float()))
	case reflect.Bool:
		if jsVal.Type() != js.TypeBoolean {
			return fmt.Errorf("expected boolean, got %s", jsVal.Type().String())
		}
		field.SetBool(jsVal.Bool())
	case reflect.Float32, reflect.Float64:
		if jsVal.Type() != js.TypeNumber {
			return fmt.Errorf("expected number for float field, got %s", jsVal.Type().String())
		}
		field.SetFloat(jsVal.Float())
	case reflect.Slice:
		return fillSlice(field, jsVal)
	case reflect.Struct:
		return Decode(jsVal, field.Addr().Interface())
	case reflect.Map:
		return fillMap(field, jsVal)
	default:
		return fmt.Errorf("unsupported field type %s", field.Type().String())
	}
	return nil
}

// fillSlice populates a Go slice from a js.Value array.
func fillSlice(slice reflect.Value, jsVal js.Value) error {
	if !jsVal.InstanceOf(js.Global().Get("Array")) {
		return fmt.Errorf("expected array, got %s", jsVal.Type().String())
	}
	length := jsVal.Length()
	newSlice := reflect.MakeSlice(slice.Type(), length, length)
	for i := 0; i < length; i++ {
		err := setFieldValue(newSlice.Index(i), jsVal.Index(i))
		if err != nil {
			return fmt.Errorf("error at array index %d: %w", i, err)
		}
	}
	slice.Set(newSlice)
	return nil
}

// fillMap populates a Go map from a js.Value object assumed to represent a map.
func fillMap(mapVal reflect.Value, jsVal js.Value) error {
	keys := jsVal.Call("keys")
	length := keys.Length()
	newMap := reflect.MakeMap(mapVal.Type())
	for i := 0; i < length; i++ {
		keyJs := keys.Index(i)
		key := reflect.New(mapVal.Type().Key()).Elem()
		if key.Kind() == reflect.String {
			key.SetString(keyJs.String())
		} else {
			return fmt.Errorf("map keys must be strings")
		}
		valJs := jsVal.Get(keyJs.String())
		val := reflect.New(mapVal.Type().Elem()).Elem()
		err := setFieldValue(val, valJs)
		if err != nil {
			return fmt.Errorf("error at map key '%s': %w", keyJs.String(), err)
		}
		newMap.SetMapIndex(key, val)
	}
	mapVal.Set(newMap)
	return nil
}
