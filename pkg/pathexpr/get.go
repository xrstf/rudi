// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package pathexpr

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/eval/types"
)

func Get(value any, path Path) (any, error) {
	if len(path) == 0 {
		return value, nil
	}

	for _, step := range path {
		nativeVal, err := types.UnwrapType(value)
		if err != nil {
			return nil, fmt.Errorf("cannot descend with %s into %T: %w", step, value, err)
		}

		if valueAsSlice, ok := nativeVal.([]any); ok {
			index, ok := toIntegerStep(step)
			if !ok {
				return nil, fmt.Errorf("cannot use %v as an array index", step)
			}

			if index < 0 || index >= len(valueAsSlice) {
				return nil, fmt.Errorf("index %d out of bounds", index)
			}

			value = valueAsSlice[index]
			continue
		}

		if valueAsObject, ok := nativeVal.(map[string]any); ok {
			key, ok := toStringStep(step)
			if !ok {
				return nil, fmt.Errorf("cannot use %v as an object key", step)
			}

			var exists bool
			value, exists = valueAsObject[key]
			if !exists {
				return nil, fmt.Errorf("no such key: %q", key)
			}

			continue
		}

		return nil, fmt.Errorf("cannot descend with %v (%T) into %T", step, step, value)
	}

	return value, nil
}
