// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package equality

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

var ErrIncompatibleTypes = errors.New("types are incompatible")

func deliteral(val any) any {
	lit, ok := val.(ast.Literal)
	if ok {
		return lit.LiteralValue()
	}

	return val
}

// Re-use a pedantic coalescer to not repeat the logic to turn int/int32/int64 into int64.
var typeChecker = coalescing.NewPedantic()

func EqualCoalesced(c coalescing.Coalescer, left, right any) (bool, error) {
	if c == nil {
		c = coalescing.NewStrict()
	}

	left = deliteral(left)
	right = deliteral(right)

	// if either of the sides is a null, convert the other to null
	matched, equal, err := nullishEqualCoalesced(c, left, right)
	if err != nil {
		return false, err
	}
	if matched {
		return equal, nil
	}

	// if either of the sides is a bool, convert the other to bool
	matched, equal, err = boolishEqualCoalesced(c, left, right)
	if err != nil {
		return false, err
	}
	if matched {
		return equal, nil
	}

	// if either of the sides is a int, convert the other to a int
	matched, equal, err = intEqualCoalesced(c, left, right)
	if err != nil {
		return false, err
	}
	if matched {
		return equal, nil
	}

	// if either of the sides is a float, convert the other to a float
	matched, equal, err = floatEqualCoalesced(c, left, right)
	if err != nil {
		return false, err
	}
	if matched {
		return equal, nil
	}

	// if either of the sides is a string, convert the other to a string
	matched, equal, err = stringishEqualCoalesced(c, left, right)
	if err != nil {
		return false, err
	}
	if matched {
		return equal, nil
	}

	// if either of the sides is a vector, convert the other to a vector
	matched, equal, err = vectorishEqualCoalesced(c, left, right)
	if err != nil {
		return false, err
	}
	if matched {
		return equal, nil
	}

	// now only objects are left
	matched, equal, err = objectishEqualCoalesced(c, left, right)
	if err != nil {
		return false, err
	}
	if matched {
		return equal, nil
	}

	return false, fmt.Errorf("cannot compare with %T with %T", left, right)
}

func nullishEqualCoalesced(c coalescing.Coalescer, left any, right any) (matched bool, equal bool, err error) {
	leftOk := left == nil
	rightOk := right == nil

	if !leftOk && !rightOk {
		return false, false, nil
	}

	matched = true

	if leftOk && rightOk {
		return matched, true, nil
	}

	var other any

	if leftOk {
		other = right
	} else {
		other = left
	}

	isNullish, err := c.ToNull(other)
	if err != nil {
		return matched, false, err
	}

	return matched, isNullish, nil
}

func boolishEqualCoalesced(c coalescing.Coalescer, left any, right any) (matched bool, equal bool, err error) {
	leftBool, leftOk := left.(bool)
	rightBool, rightOk := right.(bool)

	if !leftOk && !rightOk {
		return false, false, nil
	}

	matched = true

	if leftOk && rightOk {
		return matched, leftBool == rightBool, nil
	}

	var (
		a bool
		b any
	)

	if leftOk {
		a = leftBool
		b = right
	} else {
		a = rightBool
		b = left
	}

	bValue, err := c.ToBool(b)
	if err != nil {
		return matched, false, err
	}

	return matched, a == bValue, nil
}

func intEqualCoalesced(c coalescing.Coalescer, left any, right any) (matched bool, equal bool, err error) {
	leftInt, leftErr := typeChecker.ToInt64(left)
	rightInt, rightErr := typeChecker.ToInt64(right)

	if leftErr != nil && rightErr != nil {
		return false, false, nil
	}

	matched = true

	if leftErr == nil && rightErr == nil {
		return matched, leftInt == rightInt, nil
	}

	var (
		a int64
		b any
	)

	if leftErr == nil {
		a = leftInt
		b = right
	} else {
		a = rightInt
		b = left
	}

	bValue, err := c.ToInt64(b)
	if err != nil {
		return matched, false, err
	}

	return matched, a == bValue, nil
}

