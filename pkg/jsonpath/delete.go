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
	// $var[1], $var.foo, $var["foo"], $var[(+ 1 2)]
	case SingleStep:
		switch keyThing := foundKeyThings.(type) {
		case int:
			asVector, ok := dest.([]any)
			if !ok {
				panic("VectorStep should have errored on a non-vector value.")
			}

			deleted, err := deleteInternal(foundValueThings, remainingSteps)
			if err != nil {
				return nil, err
			}

			asVector[keyThing] = deleted

			return asVector, nil

		case string:
			asObject, ok := dest.(map[string]any)
			if !ok {
				panic("ObjectStep should have errored on a non-object value.")
			}

			deleted, err := deleteInternal(foundValueThings, remainingSteps)
			if err != nil {
				return nil, err
			}

			asObject[keyThing] = deleted

			return asObject, nil

		default:
			panic(fmt.Sprintf("SingleStep should have returned int index or string key, but returned %v (%T)", foundKeyThings, foundKeyThings))
		}

	// $var[?(…)]
	case FilterStep:
		foundValues := foundValueThings.([]any)
		if len(foundValues) == 0 {
			return dest, nil
		}

		foundsKeys, ok := foundKeyThings.([]string)
		if ok {
			asObject, ok := dest.(map[string]any)
			if !ok {
				panic("ObjectStep should have errored on a non-object value.")
			}

			for idx, key := range foundsKeys {
				deleted, err := deleteInternal(foundValues[idx], remainingSteps)
				if err != nil {
					return nil, err
				}

				asObject[key] = deleted
			}

			return asObject, nil
		}

		foundIndexes, ok := foundKeyThings.([]int)
		if ok {
			asVector, ok := dest.([]any)
			if !ok {
				panic("VectorStep should have errored on a non-vector value.")
			}

			for idx, vectorIndex := range foundIndexes {
				deleted, err := deleteInternal(foundValues[idx], remainingSteps)
				if err != nil {
					return nil, err
				}

				asVector[vectorIndex] = deleted
			}

			return asVector, nil
		}

		panic(fmt.Sprintf("FilterStep should have returned []int or []string, but returned %v (%T)", foundKeyThings, foundKeyThings))

	default:
		panic(fmt.Sprintf("Unknown path step type %T", thisStep))
	}
}

func deleteFromLeaf(dest any, step Step, foundKeyThings any) (any, error) {
	switch step.(type) {
	case SingleStep:
		switch keyThing := foundKeyThings.(type) {
		case int:
			asVector, ok := dest.([]any)
			if !ok {
				panic("SingleStep should have errored on a non-vector value.")
			}

			return removeSliceItem(asVector, keyThing), nil

		case string:
			asObject, ok := dest.(map[string]any)
			if !ok {
				panic("SingleStep should have errored on a non-object value.")
			}

			delete(asObject, keyThing)
			return asObject, nil

		default:
			panic(fmt.Sprintf("SingleStep should have returned int index or string key, but returned %v (%T)", foundKeyThings, foundKeyThings))
		}

	default:
		panic(fmt.Sprintf("Unknown path step type %T", step))
	}
}
