// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package equality

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

var ErrIncompatibleTypes = errors.New("types are incompatible")

func StrictEqual(left, right ast.Literal) (bool, error) {
	switch leftAsserted := left.(type) {
	case ast.Null:
		return nullStrictEquals(leftAsserted, right)
	case ast.Bool:
		return boolStrictEquals(leftAsserted, right)
	case ast.String:
		return stringStrictEquals(leftAsserted, right)
	case ast.Number:
		return numberStrictEquals(leftAsserted, right)
	case ast.Vector:
		return vectorStrictEquals(leftAsserted, right)
	case ast.Object:
		return objectStrictEquals(leftAsserted, right)
	default:
		return false, fmt.Errorf("cannot compare with %T with %T", left, right)
	}
}

func boolStrictEquals(left ast.Bool, right ast.Literal) (bool, error) {
	rightValue, ok := right.(ast.Bool)
	if !ok {
		return false, ErrIncompatibleTypes
	}

	return left.Equal(rightValue), nil
}

func nullStrictEquals(left ast.Null, right any) (bool, error) {
	rightValue, ok := right.(ast.Null)
	if !ok {
		return false, ErrIncompatibleTypes
	}

	return left.Equal(rightValue), nil
}

func stringStrictEquals(left ast.String, right any) (bool, error) {
	rightValue, ok := right.(ast.String)
	if !ok {
		return false, ErrIncompatibleTypes
	}

	return left.Equal(rightValue), nil
}

func numberStrictEquals(left ast.Number, right any) (bool, error) {
	rightValue, ok := right.(ast.Number)
	if !ok {
		return false, ErrIncompatibleTypes
	}

	return left.Equal(rightValue), nil
}

func vectorStrictEquals(left ast.Vector, right any) (bool, error) {
	rightValue, ok := right.(ast.Vector)
	if !ok {
		return false, ErrIncompatibleTypes
	}

	if len(left.Data) != len(rightValue.Data) {
		return false, nil
	}

	for i, leftItem := range left.Data {
		rightItem := rightValue.Data[i]

		leftWrapped, err := types.WrapNative(leftItem)
		if err != nil {
			return false, ErrIncompatibleTypes
		}

		rightWrapped, err := types.WrapNative(rightItem)
		if err != nil {
			return false, ErrIncompatibleTypes
		}

		// wrapping always returns literals, so type assertions are safe here
		equal, err := StrictEqual(leftWrapped, rightWrapped)
		if err != nil {
			return false, err
		}

		if !equal {
			return false, nil
		}
	}

	return true, nil
}

func objectStrictEquals(left ast.Object, right any) (bool, error) {
	rightValue, ok := right.(ast.Object)
	if !ok {
		return false, ErrIncompatibleTypes
	}

	if len(left.Data) != len(rightValue.Data) {
		return false, nil
	}

	keysSeen := map[string]struct{}{}

	for key, leftItem := range left.Data {
		rightItem, exists := rightValue.Data[key]
		if !exists {
			return false, nil
		}

		keysSeen[key] = struct{}{}

		leftWrapped, err := types.WrapNative(leftItem)
		if err != nil {
			return false, ErrIncompatibleTypes
		}

		rightWrapped, err := types.WrapNative(rightItem)
		if err != nil {
			return false, ErrIncompatibleTypes
		}

		// wrapping always returns literals, so type assertions are safe here
		equal, err := StrictEqual(leftWrapped, rightWrapped)
		if err != nil {
			return false, err
		}

		if !equal {
			return false, nil
		}
	}

	for key := range rightValue.Data {
		delete(keysSeen, key)
	}

	return len(keysSeen) == 0, nil
}
