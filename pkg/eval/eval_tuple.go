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

	args := tup.Expressions[1:]

	funcName := identifier.Name
	function, ok := ctx.GetFunction(funcName)
	if !ok {
		return ctx, nil, fmt.Errorf("unknown function %s", funcName)
	}

	// prepare handling a possible bang (like `(append! .foo 12)`)
	var updateSymbol *ast.Symbol
	if identifier.Bang {
		if len(args) == 0 {
			return ctx, nil, fmt.Errorf("%s must have at least 1 symbol argument", identifier.String())
		}

		firstArg := args[0]
		symbol, ok := firstArg.(ast.Symbol)
		if !ok {
			return ctx, nil, fmt.Errorf("%s must use Symbol as first argument, got %T", identifier.String(), firstArg)
		}

		updateSymbol = &symbol
	}

	// call the function
	newContext, result, err := function(ctx, tup.Expressions[1:])
	if err != nil {
		return ctx, nil, fmt.Errorf("%s: %w", funcName, err)
	}

	// if desired, update the symbol's value
	if updateSymbol != nil {
		if updateSymbol.Variable != nil {
			varName := string(*updateSymbol.Variable)
			newContext = ctx.WithVariable(varName, result)
		} else {
			ctx.GetDocument().Set(result)
		}
	}

	if tup.PathExpression != nil {
		deeper, err := TraversePathExpression(ctx, result, tup.PathExpression)
		if err != nil {
			return ctx, nil, err
		}

		return ctx, deeper, nil
	}

	return newContext, result, nil
}
