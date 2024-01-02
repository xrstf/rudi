// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

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

	steps, err := traverse(value, path)
	if err != nil {
		return nil, err
	}

	// This should never happen because of the earlier check.
	if len(steps) == 0 {
		return value, nil
	}

	lastStep := steps[len(steps)-1]

	return *lastStep, nil
}
