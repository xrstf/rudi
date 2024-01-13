// SPDX-FileCopyrightText: 2024 Christoph Mewes
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
	return Patch(dest, path, func(_ bool, _ any) (any, error) {
		return newValue, nil
	})
}

type PatchFunc func(exists bool, val any) (any, error)

func Patch(dest any, path Path, patchValue PatchFunc) (any, error) {
	if !path.IsValid() {
		return nil, errors.New("invalid path")
	}

	return patch(dest, true, path, patchValue)
}

func patch(dest any, exists bool, path Path, patchValue PatchFunc) (any, error) {
	if len(path) == 0 {
		return patchValue(exists, dest)
	}

	thisStep := path[0]
	remainingSteps := path[1:]

	foundKeyThings, foundValueThings, err := traverseSingleStep(dest, thisStep)
	if err != nil && !errors.Is(err, noSuchKeyErr) && !errors.Is(err, indexOutOfBoundsErr) {
		return nil, err
	}

	switch thisStep.(type) {
	// $var[…]
	case SingularVectorStep:
		// nil values (or non-existing values) can be turned into vectors
		if dest == nil {
			dest = []any{}
		}

		asVector, ok := dest.([]any)
		if !ok {
			panic("VectorStep should have errored on a non-vector value.")
		}

		return patchFoundVectorValue(asVector, foundKeyThings, foundValueThings, err == nil, remainingSteps, patchValue)

	// $var[?(…)]
	case MultiVectorStep:
		// nil values (or non-existing values) can be turned into vectors
		if dest == nil {
			dest = []any{}
		}

		asVector, ok := dest.([]any)
		if !ok {
			panic("VectorStep should have errored on a non-vector value.")
		}

		foundValues := foundValueThings.([]any)

		for idx, vectorIndex := range foundKeyThings.([]int) {
			var err error
			asVector, err = patchFoundVectorValue(asVector, vectorIndex, foundValues[idx], true, remainingSteps, patchValue)
			if err != nil {
				return nil, err
			}
		}

		return asVector, nil

	// $var.…
	case SingularObjectStep:
		// nil values (or non-existing values) can be turned into objects
		if dest == nil {
			dest = map[string]any{}
		}

		asObject, ok := dest.(map[string]any)
		if !ok {
			panic("ObjectStep should have errored on a non-object value.")
		}

		return patchFoundObjectValue(asObject, foundKeyThings, foundValueThings, err == nil, remainingSteps, patchValue)

	// $var[?(…)]
	case MultiObjectStep:
		// nil values (or non-existing values) can be turned into objects
		if dest == nil {
			dest = map[string]any{}
		}

		asObject, ok := dest.(map[string]any)
		if !ok {
			panic("ObjectStep should have errored on a non-object value.")
		}

		foundValues := foundValueThings.([]any)

		for idx, key := range foundKeyThings.([]string) {
			var err error
			asObject, err = patchFoundObjectValue(asObject, key, foundValues[idx], true, remainingSteps, patchValue)
			if err != nil {
				return nil, err
			}
		}

		return asObject, nil

	default:
		panic(fmt.Sprintf("Unknown path step type %T", thisStep))
	}
}

func patchFoundVectorValue(dest []any, index any, existingValue any, existed bool, remainingSteps Path, patchValue PatchFunc) ([]any, error) {
	idx, ok := index.(int)
	if !ok {
		panic("VectorStep did not return an int index as first return value.")
	}
	if idx < 0 {
		return nil, fmt.Errorf("invalid index %d: %w", idx, indexOutOfBoundsErr)
	}

	patched, err := patch(existingValue, existed, remainingSteps, patchValue)
	if err != nil {
		return nil, err
	}

	// expand destination to make room for the target index
	for len(dest) < idx+1 {
		dest = append(dest, nil)
	}

	dest[idx] = patched

	return dest, nil
}

func patchFoundObjectValue(dest map[string]any, anyKey any, existingValue any, existed bool, remainingSteps Path, patchValue PatchFunc) (map[string]any, error) {
	key, ok := anyKey.(string)
	if !ok {
		panic("ObjectStep did not return a string key as first return value.")
	}

	patched, err := patch(existingValue, existed, remainingSteps, patchValue)
	if err != nil {
		return nil, err
	}

	dest[key] = patched

	return dest, nil
}
