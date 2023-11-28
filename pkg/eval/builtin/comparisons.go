// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/equality"
	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func eqFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 2 {
		return nil, fmt.Errorf("expected exactly 2 arguments, got %d", size)
	}

	_, leftData, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	_, rightData, err := eval.EvalExpression(ctx, args[1])
	if err != nil {
		return nil, fmt.Errorf("argument #1: %w", err)
	}

	equal, err := equality.EqualCoalesced(ctx.Coalesce(), leftData, rightData)
	if err != nil {
		return nil, err
	}

	return equal, nil
}

func likeFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 2 {
		return nil, fmt.Errorf("expected exactly 2 arguments, got %d", size)
	}

	_, leftData, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	_, rightData, err := eval.EvalExpression(ctx, args[1])
	if err != nil {
		return nil, fmt.Errorf("argument #1: %w", err)
	}

	equal, err := equality.EqualCoalesced(coalescing.NewHumane(), leftData, rightData)
	if err != nil {
		return nil, err
	}

	return equal, nil
}

type intProcessor func(left, right int64) (bool, error)
type floatProcessor func(left, right float64) (bool, error)

func makeNumberComparatorFunc(inter intProcessor, floater floatProcessor, desc string) types.Function {
	return types.BasicFunction(func(ctx types.Context, args []ast.Expression) (any, error) {
		if size := len(args); size != 2 {
			return nil, fmt.Errorf("expected 2 argument(s), got %d", size)
		}

		numbers, _, err := evalNumericalExpressions(ctx, args)
		if err != nil {
			return nil, err
		}

		leftInt, leftOk := numbers[0].ToInteger()
		rightInt, rightOk := numbers[1].ToInteger()

		if leftOk != rightOk {
			return nil, errors.New("cannot compare floats to integers")
		}

		if leftOk && rightOk {
			return inter(leftInt, rightInt)
		}

		leftFloat := numbers[0].ToFloat()
		rightFloat := numbers[1].ToFloat()

		return floater(leftFloat, rightFloat)
	}, desc)
}
