// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/deepcopy"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/pathexpr"
)

func EvalTuple(ctx types.Context, tup ast.Tuple) (types.Context, any, error) {
	// Function calls are the only place where we check if the Go context has been cancelled.
	// This error should not be caught and swallowed by any other function, like `try` or `default`.
	if err := ctx.Context().Err(); err != nil {
		return ctx, nil, err
	}

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

type BangHandler interface {
	// All functions work fine with the default bang handler ("set!", "append!", ...), except
	// for "delete!", which requires special handling to make it work as expected. Custom bang
	// handlers are useful to introducing side effects explicitly (so it becomes very clear if a
	// function in Rudi has side effects or not).
	BangHandler(ctx types.Context, args []ast.Expression, value any) (types.Context, any, error)
}

func EvalFunctionCall(ctx types.Context, fun ast.Identifier, args []ast.Expression) (types.Context, any, error) {
	funcName := fun.Name
	function, ok := ctx.GetFunction(funcName)
	if !ok {
		return ctx, nil, fmt.Errorf("unknown function %s", funcName)
	}

	// call the function
	result, err := function.Evaluate(ctx, args)
	if err != nil {
		return ctx, nil, fmt.Errorf("%s: %w", funcName, err)
	}

	resultCtx := ctx

	// if desired, update the context and introduce side effects
	if fun.Bang {
		// "delete!" has a special behaviour for the bang modifier, so do possibly some other functions.
		if custom, ok := function.(BangHandler); ok {
			return custom.BangHandler(ctx, args, result)
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

		// We always return the computed value, no matter how deep we inject it into the symbol;
		// but for setting the new variable/document, we need the _whole_ new value, which might be
		// the result of combining the current + setting a value deep somewhere.
		updatedValue := result

		// if the symbol has a path to traverse, do so
		if updateSymbol.PathExpression != nil {
			// pre-evaluate the path expression
			pathExpr, err := EvalPathExpression(ctx, updateSymbol.PathExpression)
			if err != nil {
				return ctx, nil, fmt.Errorf("argument #0: invalid path expression: %w", err)
			}

			// get the current value of the symbol
			var currentValue any

			if updateSymbol.Variable != nil {
				varName := string(*updateSymbol.Variable)

				// a non-existing variable is fine, this is how you define new variables in the first place
				currentValue, _ = ctx.GetVariable(varName)
			} else {
				currentValue = ctx.GetDocument().Data()
			}

			currentValue, err = deepcopy.Clone(currentValue)
			if err != nil {
				return ctx, nil, err
			}

			// apply the path expression
			updatedValue, err = pathexpr.Set(currentValue, pathexpr.FromEvaluatedPath(*pathExpr), updatedValue)
			if err != nil {
				return ctx, nil, fmt.Errorf("cannot set value in %T at %s: %w", currentValue, pathExpr, err)
			}
		}

		if updateSymbol.Variable != nil {
			varName := string(*updateSymbol.Variable)
			resultCtx = resultCtx.WithVariable(varName, updatedValue)
		} else {
			ctx.GetDocument().Set(updatedValue)
		}
	}

	return resultCtx, result, nil
}
