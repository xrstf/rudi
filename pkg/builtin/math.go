// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

func integerAddFunction(base int64, extra ...int64) (any, error) {
	for _, num := range extra {
		base += num
	}

	return base, nil
}

func numberAddFunction(base ast.Number, extra ...ast.Number) (any, error) {
	sum := base.MustToFloat()
	for _, num := range extra {
		sum += num.MustToFloat()
	}

	return sum, nil
}

func integerSubFunction(base int64, extra ...int64) (any, error) {
	for _, num := range extra {
		base -= num
	}

	return base, nil
}

func numberSubFunction(base ast.Number, extra ...ast.Number) (any, error) {
	diff := base.MustToFloat()
	for _, num := range extra {
		diff -= num.MustToFloat()
	}

	return diff, nil
}

func integerMultFunction(base int64, extra ...int64) (any, error) {
	for _, num := range extra {
		base *= num
	}

	return base, nil
}

func numberMultFunction(base ast.Number, extra ...ast.Number) (any, error) {
	product := base.MustToFloat()
	for _, num := range extra {
		product *= num.MustToFloat()
	}

	return product, nil
}

func integerDivFunction(base int64, extra ...int64) (any, error) {
	for _, num := range extra {
		if num == 0 {
			return nil, errors.New("division by zero")
		}

		base /= num
	}

	return base, nil
}

func numberDivFunction(base ast.Number, extra ...ast.Number) (any, error) {
	result := base.MustToFloat()

	for _, num := range extra {
		if num.MustToFloat() == 0 {
			return nil, errors.New("division by zero")
		}

		result /= num.MustToFloat()
	}

	return result, nil
}
