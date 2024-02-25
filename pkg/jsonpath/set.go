// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import (
	"errors"
	"fmt"
	"reflect"
)

func Set(dest any, path Path, newValue any) (any, error) {
	return Patch(dest, path, func(_ bool, _ any, _ any) (any, error) {
		return newValue, nil
	})
}

type PatchFunc func(exists bool, key any, val any) (any, error)

func Patch(dest any, path Path, patchValue PatchFunc) (any, error) {
	if !path.IsValid() {
		return nil, errors.New("invalid path")
	}

	return patch(dest, nil, true, path, patchValue)
}

func patch(dest any, key any, exists bool, path Path, patchValue PatchFunc) (any, error) {
	if len(path) == 0 {
		return patchValue(exists, key, dest)
	}

	thisStep := path[0]
	remainingSteps := path[1:]

	foundKeyThings, foundValueThings, destKind, err := traverseStep(dest, thisStep)
	if err != nil && !errors.Is(err, noSuchKeyErr) && !errors.Is(err, indexOutOfBoundsErr) {
		return nil, err
	}

	switch thisStep.(type) {
	// $var[1], $var.foo, $var["foo"], $var[(+ 1 2)]
	case SingleStep:
		switch destKind {
		// arrays, slices
		// case kindList:
		// 	return patchFoundVectorValue(asVector, foundKeyThings, foundValueThings, err == nil, remainingSteps, patchValue)
		}

		switch foundKeyThings.(type) {
		case int:
			// nil values (or non-existing values) can be turned into vectors
			if dest == nil {
				dest = []any{}
			}

			asVector, ok := dest.([]any)
			if !ok {
				panic("VectorStep should have errored on a non-vector value.")
			}

			return patchFoundVectorValue(asVector, foundKeyThings, foundValueThings, err == nil, remainingSteps, patchValue)

		case string:
			// nil values (or non-existing values) can be turned into objects
			if dest == nil {
				dest = map[string]any{}
			}

			asObject, ok := dest.(map[string]any)
			if !ok {
				panic("ObjectStep should have errored on a non-object value.")
			}

			return patchFoundObjectValue(asObject, foundKeyThings, foundValueThings, err == nil, remainingSteps, patchValue)

		default:
			panic(fmt.Sprintf("SingleStep should have returned int index or string key, but returned %v (%T)", foundKeyThings, foundKeyThings))
		}

	// $var[?(…)]
	case FilterStep:
		foundValues := foundValueThings.([]any)
		if len(foundValues) == 0 {
			return dest, nil
		}

		foundsKeys, ok := foundKeyThings.([]string)
		if ok {
			// nil values (or non-existing values) can be turned into objects
			if dest == nil {
				dest = map[string]any{}
			}

			asObject, ok := dest.(map[string]any)
			if !ok {
				panic("ObjectStep should have errored on a non-object value.")
			}

			for idx, key := range foundsKeys {
				var err error
				asObject, err = patchFoundObjectValue(asObject, key, foundValues[idx], true, remainingSteps, patchValue)
				if err != nil {
					return nil, err
				}
			}

			return asObject, nil
		}

		foundIndexes, ok := foundKeyThings.([]int)
		if ok {
			// nil values (or non-existing values) can be turned into vectors
			if dest == nil {
				dest = []any{}
			}

			asVector, ok := dest.([]any)
			if !ok {
				panic("VectorStep should have errored on a non-vector value.")
			}

			for idx, vectorIndex := range foundIndexes {
				var err error
				asVector, err = patchFoundVectorValue(asVector, vectorIndex, foundValues[idx], true, remainingSteps, patchValue)
				if err != nil {
					return nil, err
				}
			}

			return asVector, nil
		}

		panic(fmt.Sprintf("FilterStep should have returned []int or []string, but returned %v (%T)", foundKeyThings, foundKeyThings))

	default:
		panic(fmt.Sprintf("Unknown path step type %T", thisStep))
	}
}

func patchFoundVectorValue(dest []any, index any, existingValue any, existed bool, remainingSteps Path, patchValue PatchFunc) ([]any, error) {
	idx, ok := index.(int)
	if !ok {
		panic("VectorStep did not return an int index as first return value.")
	}
	if idx < 0 {
		return nil, fmt.Errorf("invalid index %d: %w", idx, indexOutOfBoundsErr)
	}

	patched, err := patch(existingValue, idx, existed, remainingSteps, patchValue)
	if err != nil {
		return nil, err
	}

	// expand destination to make room for the target index
	for len(dest) < idx+1 {
		dest = append(dest, nil)
	}

	dest[idx] = patched

	return dest, nil
}

