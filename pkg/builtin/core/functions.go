// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package core

import (
	"context"
	"errors"
	"fmt"

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
		"delete":  functions.NewBuilder(deleteFunction).WithBangHandler(overwriteEverythingBangHandler).WithDescription("removes a key from an object or an item from a vector").Build(),
		"do":      functions.NewBuilder(DoFunction).WithDescription("eval a sequence of statements where only one expression is valid").Build(),
		"empty?":  functions.NewBuilder(isEmptyFunction).WithCoalescer(humaneCoalescer).WithDescription("returns true when the given value is empty-ish (0, false, null, \"\", ...)").Build(),
		"error":   functions.NewBuilder(errorFunction, fmtErrorFunction).WithDescription("returns an error").Build(),
		"has?":    functions.NewBuilder(hasFunction).WithDescription("returns true if the given symbol's path expression points to an existing value").Build(),
		"if":      functions.NewBuilder(ifElseFunction, ifFunction).WithDescription("evaluate one of two expressions based on a condition").Build(),
		"set":     functions.NewBuilder(setFunction).WithBangHandler(overwriteEverythingBangHandler).WithDescription("set a value in a variable/document, only really useful with ! modifier (set!)").Build(),
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
	if test {
		_, result, err := ctx.Runtime().EvalExpression(ctx, yes)
		return result, err
	}

	return nil, nil
}

func ifElseFunction(ctx types.Context, test bool, yes, no ast.Expression) (any, error) {
	if test {
		_, result, err := ctx.Runtime().EvalExpression(ctx, yes)
		return result, err
	}

	_, result, err := ctx.Runtime().EvalExpression(ctx, no)
	return result, err
}

// NB: Variadic functions always require at least 1 argument in Rudi to match.
// This function, doing basically "nothing", is re-used by other Rudi functions and therefore
// exported.
func DoFunction(ctx types.Context, args ...ast.Expression) (any, error) {
	var (
		tupleCtx = ctx
		result   any
		err      error
	)

	for _, arg := range args {
		tupleCtx, result, err = ctx.Runtime().EvalExpression(tupleCtx, arg)
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

	// pre-evaluate the path
	evaluatedPath, err := pathexpr.Eval(ctx, pathExpr)
	if err != nil {
		return nil, fmt.Errorf("invalid path expression: %w", err)
	}

	// evaluate the base value
	_, value, err := ctx.Runtime().EvalExpression(ctx, expr)
	if err != nil {
		return nil, err
	}

	_, err = pathexpr.Traverse(value, *evaluatedPath)
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

	_, value, err = ctx.Runtime().EvalExpression(ctx, fallback)
	if err != nil {
		return nil, fmt.Errorf("argument #1: %w", err)
	}

	return value, nil
}

func tryFunction(ctx types.Context, test ast.Expression) (any, error) {
	_, result, err := ctx.Runtime().EvalExpression(ctx, test)
	if err != nil {
		return nil, keepContextCanceled(err)
	}

	return result, nil
}

func tryWithFallbackFunction(ctx types.Context, test ast.Expression, fallback ast.Expression) (any, error) {
	_, result, err := ctx.Runtime().EvalExpression(ctx, test)
	if err != nil {
		_, result, err = ctx.Runtime().EvalExpression(ctx, fallback)
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

	// pre-evaluate the path (assuming it's cheaper to calculate than the main expression)
	evaluatedPath, err := pathexpr.Eval(ctx, pathExpr)
	if err != nil {
		return nil, fmt.Errorf("invalid path expression: %w", err)
	}

	// evaluate the target
	_, targetValue, err := ctx.Runtime().EvalExpression(ctx, expr)
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

	return jsonpath.Set(targetValue, jsonpath.FromEvaluatedPath(*evaluatedPath), value)
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

	// pre-evaluate the path
	pathExpr, err := pathexpr.Eval(ctx, pe)
	if err != nil {
		return nil, fmt.Errorf("invalid path expression: %w", err)
	}

	// evaluate the target
	_, targetValue, err := ctx.Runtime().EvalExpression(ctx, pathed.Pathless())
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

	// delete the desired path in the value
	return jsonpath.Delete(targetValue, jsonpath.FromEvaluatedPath(*pathExpr))
}

// overwriteEverythingBangHandler is the custom BangHandler for the set and delete functions.
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
// All of this applies equally to the delete function, hence both share the same bang handler.
// Since this behaviour makes set different from other regular functions, it needs a custom
// BangHandler.
func overwriteEverythingBangHandler(ctx types.Context, originalArgs []ast.Expression, value any) (types.Context, any, error) {
	if len(originalArgs) == 0 {
		return ctx, nil, errors.New("must have at least 1 symbol argument")
	}

	firstArg := originalArgs[0]
	symbol, ok := firstArg.(ast.Symbol)
	if !ok {
		return ctx, nil, fmt.Errorf("must use Symbol as first argument, got %T", firstArg)
	}

	// Since set always returns the entire data structure, all we must do here is to
	// update the target, ignoring the path expression on the symbol.
	if symbol.Variable != nil {
		varName := string(*symbol.Variable)
		ctx = ctx.WithVariable(varName, value)
	} else {
		ctx.GetDocument().Set(value)
	}

	return ctx, value, nil
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
