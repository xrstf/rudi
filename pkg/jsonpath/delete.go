// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import (
	"errors"
	"fmt"
)

type ObjectKeyDeleter interface {
	DeleteObjectKey(name string) (any, error)
}

type VectorItemDeleter interface {
	DeleteVectorItem(index int) (any, error)
}

func removeSliceItem(slice []any, index int) []any {
	return append(slice[:index], slice[index+1:]...)
}

func Delete(dest any, path Path) (any, error) {
	if !path.IsValid() {
		return nil, errors.New("invalid path")
	}

	if len(path) == 0 {
		return nil, nil
	}

	return deleteInternal(dest, path)
}

func deleteInternal(dest any, path Path) (any, error) {
	thisStep := path[0]
	remainingSteps := path[1:]

	foundKeyThings, foundValueThings, err := traverseSingleStep(dest, thisStep)
	if err != nil {
		if errors.Is(err, noSuchKeyErr) || errors.Is(err, indexOutOfBoundsErr) {
			return dest, nil
		}

		return nil, err
	}

	// we reached the level at which we want to remove the key
	if len(remainingSteps) == 0 {
		return deleteFromLeaf(dest, thisStep, foundKeyThings)
	}

	switch thisStep.(type) {
	// $var[…]
	case SingularVectorStep:
		asVector, ok := dest.([]any)
		if !ok {
			panic("VectorStep should have errored on a non-vector value.")
		}

		index, ok := foundKeyThings.(int)
		if !ok {
			panic("VectorStep should have returned an int index.")
		}

		deleted, err := deleteInternal(foundValueThings, remainingSteps)
		if err != nil {
			return nil, err
		}

		asVector[index] = deleted

		return asVector, nil

	// $var[?(…)]
	case MultiVectorStep:
		foundValues := foundValueThings.([]any)
		if len(foundValues) == 0 {
			return dest, nil
		}

		asVector, ok := dest.([]any)
		if !ok {
			panic("VectorStep should have errored on a non-vector value.")
		}

		for idx, vectorIndex := range foundKeyThings.([]int) {
			deleted, err := deleteInternal(foundValues[idx], remainingSteps)
			if err != nil {
				return nil, err
			}

			asVector[vectorIndex] = deleted
		}

		return asVector, nil

	// $var.…
	case SingularObjectStep:
		asObject, ok := dest.(map[string]any)
		if !ok {
			panic("ObjectStep should have errored on a non-object value.")
		}

		key, ok := foundKeyThings.(string)
		if !ok {
			panic("ObjectStep should have returned an string key.")
		}

		deleted, err := deleteInternal(foundValueThings, remainingSteps)
		if err != nil {
			return nil, err
		}

		asObject[key] = deleted

		return asObject, nil

	// $var[?(…)]
	case MultiObjectStep:
		foundValues := foundValueThings.([]any)
		if len(foundValues) == 0 {
			return dest, nil
		}

		asObject, ok := dest.(map[string]any)
		if !ok {
			panic("ObjectStep should have errored on a non-object value.")
		}

		for idx, key := range foundKeyThings.([]string) {
			deleted, err := deleteInternal(foundValues[idx], remainingSteps)
			if err != nil {
				return nil, err
			}

			asObject[key] = deleted
		}

		return asObject, nil

	default:
		panic(fmt.Sprintf("Unknown path step type %T", thisStep))
	}
}

func deleteFromLeaf(dest any, step Step, foundKeyThings any) (any, error) {
	switch step.(type) {
	case SingularVectorStep:
		asVector, ok := dest.([]any)
		if !ok {
			panic("VectorStep should have errored on a non-vector value.")
		}

		index, ok := foundKeyThings.(int)
		if !ok {
			panic("VectorStep should have returned an int index.")
		}

		return removeSliceItem(asVector, index), nil

	case SingularObjectStep:
		asObject, ok := dest.(map[string]any)
		if !ok {
			panic("ObjectStep should have errored on a non-object value.")
		}

		key, ok := foundKeyThings.(string)
		if !ok {
			panic("ObjectStep should have returned an string key.")
		}

		delete(asObject, key)
		return asObject, nil

	default:
		panic(fmt.Sprintf("Unknown path step type %T", step))
	}
}
