// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package equality

import (
	"errors"
	"fmt"
	"reflect"

	"go.xrstf.de/rudi/pkg/coalescing"
)

var ErrIncompatibleTypes = errors.New("types are incompatible")

type Comparer interface {
	Compare(other any) (int, error)
}

// Make return statements easier to read.
const (
	doNotCare   = 0
	Unorderable = 2
	IsEqual     = 0
	IsSmaller   = -1
	IsGreater   = 1
)

// Re-use a pedantic coalescer to not repeat the logic to turn int/int32/int64 into int64.
var typeChecker = coalescing.NewPedantic()

func Equal(c coalescing.Coalescer, left, right any) (bool, error) {
	compared, err := Compare(c, left, right)
	if err != nil {
		return false, err
	}

	return compared == IsEqual, nil
}

func Compare(c coalescing.Coalescer, left, right any) (int, error) {
	if c == nil {
		c = coalescing.NewStrict()
	}

	// if either of the sides is a null, convert the other to null
	matched, compared, err := compareNullish(c, left, right)
	if err != nil {
		return doNotCare, err
	}
	if matched {
		return compared, nil
	}

	// if either of the sides is a bool, convert the other to bool
	matched, compared, err = compareBoolish(c, left, right)
	if err != nil {
		return doNotCare, err
	}
	if matched {
		return compared, nil
	}

	// if either of the sides is a float, convert the other to a float
	matched, compared, err = compareFloatish(c, left, right)
	if err != nil {
		return doNotCare, err
	}
	if matched {
		return compared, nil
	}

	// if either of the sides is a int, convert the other to a int
	matched, compared, err = compareIntish(c, left, right)
	if err != nil {
		return doNotCare, err
	}
	if matched {
		return compared, nil
	}

	// if either of the sides is a string, convert the other to a string
	matched, compared, err = compareStringish(c, left, right)
	if err != nil {
		return doNotCare, err
	}
	if matched {
		return compared, nil
	}

	// if either of the sides is a vector, convert the other to a vector
	matched, compared, err = compareVectorish(c, left, right)
	if err != nil {
		return doNotCare, err
	}
	if matched {
		return compared, nil
	}

	// if either of the sides is an object, convert the other to an object
	matched, compared, err = compareObjectish(c, left, right)
	if err != nil {
		return doNotCare, err
	}
	if matched {
		return compared, nil
	}

	// Allow to compare 2 values of the same type, if that type implements an equal function.
	if reflect.TypeOf(left) == reflect.TypeOf(right) {
		if comparer, ok := left.(Comparer); ok {
			return comparer.Compare(right)
		}
	}

	return doNotCare, fmt.Errorf("cannot compare with %T with %T", left, right)
}

func compareNullish(c coalescing.Coalescer, left any, right any) (matched bool, compared int, err error) {
	leftOk := left == nil
	rightOk := right == nil

	if !leftOk && !rightOk {
		compared = doNotCare
		return
	}

	matched = true

	if leftOk && rightOk {
		compared = IsEqual
		return
	}

	var other any

	if leftOk {
		other = right
	} else {
		other = left
	}

	isNullish, err := c.ToNull(other)
	if err != nil {
		compared = doNotCare
		return
	}

	// both values are equal
	if isNullish {
		compared = IsEqual
		return
	}

	// the other value was not nullish, e.g. when doing "null == true"
	if leftOk {
		// null <-> something
		compared = IsSmaller
		return
	}

	// something <-> null
	compared = IsGreater
	return
}

func compareBool(left, right bool) int {
	switch {
	case left == right:
		return IsEqual
	case left && !right:
		return IsGreater
	default:
		return IsSmaller
	}
}

func compareBoolish(c coalescing.Coalescer, left any, right any) (matched bool, compared int, err error) {
	leftBool, leftOk := left.(bool)
	rightBool, rightOk := right.(bool)

	if !leftOk && !rightOk {
		compared = doNotCare
		return
	}

	matched = true

	if leftOk && rightOk {
		compared = compareBool(leftBool, rightBool)
		return
	}

	// convert the other value to a bool
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
		compared = doNotCare
		return
	}

	if bValue == a {
		compared = IsEqual
		return
	}

	if leftOk {
		rightBool = bValue
	} else {
		leftBool = bValue
	}

	compared = compareBool(leftBool, rightBool)
	return
}

func compareInt(left, right int64) int {
	switch {
	case left == right:
		return IsEqual
	case left < right:
		return IsSmaller
	default:
		return IsGreater
	}
}

func compareIntish(c coalescing.Coalescer, left any, right any) (matched bool, compared int, err error) {
	leftInt, leftErr := typeChecker.ToInt64(left)
	rightInt, rightErr := typeChecker.ToInt64(right)

	if leftErr != nil && rightErr != nil {
		compared = doNotCare
		return
	}

	matched = true

	if leftErr == nil && rightErr == nil {
		compared = compareInt(leftInt, rightInt)
		return
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
		compared = doNotCare
		return
	}

	if bValue == a {
		compared = IsEqual
		return
	}

	if leftErr == nil {
		rightInt = bValue
	} else {
		leftInt = bValue
	}

	compared = compareInt(leftInt, rightInt)
	return
}

