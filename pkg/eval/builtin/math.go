// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func evalNumericalExpressions(ctx types.Context, args []ast.Expression) (values []ast.Number, int64only bool, err error) {
	values = make([]ast.Number, len(args))
	int64only = true

	for i, arg := range args {
		_, evaluated, err := eval.EvalExpression(ctx, arg)
		if err != nil {
			return nil, false, fmt.Errorf("argument #%d: %w", i, err)
		}

		num, err := ctx.Coalesce().ToNumber(evaluated)
		if err != nil {
			return nil, false, fmt.Errorf("argument #%d: %w", i, err)
		}

		values[i] = num

		if _, isInt := num.ToInteger(); !isInt {
			int64only = false
		}
	}

	return values, int64only, nil
}

func sumFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size < 2 {
		return nil, fmt.Errorf("expected 2+ arguments, got %d", size)
	}

	values, int64sOnly, err := evalNumericalExpressions(ctx, args)
	if err != nil {
		return nil, err
	}

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
		sum += num.ToFloat()
	}

	return sum, nil
}

func subFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size < 2 {
		return nil, fmt.Errorf("expected 2+ arguments, got %d", size)
	}

	values, int64sOnly, err := evalNumericalExpressions(ctx, args)
	if err != nil {
		return nil, err
	}

	if int64sOnly {
		difference, _ := values[0].ToInteger()
		for _, num := range values[1:] {
			val, _ := num.ToInteger()
			difference -= val
		}

		return difference, nil
	}

	difference := values[0].ToFloat()
	for _, num := range values[1:] {
		difference -= num.ToFloat()
	}

	return difference, nil
}

func multiplyFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size < 2 {
		return nil, fmt.Errorf("expected 2+ arguments, got %d", size)
	}

	values, int64sOnly, err := evalNumericalExpressions(ctx, args)
	if err != nil {
		return nil, err
	}

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
		product *= num.ToFloat()
	}

	return product, nil
}

func divideFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size < 2 {
		return nil, fmt.Errorf("expected 2+ arguments, got %d", size)
	}

	values, _, err := evalNumericalExpressions(ctx, args)
	if err != nil {
		return nil, err
	}

	result := values[0].ToFloat()

	for _, num := range values[1:] {
		divisor := num.ToFloat()
		if divisor == 0 {
			return nil, errors.New("division by zero")
		}
		result /= divisor
	}

	return result, nil
}
