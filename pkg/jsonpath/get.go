// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import (
	"errors"
)

func Get(value any, path Path) (any, error) {
	if !path.IsValid() {
		return nil, errors.New("invalid path")
	}

	if len(path) == 0 {
		return value, nil
	}

	if path.HasFilterSteps() {
		return getFiltered(value, path)
	}

	return getSingle(value, path)
}

func getSingle(value any, path Path) (any, error) {
	for _, step := range path {
		_, newValue, err := traverseStep(value, step.(SingleStep))
		if err != nil {
			return nil, err
		}

		value = newValue
	}

	return value, nil
}

func getFiltered(value any, path Path) ([]any, error) {
	currentLeafValues := []any{value}

	for _, step := range path {
		newLeafValues := []any{}

		for _, val := range currentLeafValues {
			_, result, err := traverseStep(val, step)
			if err != nil {
				if ignoreErrorInFilters(err) {
					continue
				}

				return nil, err
			}

			if isFilterStep(step) {
				newValues, ok := result.([]any)
				if !ok {
					panic("isFilterStep is out of sync with path.HasFilterSteps()")
				}

				newLeafValues = append(newLeafValues, newValues...)
			} else {
				newLeafValues = append(newLeafValues, result)
			}
		}

		currentLeafValues = newLeafValues
	}

	return currentLeafValues, nil
}
