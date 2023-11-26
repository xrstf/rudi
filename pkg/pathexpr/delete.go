// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package pathexpr

import (
	"errors"
	"fmt"
)

func removeSliceItem(slice []any, index int) []any {
	return append(slice[:index], slice[index+1:]...)
}

func Delete(dest any, path Path) (any, error) {
	if len(path) == 0 {
		return nil, nil
	}

	thisStep := path[0]
	remainingSteps := path[1:]

	// we reached the level at which we want to remove the key
	if len(remainingSteps) == 0 {
		// [index]
		if index, ok := toIntegerStep(thisStep); ok {
			if index < 0 {
				return nil, fmt.Errorf("index %d out of bounds", index)
			}

			if slice, ok := dest.([]any); ok {
				if index >= len(slice) {
					return nil, fmt.Errorf("index %d out of bounds", index)
				}

				return removeSliceItem(slice, index), nil
			}

			return nil, fmt.Errorf("cannot delete index from %T", dest)
		}

		// .key
		if key, ok := toStringStep(thisStep); ok {
			if object, ok := dest.(map[string]any); ok {
				delete(object, key)
				return object, nil
			}

			return nil, fmt.Errorf("cannot delete key from %T", dest)
		}

		return nil, fmt.Errorf("can only remove object keys or slice items, not %T", thisStep)
	}

	// [index]...
	if index, ok := toIntegerStep(thisStep); ok {
		if index < 0 {
			return nil, fmt.Errorf("index %d out of bounds", index)
		}

		if slice, ok := dest.([]any); ok {
			if index >= len(slice) {
				return nil, fmt.Errorf("index %d out of bounds", index)
			}

			existingValue := slice[index]

			updatedValue, err := Delete(existingValue, remainingSteps)
			if err != nil {
				return nil, err
			}

			slice[index] = updatedValue

			return slice, nil
		}

		return nil, fmt.Errorf("cannot descend with [%d] into %T", index, dest)
	}

	// .key
	if key, ok := toStringStep(thisStep); ok {
		if object, ok := dest.(map[string]any); ok {
			// getting the empty value for non-existing keys is fine
			existingValue := object[key]

			updatedValue, err := Delete(existingValue, remainingSteps)
			if err != nil {
				return nil, err
			}

			object[key] = updatedValue

			return object, nil
		}

		return nil, fmt.Errorf("cannot descend with [%s] into %T", key, dest)
	}

	return nil, errors.New("invalid path step: neither key nor index")
}
