// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import (
	"errors"
	"fmt"
	"sort"
)

var (
	noSuchKeyErr        = errors.New("no such key")
	indexOutOfBoundsErr = errors.New("index out of bounds")
	invalidStepErr      = errors.New("cannot use this step type to traverse")
	untraversableErr    = errors.New("does not support traversing into this type")
)

func ignoreErrorInDynamic(err error) bool {
	// invalidStepErr are not silently swallowed!

	return errors.Is(err, noSuchKeyErr) || errors.Is(err, indexOutOfBoundsErr) || errors.Is(err, untraversableErr)
}

func traverseSingleStep(value any, step Step) (any, any, error) {
	if valueAsVector, ok := value.([]any); ok {
		return traverseVectorSingleStep(valueAsVector, step)
	}

	// if vectorReader, ok := value.(VectorReader); ok {
	// 	index, ok := toIntegerStep(step)
	// 	if ok {
	// 		var err error

	// 		value, err = vectorReader.GetVectorItem(index)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("cannot descend with %v (%T) into %T: %w", step, step, value, err)
	// 		}

	// 		continue
	// 	}
	// }

	if valueAsObject, ok := value.(map[string]any); ok {
		return traverseObjectSingleStep(valueAsObject, step)
	}

	// if objectReader, ok := value.(ObjectReader); ok {
	// 	key, ok := toStringStep(step)
	// 	if ok {
	// 		var err error

	// 		value, err = objectReader.GetObjectKey(key)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("cannot descend with %v (%T) into %T: %w", step, step, value, err)
	// 		}

	// 		continue
	// 	}
	// }

	return nil, nil, fmt.Errorf("cannot traverse %T: %w", value, untraversableErr)
}

func traverseVectorSingleStep(value []any, step Step) (any, any, error) {
	if vectorStep, ok := step.(VectorStep); ok {
		index := vectorStep.Index()
		if index < 0 || index >= len(value) {
			return nil, nil, fmt.Errorf("invalid index %d: %w", index, indexOutOfBoundsErr)
		}

		return index, value[index], nil
	}

	if vectorStep, ok := step.(DynamicVectorStep); ok {
		indexes := []int{}
		values := []any{}

		for index, val := range value {
			keep, err := vectorStep.Keep(index, val)
			if err != nil {
				return nil, nil, err
			}

			if keep {
				indexes = append(indexes, index)
				values = append(values, val)
			}
		}

		return indexes, values, nil
	}

	return nil, nil, fmt.Errorf("invalid step %T: %w", step, invalidStepErr)
}

func traverseObjectSingleStep(value map[string]any, step Step) (any, any, error) {
	if objStep, ok := step.(ObjectStep); ok {
		key := objStep.Key()

		val, exists := value[key]
		if !exists {
			return nil, nil, fmt.Errorf("invalid key %q: %w", key, noSuchKeyErr)
		}

		return key, val, nil
	}

	if objStep, ok := step.(DynamicObjectStep); ok {
		// To allow side effects in the dynamic step to work consistently,
		// we need to loop over the object in a consistent way.
		orderedKeys := orderedObjectKeys(value)

		selectedKeys := []string{}
		selectedValues := []any{}

		for _, key := range orderedKeys {
			val := value[key]

			keep, err := objStep.Keep(key, val)
			if err != nil {
				return nil, nil, err
			}

			if keep {
				selectedKeys = append(selectedKeys, key)
				selectedValues = append(selectedValues, val)
			}
		}

		return selectedKeys, selectedValues, nil
	}

	return nil, nil, fmt.Errorf("cannot traverse %T: %w", value, untraversableErr)
}

func orderedObjectKeys(obj map[string]any) []string {
	allKeys := make([]string, len(obj))
	i := 0

	for k := range obj {
		allKeys[i] = k
		i++
	}

	// Using the slightly faster slices.Sort would bump our min Go version to 1.21.
	sort.Strings(allKeys)

	return allKeys
}
