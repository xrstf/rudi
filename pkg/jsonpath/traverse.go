// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import (
	"errors"
	"fmt"
	"reflect"
)

// traverse takes any sort of value and a path and will return a slice
// of pointers to all the values the function visited along the way.
// For a successful traversal (err==nil), the resulting slice will have
// the same length as the input path.
// If the path points to a nil value, the last element of the return
// slice will be nil.
func traverse(value any, path Path) ([]*any, error) {
	if len(path) == 0 {
		return nil, nil
	}

	steps := []*any{}

	// fmt.Printf("Getting %v…\n", path)

	for _, step := range path {
		if value == nil {
			return nil, errors.New("cannot descend into nil")
		}

		// fmt.Printf("* value: %v\n", value)

		// determine the current value's type
		valueType := reflect.TypeOf(value)
		elemType := valueType

		// get the type's kind
		valueKind := valueType.Kind()
		elemKind := valueKind

		// fmt.Printf("  kitty: %v (%v)\n", valueType, elemType)

		rValue := reflect.ValueOf(value)

		// unwrap pointer types to their underlying types (*int => int)
		if valueKind == reflect.Pointer {
			if rValue.IsNil() {
				return nil, errors.New("cannot descend into nil")
			}

			elemType = valueType.Elem()
			elemKind = elemType.Kind()

			// dereference the pointer
			rValue = rValue.Elem()

			// fmt.Printf("  pkitty: %v (%v)\n", elemType, elemKind)
		}

		var err error

		switch elemKind {
		case reflect.Slice, reflect.Array:
			value, err = traverseSlice(rValue, step)
			if err != nil {
				return nil, err
			}

		case reflect.Map:
			value, err = traverseMap(rValue, step)
			if err != nil {
				return nil, err
			}

		case reflect.Struct:
			value, err = traverseStruct(value, rValue, step)
			if err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("cannot descend with %v (%T) into %T", step, step, value)
		}

		steps = append(steps, &value)
	}

	return steps, nil
}

func traverseSlice(value reflect.Value, step any) (any, error) {
	index, ok := toIntegerStep(step)
	if !ok {
		return nil, fmt.Errorf("cannot use %v as an array index", step)
	}

	if index < 0 || index >= value.Len() {
		return nil, fmt.Errorf("index %d out of bounds", index)
	}

	return value.Index(index).Interface(), nil
}

func traverseMap(value reflect.Value, step any) (any, error) {
	key, ok := toStringStep(step)
	if !ok {
		return nil, fmt.Errorf("cannot use %v as an object key", step)
	}

	indexValue := value.MapIndex(reflect.ValueOf(key))
	if indexValue == (reflect.Value{}) {
		return nil, fmt.Errorf("no such key: %q", key)
	}

	return indexValue.Interface(), nil
}

func traverseStruct(value any, rValue reflect.Value, step any) (any, error) {
	if vectorReader, ok := value.(VectorReader); ok {
		index, ok := toIntegerStep(step)
		if ok {
			value, err := vectorReader.GetVectorItem(index)
			if err != nil {
				return nil, fmt.Errorf("cannot descend with %v (%T) into %T: %w", step, step, value, err)
			}

			return value, nil
		}
	}

	if objectReader, ok := value.(ObjectReader); ok {
		key, ok := toStringStep(step)
		if ok {
			value, err := objectReader.GetObjectKey(key)
			if err != nil {
				return nil, fmt.Errorf("cannot descend with %v (%T) into %T: %w", step, step, value, err)
			}

			return value, nil
		}
	}

	key, ok := toStringStep(step)
	if !ok {
		return nil, fmt.Errorf("cannot use %v as an object key", step)
	}

	fieldValue := rValue.FieldByName(key)
	if fieldValue == (reflect.Value{}) || !fieldValue.CanInterface() {
		return nil, fmt.Errorf("no such field: %q", key)
	}

	return fieldValue.Interface(), nil
}
