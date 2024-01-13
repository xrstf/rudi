// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import (
	"errors"
)

type ObjectReader interface {
	GetObjectKey(name string) (any, error)
}

type VectorReader interface {
	GetVectorItem(index int) (any, error)
}

func Get(value any, path Path) (any, error) {
	if !path.IsValid() {
		return nil, errors.New("invalid path")
	}

	if len(path) == 0 {
		return value, nil
	}

	if path.HasMultiSteps() {
		return getDynamic(value, path)
	}

	return getSingular(value, path)
}

func getSingular(value any, path Path) (any, error) {
	for _, step := range path {
		_, newValue, err := traverseSingleStep(value, step)
		if err != nil {
			return nil, err
		}

		value = newValue
	}

	return value, nil
}

func getDynamic(value any, path Path) ([]any, error) {
	currentLeafValues := []any{value}

	for _, step := range path {
		newLeafValues := []any{}

		for _, val := range currentLeafValues {
			_, result, err := traverseSingleStep(val, step)
			if err != nil {
				if ignoreErrorInDynamic(err) {
					continue
				}

				return nil, err
			}

			if isMultiStep(step) {
				newValues, ok := result.([]any)
				if !ok {
					panic("isDynamicStep is out of sync with path.IsDynamic()")
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
