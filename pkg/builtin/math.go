// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

func integerAddFunction(a, b int64, extra ...int64) (any, error) {
	sum := a + b
	for _, num := range extra {
		sum += num
	}

	return sum, nil
}

func numberAddFunction(a, b ast.Number, extra ...ast.Number) (any, error) {
	sum := a.MustToFloat() + b.MustToFloat()
	for _, num := range extra {
		sum += num.MustToFloat()
	}

	return sum, nil
}

func integerSubFunction(a, b int64, extra ...int64) (any, error) {
	diff := a - b
	for _, num := range extra {
		diff -= num
	}

	return diff, nil
}

func numberSubFunction(a, b ast.Number, extra ...ast.Number) (any, error) {
	diff := a.MustToFloat() - b.MustToFloat()
	for _, num := range extra {
		diff -= num.MustToFloat()
	}

	return diff, nil
}

func integerMultFunction(a, b int64, extra ...int64) (any, error) {
	product := a * b
	for _, num := range extra {
		product *= num
	}

	return product, nil
}

func numberMultFunction(a, b ast.Number, extra ...ast.Number) (any, error) {
	product := a.MustToFloat() * b.MustToFloat()
	for _, num := range extra {
		product *= num.MustToFloat()
	}

	return product, nil
}

func integerDivFunction(a, b int64, extra ...int64) (any, error) {
	if b == 0 {
		return nil, errors.New("division by zero")
	}

	result := a / b

	for _, num := range extra {
		if num == 0 {
			return nil, errors.New("division by zero")
		}

		result /= num
	}

	return result, nil
}

func numberDivFunction(a, b ast.Number, extra ...ast.Number) (any, error) {
	if b.MustToFloat() == 0 {
		return nil, errors.New("division by zero")
	}

	result := a.MustToFloat() / b.MustToFloat()

	for _, num := range extra {
		if num.MustToFloat() == 0 {
			return nil, errors.New("division by zero")
		}

		result /= num.MustToFloat()
	}

	return result, nil
}