func compareFloat(left, right float64) int {
	switch {
	case left == right:
		return IsEqual
	case left < right:
		return IsSmaller
	default:
		return IsGreater
	}
}

func compareFloatish(c coalescing.Coalescer, left any, right any) (matched bool, compared int, err error) {
	leftFloat, leftErr := typeChecker.ToFloat64(left)
	rightFloat, rightErr := typeChecker.ToFloat64(right)

	if leftErr != nil && rightErr != nil {
		compared = doNotCare
		return
	}

	matched = true

	if leftErr == nil && rightErr == nil {
		compared = compareFloat(leftFloat, rightFloat)
		return
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
		compared = doNotCare
		return
	}

	if bValue == a {
		compared = IsEqual
		return
	}

	if leftErr == nil {
		rightFloat = bValue
	} else {
		leftFloat = bValue
	}

	compared = compareFloat(leftFloat, rightFloat)
	return
}

func compareString(left, right string) int {
	switch {
	case left == right:
		return IsEqual
	case left < right:
		return IsSmaller
	default:
		return IsGreater
	}
}

func compareStringish(c coalescing.Coalescer, left any, right any) (matched bool, compared int, err error) {
	leftString, leftOk := left.(string)
	rightString, rightOk := right.(string)

	if !leftOk && !rightOk {
		compared = doNotCare
		return
	}

	matched = true

	if leftOk && rightOk {
		compared = compareString(leftString, rightString)
		return
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
		compared = doNotCare
		return
	}

	if bValue == a {
		compared = IsEqual
		return
	}

	if leftOk {
		rightString = bValue
	} else {
		leftString = bValue
	}

	compared = compareString(leftString, rightString)
	return
}

func compareVectorish(c coalescing.Coalescer, left any, right any) (matched bool, compared int, err error) {
	leftVector, leftErr := typeChecker.ToVector(left)
	rightVector, rightErr := typeChecker.ToVector(right)

	if leftErr != nil && rightErr != nil {
		compared = doNotCare
		return
	}

	var equal bool
	matched = true

	if leftErr == nil && rightErr == nil {
		equal, err = equalVectors(c, leftVector, rightVector)
		if err != nil {
			compared = doNotCare
			return
		}

		if equal {
			compared = IsEqual
			return
		}

		// vectors cannot be ordered (vec1 < vec2)
		compared = Unorderable
		return
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

	// vector conversion is only allowed if the A vector is empty, so that [] == {} depending on the coalescer
	if len(a) > 0 {
		compared = doNotCare
		err = ErrIncompatibleTypes
		return
	}

	bVector, err := c.ToVector(b)
	if err != nil {
		compared = doNotCare
		return
	}

	equal, err = equalVectors(c, a, bVector)
	if equal {
		compared = IsEqual
		return
	}

	// vectors cannot be ordered (vec1 < vec2)
	compared = Unorderable
	return
}

func equalVectors(c coalescing.Coalescer, left, right []any) (bool, error) {
	if len(left) != len(right) {
		return false, nil
	}

	for i, leftItem := range left {
		rightItem := right[i]

		// wrapping always returns literals, so type assertions are safe here
		equal, err := Equal(c, leftItem, rightItem)
		if err != nil {
			return false, err
		}
		if !equal {
			return false, nil
		}
	}

	return true, nil
}

func compareObjectish(c coalescing.Coalescer, left any, right any) (matched bool, compared int, err error) {
	leftObject, leftErr := typeChecker.ToObject(left)
	rightObject, rightErr := typeChecker.ToObject(right)

	if leftErr != nil && rightErr != nil {
		compared = doNotCare
		return
	}

	var equal bool
	matched = true

	if leftErr == nil && rightErr == nil {
		equal, err = equalObjects(c, leftObject, rightObject)
		if err != nil {
			compared = doNotCare
			return
		}

		if equal {
			compared = IsEqual
			return
		}

		// objects cannot be ordered (obj1 < obj2)
		compared = Unorderable
		return
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

	// object conversion is only allowed if the A object is empty, so that [] == {} depending on the coalescer
	if len(a) > 0 {
		compared = doNotCare
		err = ErrIncompatibleTypes
		return
	}

	bObject, err := c.ToObject(b)
	if err != nil {
		compared = doNotCare
		return
	}

	equal, err = equalObjects(c, a, bObject)
	if equal {
		compared = IsEqual
		return
	}

	// objects cannot be ordered (obj1 < obj2)
	compared = Unorderable
	return
}

func equalObjects(c coalescing.Coalescer, left, right map[string]any) (bool, error) {
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
		equal, err := Equal(c, leftItem, rightItem)
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
