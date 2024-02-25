// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

var (
	noSuchKeyErr        = errors.New("no such key")
	indexOutOfBoundsErr = errors.New("index out of bounds")
	invalidStepErr      = errors.New("cannot use this step type to traverse")
	untraversableErr    = errors.New("does not support traversing into this type")
)

func ignoreErrorInFilters(err error) bool {
	// invalidStepErr are not silently swallowed!
	return errors.Is(err, noSuchKeyErr) || errors.Is(err, indexOutOfBoundsErr) || errors.Is(err, untraversableErr)
}

type variableKind int

const (
	kindList variableKind = iota
	kindMap
	kindStuct
)

func traverseStep(value any, step Step) (any, any, variableKind, error) {
	if value == nil {
		return nil, nil, 0, fmt.Errorf("cannot traverse into null: %w", untraversableErr)
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
			return nil, nil, 0, errors.New("cannot descend into nil")
		}

		elemType = valueType.Elem()
		elemKind = elemType.Kind()

		// dereference the pointer
		rValue = rValue.Elem()

		// fmt.Printf("  pkitty: %v (%v)\n", elemType, elemKind)
	}

	switch elemKind {
	case reflect.Slice, reflect.Array:
		keys, results, err := traverseIndexableStep(rValue, step)
		return keys, results, kindList, err

	case reflect.Map:
		keys, results, err := traverseMapStep(rValue, step)
		return keys, results, kindMap, err

	case reflect.Struct:
		keys, results, err := traverseStructStep(rValue, step)
		return keys, results, kindStuct, err

	default:
		return nil, nil, 0, fmt.Errorf("cannot traverse %T: %w", value, untraversableErr)
	}
}

func traverseIndexableStep(value reflect.Value, step Step) (key any, result any, err error) {
	switch asserted := step.(type) {
	case SingleStep:
		return traverseIndexableSingleStep(value, asserted)
	case FilterStep:
		return traverseIndexableFilterStep(value, asserted)
	default:
		panic(fmt.Sprintf("Unknown path type %T.", step))
	}
}

func traverseIndexableSingleStep(value reflect.Value, step SingleStep) (key any, result any, err error) {
	index, ok := step.ToIndex()
	if !ok {
		return nil, nil, fmt.Errorf("cannot use step %v to traverses into vectors", step)
	}

	// this is not an out of bounds because negative indexes should not be silently swallowed
	if index < 0 {
		return index, nil, fmt.Errorf("invalid index %d", index)
	}

	if index >= value.Len() {
		return index, nil, fmt.Errorf("invalid index %d: %w", index, indexOutOfBoundsErr)
	}

	return index, value.Index(index).Interface(), nil
}

func traverseIndexableFilterStep(value reflect.Value, step FilterStep) (key any, result any, err error) {
	indexes := []int{}
	values := []any{}

	for index := 0; index < value.Len(); index++ {
		val := value.Index(index).Interface()

		keep, err := step.Keep(index, val)
		if err != nil {
			// Removing the error's type is important so further up the call chain we can distinguish
			// between "$var.foo" with .foo not existing, or $var[?(eq? .foo 1)]; otherwise too many
			// errors would be swallowed.
			return nil, nil, errors.New(err.Error())
		}

		if keep {
			indexes = append(indexes, index)
			values = append(values, val)
		}
	}

	return indexes, values, nil
}

func traverseMapStep(value reflect.Value, step Step) (key any, result any, err error) {
	switch asserted := step.(type) {
	case SingleStep:
		return traverseMapSingleStep(value, asserted)
	case FilterStep:
		return traverseMapFilterStep(value, asserted)
	default:
		panic(fmt.Sprintf("Unknown path type %T.", step))
	}
}

func traverseMapSingleStep(value reflect.Value, step SingleStep) (key any, result any, err error) {
	key, ok := step.ToKey()
	if !ok {
		return nil, nil, fmt.Errorf("cannot use step %v to traverses into objects", step)
	}

	keyValue := value.MapIndex(reflect.ValueOf(key))
	if keyValue == (reflect.Value{}) {
		return key, nil, fmt.Errorf("invalid key %q: %w", key, noSuchKeyErr)
	}

	return key, keyValue.Interface(), nil
}