func floatEqualCoalesced(c coalescing.Coalescer, left any, right any) (matched bool, equal bool, err error) {
	leftFloat, leftErr := typeChecker.ToFloat64(left)
	rightFloat, rightErr := typeChecker.ToFloat64(right)

	if leftErr != nil && rightErr != nil {
		return false, false, nil
	}

	matched = true

	if leftErr == nil && rightErr == nil {
		return matched, leftFloat == rightFloat, nil
	}

	var (
		a float64
		b any
	)

	if leftErr == nil {
		a = leftFloat
		b = right
	} else {
		a = rightFloat
		b = left
	}

	bValue, err := c.ToFloat64(b)
	if err != nil {
		return matched, false, err
	}

	return matched, a == bValue, nil
}

func stringishEqualCoalesced(c coalescing.Coalescer, left any, right any) (matched bool, equal bool, err error) {
	leftString, leftOk := left.(string)
	rightString, rightOk := right.(string)

	if !leftOk && !rightOk {
		return false, false, nil
	}

	matched = true

	if leftOk && rightOk {
		return matched, leftString == rightString, nil
	}

	var (
		a string
		b any
	)

	if leftOk {
		a = leftString
		b = right
	} else {
		a = rightString
		b = left
	}

	bValue, err := c.ToString(b)
	if err != nil {
		return matched, false, err
	}

	return matched, a == bValue, nil
}

func vectorishEqualCoalesced(c coalescing.Coalescer, left any, right any) (matched bool, equal bool, err error) {
	leftVector, leftErr := typeChecker.ToVector(left)
	rightVector, rightErr := typeChecker.ToVector(right)

	if leftErr != nil && rightErr != nil {
		return false, false, nil
	}

	matched = true

	if leftErr == nil && rightErr == nil {
		equal, err := vectorEqualCoalesced(c, leftVector, rightVector)

		return matched, equal, err
	}

	var (
		a []any
		b any
	)

	if leftErr == nil {
		a = leftVector
		b = right
	} else {
		a = rightVector
		b = left
	}

	// vector conversion is only allowed if the a vector is empty, so that [] == {} depending on the coalescer
	if len(a) > 0 {
		return matched, false, ErrIncompatibleTypes
	}

	bVector, err := c.ToVector(b)
	if err != nil {
		return matched, false, ErrIncompatibleTypes
	}

	equal, err = vectorEqualCoalesced(c, a, bVector)

	return matched, equal, err
}

func vectorEqualCoalesced(c coalescing.Coalescer, left, right []any) (bool, error) {
	if len(left) != len(right) {
		return false, nil
	}

	for i, leftItem := range left {
		rightItem := right[i]

		// wrapping always returns literals, so type assertions are safe here
		equal, err := EqualCoalesced(c, leftItem, rightItem)
		if err != nil {
			return false, err
		}
		if !equal {
			return false, nil
		}
	}

	return true, nil
}

func objectishEqualCoalesced(c coalescing.Coalescer, left any, right any) (matched bool, equal bool, err error) {
	leftObject, leftErr := typeChecker.ToObject(left)
	rightObject, rightErr := typeChecker.ToObject(right)

	if leftErr != nil && rightErr != nil {
		return false, false, nil
	}

	matched = true

	if leftErr == nil && rightErr == nil {
		equal, err := objectEqualCoalesced(c, leftObject, rightObject)

		return matched, equal, err
	}

	var (
		a map[string]any
		b any
	)

	if leftErr == nil {
		a = leftObject
		b = right
	} else {
		a = rightObject
		b = left
	}

	// vector conversion is only allowed if the a vector is empty, so that [] == {} depending on the coalescer
	if len(a) > 0 {
		return matched, false, ErrIncompatibleTypes
	}

	bObject, err := c.ToObject(b)
	if err != nil {
		return matched, false, ErrIncompatibleTypes
	}

	equal, err = objectEqualCoalesced(c, a, bObject)

	return matched, equal, err
}

func objectEqualCoalesced(c coalescing.Coalescer, left, right map[string]any) (bool, error) {
	if len(left) != len(right) {
		return false, nil
	}

	keysSeen := map[string]struct{}{}

	for key, leftItem := range left {
		rightItem, exists := right[key]
		if !exists {
			return false, nil
		}

		keysSeen[key] = struct{}{}

		// wrapping always returns literals, so type assertions are safe here
		equal, err := EqualCoalesced(c, leftItem, rightItem)
		if err != nil {
			return false, err
		}
		if !equal {
			return false, nil
		}
	}

	for key := range right {
		delete(keysSeen, key)
	}

	return len(keysSeen) == 0, nil
}
