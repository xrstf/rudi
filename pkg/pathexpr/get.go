// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package pathexpr

import (
	"fmt"
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

	for _, step := range path {
		if valueAsSlice, ok := value.([]any); ok {
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

		if valueAsObject, ok := value.(map[string]any); ok {
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

		return nil, fmt.Errorf("cannot descend with %v (%T) into %T", step, step, value)
	}

	return value, nil
}
