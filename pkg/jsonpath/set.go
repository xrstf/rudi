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
		case kindList:
			idx, ok := foundKeyThings.(int)
			if !ok {
				panic(fmt.Sprintf("Slice/array key is not an integer, but %T?", foundKeyThings))
			}

			patched, err := patch(foundValueThings, idx, err == nil, remainingSteps, patchValue)
			if err != nil {
				return nil, err
			}

			return setListItem(dest, idx, patched)

		case kindMap:
			patched, err := patch(foundValueThings, foundKeyThings, err == nil, remainingSteps, patchValue)
			if err != nil {
				return nil, err
			}

			return setMapItem(dest, foundKeyThings, patched)

		case kindStuct:
			fieldName, ok := foundKeyThings.(string)
			if !ok {
				panic(fmt.Sprintf("Struct field name is not a string, but %T?", foundKeyThings))
			}

			patched, err := patch(foundValueThings, foundKeyThings, err == nil, remainingSteps, patchValue)
			if err != nil {
				return nil, err
			}

			return setStructField(dest, fieldName, patched)

		default:
			panic(fmt.Sprintf("SingleStep returned unimplemented destination kind %v", destKind))
		}

	// $var[?(…)]
	case FilterStep:
		// this step found nothing, so there is no values to be updated and we can stop
		foundValues := foundValueThings.([]any)
		if len(foundValues) == 0 {
			return dest, nil
		}

		switch destKind {
		// arrays, slices
		case kindList:
			foundIndexes, ok := foundKeyThings.([]int)
			if !ok {
				panic(fmt.Sprintf("Slice/array keys are not []int, but %T?", foundKeyThings))
			}

			for idx, listIndex := range foundIndexes {
				var err error
				dest, err = patchFoundListItem(dest, listIndex, foundValues[idx], true, remainingSteps, patchValue)
				if err != nil {
					return nil, err
				}
			}

			return dest, nil

		case kindMap:
			foundKeys, ok := foundKeyThings.([]any)
			if !ok {
				panic(fmt.Sprintf("Map keys are not []any, but %T?", foundKeyThings))
			}

			for idx, key := range foundKeys {
				var err error
				dest, err = patchFoundMapItem(dest, key, foundValues[idx], true, remainingSteps, patchValue)
				if err != nil {
					return nil, err
				}
			}

			return dest, nil

		default:
			panic(fmt.Sprintf("FilterStep returned unimplemented destination kind %v", destKind))
		}

	default:
		panic(fmt.Sprintf("Unknown path step type %T", thisStep))
	}
}

func patchFoundListItem(dest any, index int, existingValue any, existed bool, remainingSteps Path, patchValue PatchFunc) (any, error) {
	if index < 0 {
		panic(fmt.Sprintf("Found negative index %d in slice?", index))
	}

	patched, err := patch(existingValue, index, existed, remainingSteps, patchValue)
	if err != nil {
		return nil, err
	}

	return setListItem(dest, index, patched)
}

func patchFoundMapItem(dest any, key any, existingValue any, existed bool, remainingSteps Path, patchValue PatchFunc) (any, error) {
	patched, err := patch(existingValue, key, existed, remainingSteps, patchValue)
	if err != nil {
		return nil, err
	}

	return setMapItem(dest, key, patched)
}

func setStructField(dest any, fieldName string, newValue any) (any, error) {
	rDest, wasPointer := unpointer(dest)

	// fmt.Printf("input rDest: %s\n", rValueString(rDest))

	if !rDest.CanSet() {
		return nil, fmt.Errorf("cannot set field in %T (must call this function with a pointer)", dest)
	}

	rFieldValue := rDest.FieldByName(fieldName)
	if rFieldValue == (reflect.Value{}) || !rFieldValue.CanInterface() {
		return nil, fmt.Errorf("no such field: %q", fieldName)
	}

	// fmt.Printf("field value: %s\n", rValueString(rFieldValue))

	// update the value, including auto-pointer and auto-dereferencing magic
	if err := setReflectValueAdjusted(rFieldValue, newValue); err != nil {
		return nil, err
	}

	if wasPointer {
		rDest = rDest.Addr()
	}

	return rDest.Interface(), nil
}

