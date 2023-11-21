// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package pathexpr

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/eval/types"
)

func Set(dest any, path Path, newValue any) (any, error) {
	if len(path) == 0 {
		return newValue, nil
	}

	target, err := types.UnwrapType(dest)
	if err != nil {
		return nil, fmt.Errorf("cannot descend into %T", dest)
	}

	thisStep := path[0]
	remainingSteps := path[1:]

	// [index]...
	if index, ok := toIntegerStep(thisStep); ok {
		if index < 0 {
			return nil, fmt.Errorf("index %d out of bounds", index)
		}

		if slice, ok := target.([]any); ok {
			if index >= len(slice) {
				return nil, fmt.Errorf("index %d out of bounds", index)
			}

			existingValue := slice[index]

			updatedValue, err := Set(existingValue, remainingSteps, newValue)
			if err != nil {
				return nil, err
			}

			slice[index] = updatedValue

			return slice, nil
		}

		return nil, fmt.Errorf("cannot descend with [%d] into %T", index, target)
	}

	// .key
	if key, ok := toStringStep(thisStep); ok {
		if object, ok := target.(map[string]any); ok {
			// getting the empty value for non-existing keys is fine
			existingValue := object[key]

			updatedValue, err := Set(existingValue, remainingSteps, newValue)
			if err != nil {
				return nil, err
			}

			object[key] = updatedValue

			return object, nil
		}

		// nulls can be turned into objects
		if target == nil {
			updatedValue, err := Set(nil, remainingSteps, newValue)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				key: updatedValue,
			}, nil
		}

		return nil, fmt.Errorf("cannot descend with [%s] into %T", key, target)
	}

	return nil, errors.New("invalid path step: neither key nor index")
}
