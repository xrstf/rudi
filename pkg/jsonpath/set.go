// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import (
	"errors"
	"fmt"
)

type ObjectWriter interface {
	ObjectReader
	SetObjectKey(name string, value any) (any, error)
}

type VectorWriter interface {
	VectorReader
	SetVectorItem(index int, value any) (any, error)
}

func Set(dest any, path Path, newValue any) (any, error) {
	if len(path) == 0 {
		return newValue, nil
	}

	thisStep := path[0]
	remainingSteps := path[1:]

	// [index]...
	if index, ok := toIntegerStep(thisStep); ok {
		if slice, ok := dest.([]any); ok {
			if index < 0 || index >= len(slice) {
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

		if writer, ok := dest.(VectorWriter); ok {
			existingValue, err := writer.GetVectorItem(index)
			if err != nil {
				return nil, fmt.Errorf("cannot descend with [%d] into %T", index, dest)
			}

			updatedValue, err := Set(existingValue, remainingSteps, newValue)
			if err != nil {
				return nil, err
			}

			return writer.SetVectorItem(index, updatedValue)
		}

		return nil, fmt.Errorf("cannot descend with [%d] into %T", index, dest)
	}

	// .key
	if key, ok := toStringStep(thisStep); ok {
		if object, ok := dest.(map[string]any); ok {
			// getting the empty value for non-existing keys is fine
			existingValue := object[key]

			updatedValue, err := Set(existingValue, remainingSteps, newValue)
			if err != nil {
				return nil, err
			}

			object[key] = updatedValue

			return object, nil
		}

		if writer, ok := dest.(ObjectWriter); ok {
			existingValue, err := writer.GetObjectKey(key)
			if err != nil {
				return nil, fmt.Errorf("cannot descend with [%s] into %T", key, dest)
			}

			updatedValue, err := Set(existingValue, remainingSteps, newValue)
			if err != nil {
				return nil, err
			}

			return writer.SetObjectKey(key, updatedValue)
		}

		// nulls can be turned into objects
		if dest == nil {
			updatedValue, err := Set(nil, remainingSteps, newValue)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				key: updatedValue,
			}, nil
		}

		return nil, fmt.Errorf("cannot descend with [%s] into %T", key, dest)
	}

	return nil, errors.New("invalid path step: neither key nor index")
}