func patchFoundObjectValue(dest map[string]any, anyKey any, existingValue any, existed bool, remainingSteps Path, patchValue PatchFunc) (map[string]any, error) {
	key, ok := anyKey.(string)
	if !ok {
		panic("ObjectStep did not return a string key as first return value.")
	}

	patched, err := patch(existingValue, key, existed, remainingSteps, patchValue)
	if err != nil {
		return nil, err
	}

	dest[key] = patched

	return dest, nil
}

func setStructField(dest any, fieldName string, newValue any) error {
	rDest := unpointer(dest)

	if !rDest.CanSet() {
		return fmt.Errorf("cannot set field in %T (must call this function with a pointer)", dest)
	}

	rFieldValue := rDest.FieldByName(fieldName)
	if rFieldValue == (reflect.Value{}) || !rFieldValue.CanInterface() {
		return fmt.Errorf("no such field: %q", fieldName)
	}

	// update the value, including auto-pointer and auto-dereferencing magic
	if err := setReflectValueAdjusted(rFieldValue, newValue); err != nil {
		return err
	}

	return nil
}

func setListItem(dest any, index int, newValue any) (any, error) {
	if index < 0 {
		return nil, fmt.Errorf("invalid index %d: %w", index, indexOutOfBoundsErr)
	}

	rDest := unpointer(dest)

	if !rDest.CanSet() {
		return nil, fmt.Errorf("cannot set field in %T (must call this function with a pointer)", dest)
	}

	rNewValue, err := adjustPointerType(newValue, rDest.Type().Elem())
	if err != nil {
		return nil, err
	}

	// pad list to contain enough elements; this only works for slices
	if missing := (index + 1) - rDest.Cap(); missing > 0 {
		if rDest.Kind() != reflect.Slice {
			return nil, fmt.Errorf("invalid index %d: %w", index, indexOutOfBoundsErr)
		}

		// extend slice capacity
		rDest.Grow(missing)

		// fill-in zero values
		zeroVal := reflect.New(rDest.Type().Elem()).Elem()
		for i := 0; i < missing; i++ {
			rDest = reflect.Append(rDest, zeroVal)
		}
	}

	// update the value, including auto-pointer and auto-dereferencing magic
	rDest.Index(index).Set(*rNewValue)

	return rDest.Interface(), nil
}

func setMapItem(dest any, key any, newValue any) error {
	rDest := unpointer(dest)

	if !rDest.CanSet() {
		return fmt.Errorf("cannot set field in %T (must call this function with a pointer)", dest)
	}

	// adjust given key to the key type of the map
	rKey, err := adjustPointerType(key, rDest.Type().Key())
	if err != nil {
		return err
	}

	// adjust given value to the value type of the map
	rNewValue, err := adjustPointerType(newValue, rDest.Type().Elem())
	if err != nil {
		return err
	}

	rDest.SetMapIndex(*rKey, *rNewValue)

	return nil
}

func unpointer(value any) reflect.Value {
	rValue := reflect.ValueOf(value)

	// if it's a pointer, resolve its value
	if rValue.Kind() == reflect.Ptr {
		rValue = reflect.Indirect(rValue)
	}

	if rValue.Kind() == reflect.Interface {
		rValue = rValue.Elem()
	}

	return rValue
}

func adjustPointerType(value any, dest reflect.Type) (*reflect.Value, error) {
	rValue := reflect.ValueOf(value)

	// auto pointer handling: automatically convert from pointer to non-pointer

	// for better error message
	fieldType := dest.String()
	originalGivenType := "nil"
	if value != nil {
		originalGivenType = rValue.Type().String()
	}

	switch dest.Kind() {
	case reflect.Ptr:
		// turn untyped nils into typed ones
		if value == nil {
			rValue = reflect.New(dest).Elem()
		}

		// given value is not a pointer, so let's turn it into one
		if rValue.Kind() != reflect.Ptr {
			v := reflect.New(rValue.Type())
			v.Elem().Set(rValue)

			rValue = v
		}

	case reflect.Interface:
		// turn untyped nils into typed ones
		if value == nil {
			rValue = reflect.New(dest).Elem()
		}

	default:
		// catch untyped pointers (literal nils)
		if value == nil {
			return nil, errors.New("cannot set to null")
		}

		// given value is a pointer
		if rValue.Kind() == reflect.Ptr {
			if rValue.IsNil() {
				return nil, errors.New("cannot set to null")
			}

			// dereference the pointer
			rValue = rValue.Elem()
		}
	}

	if !rValue.Type().AssignableTo(dest) {
		return nil, fmt.Errorf("cannot set %s to %s", fieldType, originalGivenType)
	}

	return &rValue, nil
}

func setReflectValueAdjusted(dest reflect.Value, newValue any) error {
	rNewValue, err := adjustPointerType(newValue, dest.Type())
	if err != nil {
		return err
	}

	dest.Set(*rNewValue)

	return nil
}
