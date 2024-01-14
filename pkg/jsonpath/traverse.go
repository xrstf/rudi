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

func ignoreErrorInFilters(err error) bool {
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
		case SingleStep:
			index, key := indexOrKey(s)
			switch {
			case index != nil:
				return *index, nil, indexOutOfBoundsErr
			case key != nil:
				return *key, nil, noSuchKeyErr
			default:
				return nil, nil, fmt.Errorf("%T is neither key nor index.", s)
			}

		case FilterStep:
			return []any{}, []any{}, nil
		}
	}

	return nil, nil, fmt.Errorf("cannot traverse %T: %w", value, untraversableErr)
}

func traverseVectorSingleStep(value []any, step Step) (any, any, error) {
	if vectorStep, ok := step.(SingleStep); ok {
		index, ok := vectorStep.ToIndex()
		if !ok {
			return nil, nil, fmt.Errorf("cannot use step %v to traverses into vectors", step)
		}

		if index >= len(value) {
			return index, nil, fmt.Errorf("invalid index %d: %w", index, indexOutOfBoundsErr)
		}

		// this is not an out of bounds because negative indexes should not be silently swallowed
		if index < 0 {
			return index, nil, fmt.Errorf("invalid index %d", index)
		}

		return index, value[index], nil
	}

	if filterStep, ok := step.(FilterStep); ok {
		indexes := []int{}
		values := []any{}

		for index, val := range value {
			keep, err := filterStep.Keep(index, val)
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
	if objStep, ok := step.(SingleStep); ok {
		key, ok := objStep.ToKey()
		if !ok {
			return nil, nil, fmt.Errorf("cannot use step %v to traverses into objects", step)
		}

		val, exists := value[key]
		if !exists {
			return key, nil, fmt.Errorf("invalid key %q: %w", key, noSuchKeyErr)
		}

		return key, val, nil
	}

	if filterStep, ok := step.(FilterStep); ok {
		// To allow side effects in the dynamic step to work consistently,
		// we need to loop over the object in a consistent way.
		orderedKeys := orderedObjectKeys(value)

		selectedKeys := []string{}
		selectedValues := []any{}

		for _, key := range orderedKeys {
			val := value[key]

			keep, err := filterStep.Keep(key, val)
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