func setListItem(dest any, index int, newValue any) (any, error) {
	if index < 0 {
		return nil, fmt.Errorf("invalid index %d: %w", index, indexOutOfBoundsErr)
	}

	rDest, wasPointer := unpointer(dest)

	rNewValue, err := adjustPointerType(newValue, rDest.Type().Elem())
	if err != nil {
		return nil, err
	}

	// pad list to contain enough elements; this only works for slices;
	// this creates a completely new slice
	if missing := (index + 1) - rDest.Len(); missing > 0 {
		if rDest.Kind() != reflect.Slice {
			return nil, fmt.Errorf("invalid index %d: %w", index, indexOutOfBoundsErr)
		}

		totalLength := rDest.Len() + missing

		newSlice := reflect.MakeSlice(rDest.Type(), totalLength, totalLength)
		reflect.Copy(newSlice, rDest)

		rDest = newSlice
	}

	if rDest.Kind() == reflect.Array && !rDest.CanSet() {
		return nil, errors.New("arrays must be passed as pointers")
	}

	// update the value, including auto-pointer and auto-dereferencing magic
	rDest.Index(index).Set(*rNewValue)

	if wasPointer {
		rDest = rDest.Addr()
	}

	return rDest.Interface(), nil
}

func setMapItem(dest any, key any, newValue any) (any, error) {
	rDest, wasPointer := unpointer(dest)

	// adjust given key to the key type of the map
	rKey, err := adjustPointerType(key, rDest.Type().Key())
	if err != nil {
		return nil, err
	}

	// adjust given value to the value type of the map
	rNewValue, err := adjustPointerType(newValue, rDest.Type().Elem())
	if err != nil {
		return nil, err
	}

	rDest.SetMapIndex(*rKey, *rNewValue)

	if wasPointer {
		rDest = rDest.Addr()
	}

	return rDest.Interface(), nil
}

func rValueString(rv reflect.Value) string {
	return fmt.Sprintf("rvalue{kind=%v, type=%v, canSet=%v, canInterface=%v}", rv.Kind(), rv.Type().String(), rv.CanSet(), rv.CanInterface())
}

func rSliceValueString(rv reflect.Value) string {
	return fmt.Sprintf("rvalue{kind=%v, type=%v, cap=%d, len=%d, canSet=%v, canInterface=%v}", rv.Kind(), rv.Type().String(), rv.Cap(), rv.Len(), rv.CanSet(), rv.CanInterface())
}

func unpointer(value any) (reflect.Value, bool) {
	rValue := reflect.ValueOf(&value)
	isPointer := reflect.ValueOf(value).Kind() == reflect.Ptr

	// fmt.Printf("input rvalue: %s\n", rValueString(rValue))

	rValue = reflect.Indirect(rValue)
	// fmt.Printf(" -> unpointered : %s\n", rValueString(rValue))

	if rValue.Kind() == reflect.Interface {
		rValue = rValue.Elem()
		// fmt.Printf(" -> uninterfaced: %s\n", rValueString(rValue))

		var field reflect.Value
		if rValue.Kind() == reflect.Ptr {
			field = reflect.New(rValue.Elem().Type())
			// fmt.Printf(" -> new field pre: %s\n", rValueString(field))
			field.Elem().Set(rValue.Elem())
		} else {
			field = reflect.New(rValue.Type())
			field.Elem().Set(rValue)
		}
		// fmt.Printf(" -> new field   : %s\n", rValueString(field))

		rValue = field.Elem()
		// fmt.Printf(" -> new rValue  : %s\n", rValueString(field.Elem()))
	}

	// fmt.Printf("final input rvalue: %s\n", rValueString(rValue))

	return rValue, isPointer
}

// func unpointer(value any) (reflect.Value, bool) {
// 	return unpointerValue(reflect.ValueOf(&value))
// }

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
