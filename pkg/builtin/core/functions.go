// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package core

import (
	"context"
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/builtin/helper"
	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/deepcopy"
	"go.xrstf.de/rudi/pkg/jsonpath"
	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/functions"
	"go.xrstf.de/rudi/pkg/runtime/pathexpr"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

var (
	strictCoalescer   = coalescing.NewStrict()
	pedanticCoalescer = coalescing.NewPedantic()
	humaneCoalescer   = coalescing.NewHumane()

	Functions = types.Functions{
		"default": functions.NewBuilder(defaultFunction).WithDescription("returns the default value if the first argument is empty").Build(),
		"delete":  functions.NewBuilder(deleteFunction).WithBangHandler(deleteFunctionBangHandler).WithDescription("removes a key from an object or an item from a vector").Build(),
		"do":      functions.NewBuilder(DoFunction).WithDescription("eval a sequence of statements where only one expression is valid").Build(),
		"empty?":  functions.NewBuilder(isEmptyFunction).WithCoalescer(humaneCoalescer).WithDescription("returns true when the given value is empty-ish (0, false, null, \"\", ...)").Build(),
		"error":   functions.NewBuilder(errorFunction, fmtErrorFunction).WithDescription("returns an error").Build(),
		"has?":    functions.NewBuilder(hasFunction).WithDescription("returns true if the given symbol's path expression points to an existing value").Build(),
		"if":      functions.NewBuilder(ifElseFunction, ifFunction).WithDescription("evaluate one of two expressions based on a condition").Build(),
		"case":    functions.NewBuilder(caseFunction).WithDescription("chooses the first expression for which the test is true").Build(),
		"set":     functions.NewBuilder(setFunction).WithBangHandler(setFunctionBangHandler).WithDescription("set a value in a variable/document, most often used with ! modifier (set!)").Build(),
		"patch":   functions.NewBuilder(patchIdentifierFunction).WithBangHandler(patchFunctionBangHandler).WithDescription("applies an expression to all matched values").Build(),
		"try":     functions.NewBuilder(tryWithFallbackFunction, tryFunction).WithDescription("returns the fallback if the first expression errors out").Build(),
	}
)

func keepContextCanceled(err error) error {
	if errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func ifFunction(ctx types.Context, test bool, yes ast.Expression) (any, error) {
	return ifElseFunction(ctx, test, yes, ast.Null{})
}

func ifElseFunction(ctx types.Context, test bool, yes, no ast.Expression) (any, error) {
	if test {
		return ctx.Runtime().EvalExpression(ctx, yes)
	}

	return ctx.Runtime().EvalExpression(ctx, no)
}

func caseFunction(ctx types.Context, exprs ...ast.Expression) (any, error) {
	if len(exprs)%2 != 0 {
		return nil, errors.New("expected an even number of arguments")
	}

	for i := 0; i < len(exprs); i += 2 {
		testExpr := exprs[i]
		valueExpr := exprs[i+1]

		result, err := ctx.Runtime().EvalExpression(ctx, testExpr)
		if err != nil {
			return nil, err
		}

		enabled, err := ctx.Coalesce().ToBool(result)
		if err != nil {
			return nil, err
		}

		if enabled {
			return ctx.Runtime().EvalExpression(ctx, valueExpr)
		}
	}

	// none of the case statements matched
	return nil, nil
}

// NB: Variadic functions always require at least 1 argument in Rudi to match.
// This function, doing basically "nothing", is re-used by other Rudi functions and therefore
// exported.
func DoFunction(ctx types.Context, args ...ast.Expression) (any, error) {
	var (
		result any
		err    error
	)

	for _, arg := range args {
		result, err = ctx.Runtime().EvalExpression(ctx, arg)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func hasFunction(ctx types.Context, arg ast.Expression) (any, error) {
	pathed, ok := arg.(ast.Pathed)
	if !ok {
		return nil, fmt.Errorf("expected datatype that can hold a path expression, got %T", arg)
	}

	// separate base value expression from the path expression
	pathExpr := pathed.GetPathExpression()
	expr := pathed.Pathless()

	if pathExpr == nil {
		return nil, errors.New("argument has no path expression")
	}

	// evaluate the base value
	value, err := ctx.Runtime().EvalExpression(ctx, expr)
	if err != nil {
		return nil, err
	}

	_, err = pathexpr.Traverse(ctx, value, pathExpr)
	if err != nil {
		return false, keepContextCanceled(err)
	}

	return true, nil
}

// (default TEST:Expression FALLBACK:any)
func defaultFunction(ctx types.Context, value any, fallback ast.Expression) (any, error) {
	// this function purposefully always uses humane coalescing, but only for this check
	boolified, err := coalescing.NewHumane().ToBool(value)
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	if boolified {
		return value, nil
	}

	value, err = ctx.Runtime().EvalExpression(ctx, fallback)
	if err != nil {
		return nil, fmt.Errorf("argument #1: %w", err)
	}

	return value, nil
}

func tryFunction(ctx types.Context, test ast.Expression) (any, error) {
	result, err := ctx.Runtime().EvalExpression(ctx, test)
	if err != nil {
		return nil, keepContextCanceled(err)
	}

	return result, nil
}

func tryWithFallbackFunction(ctx types.Context, test ast.Expression, fallback ast.Expression) (any, error) {
	result, err := ctx.Runtime().EvalExpression(ctx, test)
	if err != nil {
		result, err = ctx.Runtime().EvalExpression(ctx, fallback)
		if err != nil {
			return nil, fmt.Errorf("argument #1: %w", err)
		}
	}

	return result, nil
}

// (set VAR:Variable VALUE:any)
// (set EXPR:PathExpression VALUE:any)
func setFunction(ctx types.Context, target ast.Expression, value any) (any, error) {
	pathed, ok := target.(ast.Pathed)
	if !ok {
		return nil, fmt.Errorf("expected datatype that can hold a path expression, got %T", target)
	}

	// separate base value expression from the path expression
	pathExpr := pathed.GetPathExpression()
	expr := pathed.Pathless()

	// Make sure (set) calls make sense: For non-symbols, a path must be set, because
	// "(set (foo) 42)" or "(set [1 2 3] 42)" do not make sense. For symbols on the other
	// hand, we only pre-evaluate the target if a path exists, because we still need to
	// allow variables not existing, so that "(set! $var 42)" can succeed.
	switch target.(type) {
	case ast.Symbol:
		// Set relies entirely on the bang modifier handling to actually set values
		// in variables or the global document; without the bang modifier, (set)
		// is basically a no-op and we do not even have to evaluate the symbol here.
		if pathExpr == nil {
			return value, nil
		}

	case ast.Tuple, ast.VectorNode, ast.ObjectNode:
		if pathExpr == nil {
			return nil, errors.New("no path expression provided on target value")
		}

	default:
		return nil, fmt.Errorf("unexpected target value of type %T", target)
	}

	// evaluate the target
	targetValue, err := ctx.Runtime().EvalExpression(ctx, expr)
	if err != nil {
		return nil, err
	}

	// we need to operate on a _copy_ of the value, as updating happens in the bang handler later
	// on; this is only necessary for symbols though, as functions (tuples) are expected to return
	// non-pointer data and vector/objects nodes are literals.
	if _, ok := target.(ast.Symbol); ok {
		targetValue, err = deepcopy.Clone(targetValue)
		if err != nil {
			return nil, fmt.Errorf("invalid current value: %w", err)
		}
	}

	// prepare path expression
	jsonpathExpr, err := pathexpr.ToJSONPath(ctx, pathExpr)
	if err != nil {
		return nil, fmt.Errorf("invalid path expression: %w", err)
	}

	return jsonpath.Set(targetValue, jsonpathExpr, value)
}

// setFunctionBangHandler is the custom BangHandler for the set function.
// Normally, function calls like "(append $obj.list 1)" would not return the entire object, but only
// the function result (in this case, a vector with one more element added to it). This is because
// for most functions, the path expression is evaluated before the argument is created and the
// append function is called (i.e. append doesn't even see that its first argument originates from
// $obj.list).
// If set behaved the same way, "(set $obj.value 42)" would return 42, not the entire object. That
// is not really helpful though. If you wanted to take an object and just update one value in it,
// you'd be forced to use a temporary variable ("(set! $o ....) (set! $o.value 42) $o").
// Because of this, set returns the whole updated data structure, so in the example above, the
// entire $obj. It can do this because the first argument is not pre-evaluated by the Rudi runtime,
// but passed as a raw expression (like for delete).
func setFunctionBangHandler(ctx types.Context, originalArgs []ast.Expression) (any, error) {
	if len(originalArgs) == 0 {
		return nil, errors.New("must have at least 1 symbol argument")
	}

	firstArg := originalArgs[0]
	dest, ok := firstArg.(ast.Symbol)
	if !ok {
		return nil, fmt.Errorf("must use Symbol as first argument, got %T", firstArg)
	}

	// prepare the path expression
	var pathExpr jsonpath.Path
	if dest.PathExpression != nil {
		var err error

		pathExpr, err = pathexpr.ToJSONPath(ctx, dest.PathExpression)
		if err != nil {
			return nil, fmt.Errorf("invalid path expression: %w", err)
		}
	}

	// evaluate the new value
	newValue, err := ctx.Runtime().EvalExpression(ctx, originalArgs[1])
	if err != nil {
		return nil, err
	}

	// get the current value of the symbol
	var currentValue any
	if dest.Variable != nil {
		varName := string(*dest.Variable)

		// a non-existing variable is fine, this is how new variables are defined in the first place
		currentValue, _ = ctx.GetVariable(varName)
	} else {
		currentValue = ctx.GetDocument().Data()
	}

	patchedValue, err := jsonpath.Set(currentValue, pathExpr, newValue)
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

// (patch VAR:Variable FUNC:Identifier)
// (patch EXPR:PathExpression FUNC:Identifier)
func patchIdentifierFunction(ctx types.Context, target ast.Expression, fun ast.Identifier) (any, error) {
	return _patchIdentifierFunction(ctx, true, target, fun)
}

func _patchIdentifierFunction(ctx types.Context, clone bool, target ast.Expression, fun ast.Identifier) (any, error) {
	funcName := fun.Name
	function, ok := ctx.GetFunction(funcName)
	if !ok {
		return nil, fmt.Errorf("unknown function %s", funcName)
	}

	return patchFunction(ctx, target, clone, func(_ bool, _ any, val any) (any, error) {
		return function.Evaluate(ctx, []ast.Expression{ast.Shim{Value: val}})
	})
}

// (patch VAR:Variable NAMING:Vector BODY:Expression)
// (patch EXPR:PathExpression NAMING:Vector BODY:Expression)
func patchExpressionFunction(ctx types.Context, target ast.Expression, namingVec ast.VectorNode, expr ast.Expression) (any, error) {
	return _patchExpressionFunction(ctx, true, target, namingVec, expr)
}

func _patchExpressionFunction(ctx types.Context, clone bool, target ast.Expression, namingVec ast.VectorNode, expr ast.Expression) (any, error) {
	keyVarName, valueVarName, err := helper.DecodeNamingVector(namingVec)
	if err != nil {
		return nil, fmt.Errorf("argument #1: not a valid naming vector: %w", err)
	}

	return patchFunction(ctx, target, clone, func(exists bool, key any, val any) (any, error) {
		vars := map[string]any{
			valueVarName: val,
		}

		if keyVarName != "" {
			vars[keyVarName] = key
		}

		return ctx.Runtime().EvalExpression(ctx.NewShallowScope(nil, vars), expr)
	})
}

func patchFunction(ctx types.Context, target ast.Expression, clone bool, patch jsonpath.PatchFunc) (any, error) {
	pathed, ok := target.(ast.Pathed)
	if !ok {
		return nil, fmt.Errorf("expected datatype that can hold a path expression, got %T", target)
	}

	// separate base value expression from the path expression
	pathExpr := pathed.GetPathExpression()

	// Patch calls without a path do not make sense:
	// (patch $var to-upper) is the same as (to-upper $var)
	switch target.(type) {
	case ast.Symbol, ast.Tuple, ast.VectorNode, ast.ObjectNode:
		if pathExpr == nil {
			return nil, errors.New("no path expression provided on target value")
		}

	default:
		return nil, fmt.Errorf("unexpected target value of type %T", target)
	}

	// evaluate the target
	targetValue, err := ctx.Runtime().EvalExpression(ctx, pathed.Pathless())
	if err != nil {
		return nil, err
	}

	// We need to operate on a _copy_ of the value, as updating in-place must happen
	// via bang handler and usually (patch) is not used in those cases, as it's
	// redundant (the bang handler does the same behaviour for the first argument
	// than (patch) does).
	if clone {
		if _, ok := target.(ast.Symbol); ok {
			targetValue, err = deepcopy.Clone(targetValue)
			if err != nil {
				return nil, fmt.Errorf("invalid current value: %w", err)
			}
		}
	}

	// prepare path expression
	jsonpathExpr, err := pathexpr.ToJSONPath(ctx, pathExpr)
	if err != nil {
		return nil, fmt.Errorf("invalid path expression: %w", err)
	}

	return jsonpath.Patch(targetValue, jsonpathExpr, patch)
}

// patchFunctionBangHandler is handling (patch! ...) calls, which are redundant and would throw
// an error if the regular bang handler were to take over.
func patchFunctionBangHandler(ctx types.Context, originalArgs []ast.Expression) (any, error) {
	firstArg := originalArgs[0]
	dest, ok := firstArg.(ast.Symbol)
	if !ok {
		return nil, fmt.Errorf("must use Symbol as first argument, got %T", firstArg)
	}

	switch len(originalArgs) {
	case 2:
		ident, ok := originalArgs[1].(ast.Identifier)
		if !ok {
			return nil, fmt.Errorf("argument #1: expected identifier, got %T", originalArgs[1])
		}

		return _patchIdentifierFunction(ctx, false, dest, ident)

	case 3:
		namingVec, ok := originalArgs[1].(ast.VectorNode)
		if !ok {
			return nil, fmt.Errorf("argument #1: expected vector, got %T", originalArgs[1])
		}

		return _patchExpressionFunction(ctx, false, dest, namingVec, originalArgs[2])

	default:
		return nil, fmt.Errorf("expected 2 or 3 arguments, got %d", len(originalArgs))
	}
}

// (delete TARGET:Pathed)
func deleteFunction(ctx types.Context, target ast.Expression) (any, error) {
	pathed, ok := target.(ast.Pathed)
	if !ok {
		return nil, fmt.Errorf("argument is not a type with path expression, but %T", target)
	}

	// separate base value expression from the path expression
	pe := pathed.GetPathExpression()
	if pe == nil {
		return nil, errors.New("empty path expression")
	}

	// evaluate the target
	targetValue, err := ctx.Runtime().EvalExpression(ctx, pathed.Pathless())
	if err != nil {
		return nil, err
	}

	// we need to operate on a _copy_ of the value, as updating happens in the bang handler later
	// on; this is only necessary for symbols though, as functions (tuples) are expected to return
	// non-pointer data and vector/objects nodes are literals.
	if _, ok := target.(ast.Symbol); ok {
		targetValue, err = deepcopy.Clone(targetValue)
		if err != nil {
			return nil, fmt.Errorf("invalid current value: %w", err)
		}
	}

	// prepare the path
	pathExpr, err := pathexpr.ToJSONPath(ctx, pe)
	if err != nil {
		return nil, fmt.Errorf("invalid path expression: %w", err)
	}

	// delete the desired path in the value
	return jsonpath.Delete(targetValue, pathExpr)
}

// deleteFunctionBangHandler is the custom BangHandler for the delete function. It works
// conceptionally the same as the custom bang handler for the set function.
func deleteFunctionBangHandler(ctx types.Context, originalArgs []ast.Expression) (any, error) {
	if len(originalArgs) == 0 {
		return nil, errors.New("must have at least 1 symbol argument")
	}

	firstArg := originalArgs[0]
	dest, ok := firstArg.(ast.Symbol)
	if !ok {
		return nil, fmt.Errorf("must use Symbol as first argument, got %T", firstArg)
	}

	// prepare the path expression
	var pathExpr jsonpath.Path
	if dest.PathExpression != nil {
		var err error

		pathExpr, err = pathexpr.ToJSONPath(ctx, dest.PathExpression)
		if err != nil {
			return nil, fmt.Errorf("invalid path expression: %w", err)
		}
	}

	// get the current value of the symbol
	var currentValue any
	if dest.Variable != nil {
		varName := string(*dest.Variable)

		// a non-existing variable is fine, this is how new variables are defined in the first place
		currentValue, _ = ctx.GetVariable(varName)
	} else {
		currentValue = ctx.GetDocument().Data()
	}

	patchedValue, err := jsonpath.Delete(currentValue, pathExpr)
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

func isEmptyFunction(val bool) (any, error) {
	return !val, nil
}

func errorFunction(message string) (any, error) {
	return nil, errors.New(message)
}

func fmtErrorFunction(format string, args ...any) (any, error) {
	return nil, fmt.Errorf(format, args...)
}