func traverseMapFilterStep(value reflect.Value, step FilterStep) (key any, result any, err error) {
	// To allow side effects in the dynamic step to work consistently,
	// we need to loop over the object in a consistent way.
	orderedKeys := orderedObjectKeys(value)

	selectedKeys := []string{}
	selectedValues := []any{}

	for _, key := range orderedKeys {
		keyName := key.String()
		val := value.MapIndex(key).Interface()

		// TODO: Would a check like this be needed?
		// if keyValue == (reflect.Value{}) {
		// 	return key, nil, fmt.Errorf("invalid key %q: %w", key, noSuchKeyErr)
		// }

		keep, err := step.Keep(keyName, val)
		if err != nil {
			// Removing the error's type is important so further up the call chain we can distinguish
			// between "$var.foo" with .foo not existing, or $var[?(eq? .foo 1)]; otherwise too many
			// errors would be swallowed.
			return nil, nil, errors.New(err.Error())
		}

		if keep {
			selectedKeys = append(selectedKeys, keyName)
			selectedValues = append(selectedValues, val)
		}
	}

	return selectedKeys, selectedValues, nil
}

func orderedObjectKeys(obj reflect.Value) []reflect.Value {
	allKeys := obj.MapKeys()

	// Using the slightly faster slices.Sort would bump our min Go version to 1.21.
	sort.Slice(allKeys, func(i, j int) bool {
		return allKeys[i].String() < allKeys[j].String()
	})

	return allKeys
}

func traverseStructStep(value reflect.Value, step Step) (key any, result any, err error) {
	switch asserted := step.(type) {
	case SingleStep:
		return traverseStructSingleStep(value, asserted)
	case FilterStep:
		return traverseStructFilterStep(value, asserted)
	default:
		panic(fmt.Sprintf("Unknown path type %T.", step))
	}
}

func traverseStructSingleStep(value reflect.Value, step SingleStep) (key any, result any, err error) {
	fieldName, ok := step.ToKey()
	if !ok {
		return nil, nil, fmt.Errorf("cannot use step %v to traverses into structs", step)
	}

	fieldValue := value.FieldByName(fieldName)
	if fieldValue == (reflect.Value{}) || !fieldValue.CanInterface() {
		return fieldName, nil, fmt.Errorf("no such field: %q", fieldName)
	}

	return fieldName, fieldValue.Interface(), nil
}

func traverseStructFilterStep(value reflect.Value, step FilterStep) (key any, result any, err error) {
	// To allow side effects in the dynamic step to work consistently,
	// we need to loop over the object in a consistent way.
	orderedFieldNames := orderedFieldNames(value)

	selectedKeys := []string{}
	selectedValues := []any{}

	for _, key := range orderedFieldNames {
		val := value.FieldByName(key).Interface()

		// TODO: Would a check like this be needed?
		// if fieldValue == (reflect.Value{}) || !fieldValue.CanInterface() {
		// 	return fieldName, nil, fmt.Errorf("no such field: %q", fieldName)
		// }

		keep, err := step.Keep(key, val)
		if err != nil {
			// Removing the error's type is important so further up the call chain we can distinguish
			// between "$var.foo" with .foo not existing, or $var[?(eq? .foo 1)]; otherwise too many
			// errors would be swallowed.
			return nil, nil, errors.New(err.Error())
		}

		if keep {
			selectedKeys = append(selectedKeys, key)
			selectedValues = append(selectedValues, val)
		}
	}

	return selectedKeys, selectedValues, nil
}

func orderedFieldNames(obj reflect.Value) []string {
	objType := obj.Type()

	names := make([]string, objType.NumField())
	for i := range names {
		names[i] = objType.Field(i).Name
	}

	// Using the slightly faster slices.Sort would bump our min Go version to 1.21.
	sort.Strings(names)

	return names
}
