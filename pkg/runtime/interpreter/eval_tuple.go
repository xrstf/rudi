// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package interpreter

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/deepcopy"
	"go.xrstf.de/rudi/pkg/jsonpath"
	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/pathexpr"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

func (i *interpreter) EvalTuple(ctx types.Context, tup ast.Tuple) (any, error) {
	// Function calls are the only place where we check if the Go context has been cancelled.
	// This error should not be caught and swallowed by any other function, like `try` or `default`.
	if err := ctx.GoContext().Err(); err != nil {
		return nil, err
	}

	if len(tup.Expressions) == 0 {
		return nil, errors.New("invalid tuple: tuple cannot be empty")
	}

	identifier, ok := tup.Expressions[0].(ast.Identifier)
	if !ok {
		return nil, errors.New("invalid tuple: first expression must be an identifier")
	}

	result, err := i.CallFunction(ctx, identifier, tup.Expressions[1:])
	if err != nil {
		return nil, err
	}

	deeper, err := pathexpr.Apply(ctx, result, tup.PathExpression)
	if err != nil {
		return nil, err
	}

	return deeper, nil
}

func (*interpreter) CallFunction(ctx types.Context, fun ast.Identifier, args []ast.Expression) (any, error) {
	funcName := fun.Name
	function, ok := ctx.GetFunction(funcName)
	if !ok {
		return nil, fmt.Errorf("unknown function %s", funcName)
	}

	// call the function
	result, err := function.Evaluate(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	// if desired, update the context and introduce side effects
	if fun.Bang {
		// "delete!" has a special behaviour for the bang modifier, so do possibly some other functions.
		if custom, ok := function.(types.BangHandler); ok {
			return custom.BangHandler(ctx, args, result)
		}

		// prepare handling a possible bang (like `(append! .foo 12)`)
		var updateSymbol *ast.Symbol
		if fun.Bang {
			if len(args) == 0 {
				return nil, fmt.Errorf("%s must have at least 1 symbol argument", fun.String())
			}

			firstArg := args[0]
			symbol, ok := firstArg.(ast.Symbol)
			if !ok {
				return nil, fmt.Errorf("%s must use Symbol as first argument, got %T", fun.String(), firstArg)
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
			pathExpr, err := pathexpr.Eval(ctx, updateSymbol.PathExpression)
			if err != nil {
				return nil, fmt.Errorf("argument #0: invalid path expression: %w", err)
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
				return nil, err
			}

			// apply the path expression
			updatedValue, err = jsonpath.Set(currentValue, jsonpath.FromEvaluatedPath(*pathExpr), updatedValue)
			if err != nil {
				return nil, fmt.Errorf("cannot set value in %T at %s: %w", currentValue, pathExpr, err)
			}
		}

		if updateSymbol.Variable != nil {
			varName := string(*updateSymbol.Variable)
			ctx.SetVariable(varName, updatedValue)
		} else {
			ctx.GetDocument().Set(updatedValue)
		}
	}

	return result, nil
}
