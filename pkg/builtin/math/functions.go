// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package math

import (
	"errors"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/functions"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

var (
	addRudiFunction      = functions.NewBuilder(integerAddFunction, numberAddFunction).WithDescription("returns the sum of all of its arguments").Build()
	subRudiFunction      = functions.NewBuilder(integerSubFunction, numberSubFunction).WithDescription("returns arg1 - arg2 - .. - argN").Build()
	multiplyRudiFunction = functions.NewBuilder(integerMultFunction, numberMultFunction).WithDescription("returns the product of all of its arguments").Build()
	divideRudiFunction   = functions.NewBuilder(numberDivFunction).WithDescription("returns arg1 / arg2 / .. / argN (always a floating point division, regardless of arguments)").Build()

	Functions = types.Functions{
		// These are the main functions, but within the documentation these are
		// considered to be aliases because their names cannot be used in Markdown
		// filenames.
		"+": addRudiFunction,
		"-": subRudiFunction,
		"*": multiplyRudiFunction,
		"/": divideRudiFunction,

		// aliases to make bang functions nicer (add! vs +!)
		"add":  addRudiFunction,
		"sub":  subRudiFunction,
		"mult": multiplyRudiFunction,
		"div":  divideRudiFunction,
	}
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
