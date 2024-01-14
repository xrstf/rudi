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

	if valueAsObject, ok := value.(map[string]any); ok {
		return traverseObjectSingleStep(valueAsObject, step)
	}

	if value == nil {
		switch s := step.(type) {
		case SingularVectorStep:
			return s.Index(), nil, indexOutOfBoundsErr
		case MultiVectorStep:
			return []int{}, []any{}, nil
		case SingularObjectStep:
			return s.Key(), nil, noSuchKeyErr
		case MultiObjectStep:
			return []string{}, []any{}, nil
		}
	}

	return nil, nil, fmt.Errorf("cannot traverse %T: %w", value, untraversableErr)
}

func traverseVectorSingleStep(value []any, step Step) (any, any, error) {
	if vectorStep, ok := step.(SingularVectorStep); ok {
		index := vectorStep.Index()
		if index >= len(value) {
			return index, nil, fmt.Errorf("invalid index %d: %w", index, indexOutOfBoundsErr)
		}

		// this is not an out of bounds because negative indexes should not be silently swallowed
		if index < 0 {
			return index, nil, fmt.Errorf("invalid index %d", index)
		}

		return index, value[index], nil
	}

	if vectorStep, ok := step.(MultiVectorStep); ok {
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
	if objStep, ok := step.(SingularObjectStep); ok {
		key := objStep.Key()

		val, exists := value[key]
		if !exists {
			return key, nil, fmt.Errorf("invalid key %q: %w", key, noSuchKeyErr)
		}

		return key, val, nil
	}

	if objStep, ok := step.(MultiObjectStep); ok {
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
