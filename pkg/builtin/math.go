// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/eval/util"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

type numberifiedLiteralFunction func(ctx types.Context, values []ast.Number, int64sOnly bool) (any, error)

func numberifyArgs(fun numberifiedLiteralFunction) util.LiteralFunction {
	return func(ctx types.Context, args []any) (any, error) {
		values := make([]ast.Number, len(args))
		int64sOnly := true

		for i, arg := range args {
			num, err := ctx.Coalesce().ToNumber(arg)
			if err != nil {
				return nil, fmt.Errorf("argument #%d: %w", i, err)
			}

			values[i] = num

			if _, isInt := num.ToInteger(); !isInt {
				int64sOnly = false
			}
		}

		return fun(ctx, values, int64sOnly)
	}
}

func addFunction(ctx types.Context, values []ast.Number, int64sOnly bool) (any, error) {
	if int64sOnly {
		sum := int64(0)
		for _, num := range values {
			val, _ := num.ToInteger()
			sum += val
		}

		return sum, nil
	}

	sum := float64(0)
	for _, num := range values {
		sum += num.MustToFloat()
	}

	return sum, nil
}

func subFunction(ctx types.Context, values []ast.Number, int64sOnly bool) (any, error) {
	if int64sOnly {
		difference, _ := values[0].ToInteger()
		for _, num := range values[1:] {
			val, _ := num.ToInteger()
			difference -= val
		}

		return difference, nil
	}

	difference := values[0].MustToFloat()
	for _, num := range values[1:] {
		difference -= num.MustToFloat()
	}

	return difference, nil
}

func multiplyFunction(ctx types.Context, values []ast.Number, int64sOnly bool) (any, error) {
	if int64sOnly {
		product := int64(1)
		for _, num := range values {
			factor, _ := num.ToInteger()
			product *= factor
		}

		return product, nil
	}

	product := float64(1)
	for _, num := range values {
		product *= num.MustToFloat()
	}

	return product, nil
}

func divideFunction(ctx types.Context, values []ast.Number, int64sOnly bool) (any, error) {
	result := values[0].MustToFloat()

	for _, num := range values[1:] {
		divisor := num.MustToFloat()
		if divisor == 0 {
			return nil, errors.New("division by zero")
		}

		result /= divisor
	}

	return result, nil
}
