// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package interpreter

import (
	"errors"
	"fmt"

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

	deeper, err := pathexpr.Traverse(ctx, result, tup.PathExpression)
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

	var (
		result any
		err    error
	)

	// if desired, update the context and introduce side effects
	if fun.Bang {
		result, err = callBangFunction(ctx, function, args)
	} else {
		result, err = function.Evaluate(ctx, args)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	return result, nil
}

func callBangFunction(ctx types.Context, fun types.Function, args []ast.Expression) (any, error) {
	// Functions can define their own bang handler for more specialised usecases.
	if custom, ok := fun.(types.BangHandler); ok {
		return custom.BangHandler(ctx, args)
	}

	// Regular bang handlers require at least one argument, which must be a symbol so that the
	// function result can actually be written somewhere, enabling the desired in-place behaviour.
	if len(args) == 0 {
		return nil, errors.New("must have at least 1 symbol argument")
	}

	firstArg := args[0]
	symbol, ok := firstArg.(ast.Symbol)
	if !ok {
		return nil, fmt.Errorf("first argument must be symbol, got %T", firstArg)
	}

	// prepare the path expression
	var pathExpr jsonpath.Path
	if symbol.PathExpression != nil {
		var err error

		pathExpr, err = pathexpr.ToJSONPath(ctx, symbol.PathExpression)
		if err != nil {
			return nil, fmt.Errorf("invalid path expression: %w", err)
		}
	}

	// Bang behaviour depends on the type of path expression on the symbol. If a filtered
	// expression is given, the function is called *once per match*. If the path is a simple path,
	// the function is evaluated only once.
	if pathExpr.HasFilterSteps() {
		return callBangFunctionWithFilter(ctx, fun, symbol, pathExpr, args)
	} else {
		return callBasicBangFunction(ctx, fun, symbol, pathExpr, args)
	}
}

func callBasicBangFunction(ctx types.Context, fun types.Function, dest ast.Symbol, path jsonpath.Path, args []ast.Expression) (any, error) {
	// get the current value of the symbol
	var currentValue any

	if dest.Variable != nil {
		varName := string(*dest.Variable)

		// a non-existing variable is fine, this is how new variables are defined in the first place
		currentValue, _ = ctx.GetVariable(varName)
	} else {
		currentValue = ctx.GetDocument().Data()
	}

	var result any

	patchedValue, err := jsonpath.Patch(currentValue, path, func(_ bool, _ any, _ any) (any, error) {
		var err error
		result, err = fun.Evaluate(ctx, args)
		return result, err
	})
	if err != nil {
		return nil, err
	}

	// We always return the computed value (result), no matter how deep we inject it into the symbol;
	// but for setting the new variable/document, we need the _whole_ new value (patchedValue), which
	// might be the result of combining the current + setting a value deep somewhere.

	if dest.Variable != nil {
		varName := string(*dest.Variable)
		ctx.SetVariable(varName, patchedValue)
	} else {
		ctx.GetDocument().Set(patchedValue)
	}

	return result, nil
}

func callBangFunctionWithFilter(ctx types.Context, fun types.Function, dest ast.Symbol, path jsonpath.Path, args []ast.Expression) (any, error) {
	// get the current value of the symbol
	var currentValue any

	if dest.Variable != nil {
		varName := string(*dest.Variable)

		// a non-existing variable is fine, this is how new variables are defined in the first place
		currentValue, _ = ctx.GetVariable(varName)
	} else {
		currentValue = ctx.GetDocument().Data()
	}

	// Patch all values found by traversing the path expression. For each value we find, we call fun()
	// and use the value as the first argument. This means for [{a: 1, b: [1, 2]}, {a: 1, b: [3, 4]} {a: 2}], an
	// (append! .[?(?eq .a 1)].b 42) would result in two calls to fun:
	// fun([1, 2], 42) and fun([3, 4], 42) and lead to
	// [{a: 1, b: [1, 2, 42]}, {a: 1, b: [3, 4, 42]} {a: 2}]
	// This works fine for regular functions like append, but higher-order functions like map/filter
	// need a custom bang handler because only they understand their semantics.

	patchedValue, err := jsonpath.Patch(currentValue, path, func(_ bool, _ any, val any) (any, error) {
		shimmed := ast.Shim{Value: val}
		funArgs := append([]ast.Expression{shimmed}, args[1:]...)

		return fun.Evaluate(ctx, funArgs)
	})
	if err != nil {
		return nil, err
	}

	if dest.Variable != nil {
		varName := string(*dest.Variable)
		ctx.SetVariable(varName, patchedValue)
	} else {
		ctx.GetDocument().Set(patchedValue)
	}

	return patchedValue, nil
}
