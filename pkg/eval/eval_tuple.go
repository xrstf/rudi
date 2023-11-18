// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"errors"
	"fmt"

	"go.xrstf.de/otto/pkg/eval/types"
	"go.xrstf.de/otto/pkg/lang/ast"
)

func EvalTuple(ctx types.Context, tup ast.Tuple) (types.Context, any, error) {
	if len(tup.Expressions) == 0 {
		return ctx, nil, errors.New("invalid tuple: tuple cannot be empty")
	}

	funcExpr, ok := tup.Expressions[0].(ast.Identifier)
	if !ok {
		return ctx, nil, errors.New("invalid tuple: first expression must be an identifier")
	}

	funcName := string(funcExpr)

	function, ok := ctx.GetFunction(funcName)
	if !ok {
		return ctx, nil, fmt.Errorf("unknown function %s", funcName)
	}

	// call the function
	newContext, result, err := function(ctx, tup.Expressions[1:])
	if err != nil {
		return ctx, nil, fmt.Errorf("%s: %w", funcName, err)
	}

	if tup.PathExpression != nil {
		evaluated, err := EvalPathExpression(ctx, tup.PathExpression)
		if err != nil {
			return ctx, nil, fmt.Errorf("invalid path expression: %w", err)
		}

		deeper, err := traverseEvaluatedPathExpression(ctx, result, *evaluated)
		if err != nil {
			return ctx, nil, fmt.Errorf("cannot apply path %s: %w", evaluated.String(), err)
		}

		result, err := types.WrapNative(deeper)
		if err != nil {
			return ctx, nil, err
		}

		return ctx, result, nil
	}

	return newContext, result, nil
}
