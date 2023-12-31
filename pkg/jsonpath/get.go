// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import (
	"errors"
	"fmt"
	"reflect"
)

type ObjectReader interface {
	GetObjectKey(name string) (any, error)
}

type VectorReader interface {
	GetVectorItem(index int) (any, error)
}

func Get(value any, path Path) (any, error) {
	if len(path) == 0 {
		return value, nil
	}

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

		switch elemKind {
		case reflect.Slice, reflect.Array:
			index, ok := toIntegerStep(step)
			if !ok {
				return nil, fmt.Errorf("cannot use %v as an array index", step)
			}

			if index < 0 || index >= rValue.Len() {
				return nil, fmt.Errorf("index %d out of bounds", index)
			}

			value = rValue.Index(index).Interface()

		case reflect.Map:
			key, ok := toStringStep(step)
			if !ok {
				return nil, fmt.Errorf("cannot use %v as an object key", step)
			}

			indexValue := rValue.MapIndex(reflect.ValueOf(key))
			if indexValue == (reflect.Value{}) {
				return nil, fmt.Errorf("no such key: %q", key)
			}

			value = indexValue.Interface()

		case reflect.Struct:
			if vectorReader, ok := value.(VectorReader); ok {
				index, ok := toIntegerStep(step)
				if ok {
					var err error

					value, err = vectorReader.GetVectorItem(index)
					if err != nil {
						return nil, fmt.Errorf("cannot descend with %v (%T) into %T: %w", step, step, value, err)
					}

					continue
				}
			}

			if objectReader, ok := value.(ObjectReader); ok {
				key, ok := toStringStep(step)
				if ok {
					var err error

					value, err = objectReader.GetObjectKey(key)
					if err != nil {
						return nil, fmt.Errorf("cannot descend with %v (%T) into %T: %w", step, step, value, err)
					}

					continue
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

			value = fieldValue.Interface()

		default:
			return nil, fmt.Errorf("cannot descend with %v (%T) into %T", step, step, value)
		}
	}

	return value, nil
}
