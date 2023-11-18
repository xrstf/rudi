// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package equality

import (
	"fmt"

	"go.xrstf.de/otto/pkg/eval/coalescing"
	"go.xrstf.de/otto/pkg/eval/types"
	"go.xrstf.de/otto/pkg/lang/ast"
)

// equality, but with using coalescing, so 1 == "1"
func EqualEnough(left, right ast.Literal) (bool, error) {
	// if either of the sides is a null, convert the other to null
	matched, equal, err := nullishEqualEnough(left, right)
	if matched {
		return equal, err
	}

	// if either of the sides is a bool, convert the other to bool
	matched, equal, err = boolishEqualEnough(left, right)
	if matched {
		return equal, err
	}

	// if either of the sides is a number, convert the other to a number
	matched, equal, err = numberishEqualEnough(left, right)
	if matched {
		return equal, err
	}

	// if either of the sides is a string, convert the other to a string
	matched, equal, err = stringishEqualEnough(left, right)
	if matched {
		return equal, err
	}

	// now both sides can basically just be vectors or objects

	switch leftAsserted := left.(type) {
	case ast.Vector:
		return vectorishEqualEnough(leftAsserted, right)
	case ast.Object:
		return objectishEqualEnough(leftAsserted, right)
	default:
		return false, fmt.Errorf("cannot compare with %T with %T", left, right)
	}
}

func nullishEqualEnough(left ast.Literal, right ast.Literal) (matched bool, equal bool, err error) {
	_, leftOk := left.(ast.Null)
	_, rightOk := right.(ast.Null)

	if !leftOk && !rightOk {
		return false, false, nil
	}

	matched = true

	if leftOk && rightOk {
		return matched, true, nil
	}

	var b ast.Literal

	if leftOk {
		b = right
	} else {
		b = left
	}

	bValue, err := coalescing.ToBool(b)
	if err != nil {
		return matched, false, ErrIncompatibleTypes
	}

	return matched, !bValue, nil
}

func boolishEqualEnough(left ast.Literal, right ast.Literal) (matched bool, equal bool, err error) {
	leftBool, leftOk := left.(ast.Bool)
	rightBool, rightOk := right.(ast.Bool)

	if !leftOk && !rightOk {
		return false, false, nil
	}

	matched = true

	if leftOk && rightOk {
		return matched, leftBool.Equal(rightBool), nil
	}

	var (
		a bool
		b ast.Literal
	)

	if leftOk {
		a = bool(leftBool)
		b = right
	} else {
		a = bool(rightBool)
		b = left
	}

	bValue, err := coalescing.ToBool(b)
	if err != nil {
		return matched, false, ErrIncompatibleTypes
	}

	return matched, a == bValue, nil
}

func numberishEqualEnough(left ast.Literal, right ast.Literal) (matched bool, equal bool, err error) {
	leftNumber, leftOk := left.(ast.Number)
	rightNumber, rightOk := right.(ast.Number)

	if !leftOk && !rightOk {
		return false, false, nil
	}

	matched = true

	if leftOk && rightOk {
		return matched, leftNumber.ToFloat() == rightNumber.ToFloat(), nil
	}

	var (
		a ast.Number
		b ast.Literal
	)

	if leftOk {
		a = leftNumber
		b = right
	} else {
		a = rightNumber
		b = left
	}

	bValue, err := coalescing.ToFloat64(b)
	if err != nil {
		return matched, false, ErrIncompatibleTypes
	}

	return matched, a.ToFloat() == bValue, nil
}

func stringishEqualEnough(left ast.Literal, right ast.Literal) (matched bool, equal bool, err error) {
	leftString, leftOk := left.(ast.String)
	rightString, rightOk := right.(ast.String)

	if !leftOk && !rightOk {
		return false, false, nil
	}

	matched = true

	if leftOk && rightOk {
		return matched, leftString.Equal(rightString), nil
	}

	var (
		a string
		b ast.Literal
	)

	if leftOk {
		a = string(leftString)
		b = right
	} else {
		a = string(rightString)
		b = left
	}

	bValue, err := coalescing.ToString(b)
	if err != nil {
		return matched, false, ErrIncompatibleTypes
	}

	return matched, a == bValue, nil
}

func vectorishEqualEnough(left ast.Vector, right any) (bool, error) {
	// extra: [] == {}
	rightObject, ok := right.(ast.Object)
	if ok {
		return len(left.Data) == 0 && len(rightObject.Data) == 0, nil
	}

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
		equal, err := EqualEnough(leftWrapped.(ast.Literal), rightWrapped.(ast.Literal))
		if err != nil {
			return false, err
		}

		if !equal {
			return false, nil
		}
	}

	return true, nil
}

func objectishEqualEnough(left ast.Object, right any) (bool, error) {
	// extra: [] == {}
	rightVector, ok := right.(ast.Vector)
	if ok {
		return len(left.Data) == 0 && len(rightVector.Data) == 0, nil
	}

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
		equal, err := EqualEnough(leftWrapped.(ast.Literal), rightWrapped.(ast.Literal))
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
