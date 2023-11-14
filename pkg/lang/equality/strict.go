// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package equality

import (
	"errors"
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
)

var ErrIncompatibleTypes = errors.New("types are incompatible")

func StrictEqual(left, right any) (bool, error) {
	switch leftAsserted := left.(type) {
	case ast.Null:
		return nullStrictEquals(leftAsserted, right)
	case ast.Bool:
		return boolStrictEquals(leftAsserted, right)
	case ast.String:
		return stringStrictEquals(leftAsserted, right)
	case ast.Number:
		return numberStrictEquals(leftAsserted, right)
	default:
		return false, fmt.Errorf("cannot compare with %T", left)
	}
}

func boolStrictEquals(left ast.Bool, right any) (bool, error) {
	switch asserted := right.(type) {
	case ast.Bool:
		return left == asserted, nil
	case bool:
		return bool(left) == asserted, nil
	default:
		return false, ErrIncompatibleTypes
	}
}

func nullStrictEquals(left ast.Null, right any) (bool, error) {
	switch right.(type) {
	case ast.Null:
		return true, nil
	case nil:
		return true, nil
	default:
		return false, ErrIncompatibleTypes
	}
}

func stringStrictEquals(left ast.String, right any) (bool, error) {
	switch asserted := right.(type) {
	case ast.String:
		return left == asserted, nil
	case string:
		return string(left) == asserted, nil
	default:
		return false, ErrIncompatibleTypes
	}
}

func numberStrictEquals(left ast.Number, right any) (bool, error) {
	leftIntValue, ok := left.ToInteger()
	if ok {
		switch asserted := right.(type) {
		case ast.Number:
			rightIntValue, ok := asserted.ToInteger()
			if ok {
				return leftIntValue == rightIntValue, nil
			}

			return false, ErrIncompatibleTypes
		case int64:
			return leftIntValue == asserted, nil
		case float64:
			return false, ErrIncompatibleTypes
		default:
			return false, ErrIncompatibleTypes
		}
	}

	leftFloatValue := left.ToFloat()

	switch asserted := right.(type) {
	case ast.Number:
		if asserted.IsFloat() {
			return leftFloatValue == asserted.ToFloat(), nil
		}

		return false, ErrIncompatibleTypes
	case int64:
		return false, ErrIncompatibleTypes
	case float64:
		return leftFloatValue == asserted, nil
	default:
		return false, ErrIncompatibleTypes
	}
}
