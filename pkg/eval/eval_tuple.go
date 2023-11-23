// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func EvalTuple(ctx types.Context, tup ast.Tuple) (types.Context, any, error) {
	if len(tup.Expressions) == 0 {
		return ctx, nil, errors.New("invalid tuple: tuple cannot be empty")
	}

	identifier, ok := tup.Expressions[0].(ast.Identifier)
	if !ok {
		return ctx, nil, errors.New("invalid tuple: first expression must be an identifier")
	}

	resultCtx, result, err := EvalFunctionCall(ctx, identifier, tup.Expressions[1:])
	if err != nil {
		return ctx, nil, err
	}

	if tup.PathExpression != nil {
		deeper, err := TraversePathExpression(ctx, result, tup.PathExpression)
		if err != nil {
			return ctx, nil, err
		}

		return resultCtx, deeper, nil
	}

	return resultCtx, result, nil
}

func EvalFunctionCall(ctx types.Context, fun ast.Identifier, args []ast.Expression) (types.Context, any, error) {
	funcName := fun.Name
	function, ok := ctx.GetFunction(funcName)
	if !ok {
		return ctx, nil, fmt.Errorf("unknown function %s", funcName)
	}

	// prepare handling a possible bang (like `(append! .foo 12)`)
	var updateSymbol *ast.Symbol
	if fun.Bang {
		if len(args) == 0 {
			return ctx, nil, fmt.Errorf("%s must have at least 1 symbol argument", fun.String())
		}

		firstArg := args[0]
		symbol, ok := firstArg.(ast.Symbol)
		if !ok {
			return ctx, nil, fmt.Errorf("%s must use Symbol as first argument, got %T", fun.String(), firstArg)
		}

		updateSymbol = &symbol
	}

	// call the function
	result, err := function(ctx, args)
	if err != nil {
		return ctx, nil, fmt.Errorf("%s: %w", funcName, err)
	}

	resultCtx := ctx

	// if desired, update the symbol's value
	if updateSymbol != nil {
		if updateSymbol.Variable != nil {
			varName := string(*updateSymbol.Variable)
			resultCtx = resultCtx.WithVariable(varName, result)
		} else {
			ctx.GetDocument().Set(result)
		}
	}

	return resultCtx, result, nil
}
