// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package core

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/deepcopy"
	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/functions"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/pathexpr"
)

var (
	strictCoalescer   = coalescing.NewStrict()
	pedanticCoalescer = coalescing.NewPedantic()
	humaneCoalescer   = coalescing.NewHumane()

	Functions = types.Functions{
		"default": functions.NewBuilder(defaultFunction).WithDescription("returns the default value if the first argument is empty").Build(),
		"delete":  functions.NewBuilder(deleteFunction).WithBangHandler(deleteBangHandler).WithDescription("removes a key from an object or an item from a vector").Build(),
		"do":      functions.NewBuilder(DoFunction).WithDescription("eval a sequence of statements where only one expression is valid").Build(),
		"empty?":  functions.NewBuilder(isEmptyFunction).WithCoalescer(humaneCoalescer).WithDescription("returns true when the given value is empty-ish (0, false, null, \"\", ...)").Build(),
		"error":   functions.NewBuilder(errorFunction, fmtErrorFunction).WithDescription("returns an error").Build(),
		"has?":    functions.NewBuilder(hasFunction).WithDescription("returns true if the given symbol's path expression points to an existing value").Build(),
		"if":      functions.NewBuilder(ifElseFunction, ifFunction).WithDescription("evaluate one of two expressions based on a condition").Build(),
		"set":     functions.NewBuilder(setFunction).WithDescription("set a value in a variable/document, only really useful with ! modifier (set!)").Build(),
		"try":     functions.NewBuilder(tryWithFallbackFunction, tryFunction).WithDescription("returns the fallback if the first expression errors out").Build(),
	}
)

func ifFunction(ctx types.Context, test bool, yes ast.Expression) (any, error) {
	if test {
		_, result, err := eval.EvalExpression(ctx, yes)
		return result, err
	}

	return nil, nil
}

func ifElseFunction(ctx types.Context, test bool, yes, no ast.Expression) (any, error) {
	if test {
		_, result, err := eval.EvalExpression(ctx, yes)
		return result, err
	}

	_, result, err := eval.EvalExpression(ctx, no)
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
		tupleCtx, result, err = eval.EvalExpression(tupleCtx, arg)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func hasFunction(ctx types.Context, arg ast.Expression) (any, error) {
	var (
		expr     ast.Expression
		pathExpr *ast.PathExpression
	)

	// separate base value expression from the path expression

	if symbol, ok := arg.(ast.Symbol); ok {
		pathExpr = symbol.PathExpression

		if symbol.Variable != nil {
			symbol.PathExpression = nil
		} else {
			// for bare path expressions
			symbol.PathExpression = &ast.PathExpression{}
		}

		expr = symbol
	}

	if vectorNode, ok := arg.(ast.VectorNode); ok {
		pathExpr = vectorNode.PathExpression
		vectorNode.PathExpression = nil
		expr = vectorNode
	}

	if objectNode, ok := arg.(ast.ObjectNode); ok {
		pathExpr = objectNode.PathExpression
		objectNode.PathExpression = nil
		expr = objectNode
	}

	if tuple, ok := arg.(ast.Tuple); ok {
		pathExpr = tuple.PathExpression
		tuple.PathExpression = nil
		expr = tuple
	}

	if expr == nil {
		return nil, fmt.Errorf("expected Symbol, Vector, Object or Tuple, got %T", arg)
	}

	if pathExpr == nil {
		return nil, errors.New("argument has no path expression")
	}

	// pre-evaluate the path
	evaluatedPath, err := eval.EvalPathExpression(ctx, pathExpr)
	if err != nil {
		return nil, fmt.Errorf("invalid path expression: %w", err)
	}

	// evaluate the base value
	_, value, err := eval.EvalExpression(ctx, expr)
	if err != nil {
		return nil, err
	}

	_, err = eval.TraverseEvaluatedPathExpression(value, *evaluatedPath)
	if err != nil {
		return false, nil
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

	_, value, err = eval.EvalExpression(ctx, fallback)
	if err != nil {
		return nil, fmt.Errorf("argument #1: %w", err)
	}

	return value, nil
}

func tryFunction(ctx types.Context, test ast.Expression) (any, error) {
	_, result, err := eval.EvalExpression(ctx, test)
	if err != nil {
		return nil, nil
	}

	return result, nil
}

func tryWithFallbackFunction(ctx types.Context, test ast.Expression, fallback ast.Expression) (any, error) {
	_, result, err := eval.EvalExpression(ctx, test)
	if err != nil {
		_, result, err = eval.EvalExpression(ctx, fallback)
		if err != nil {
			return nil, fmt.Errorf("argument #1: %w", err)
		}
	}

	return result, nil
}

// TODO: Allow (set (foo).bar 42)

// (set VAR:Variable VALUE:any)
// (set EXPR:PathExpression VALUE:any)
func setFunction(ctx types.Context, target, value ast.Expression) (any, error) {
	symbol, ok := target.(ast.Symbol)
	if !ok {
		return nil, fmt.Errorf("argument #0 is not a symbol, but %T", target)
	}

	// catch symbols that are technically invalid
	if symbol.Variable == nil && symbol.PathExpression == nil {
		return nil, fmt.Errorf("argument #0: must be path expression or variable, got %s", symbol.ExpressionName())
	}

	// discard any context changes within the newValue expression
	_, newValue, err := eval.EvalExpression(ctx, value)
	if err != nil {
		return nil, fmt.Errorf("argument #1: %w", err)
	}

	// Set relies entirely on the bang modifier handling to actually set values
	// in variables or the global document; without the bang modifier, (set)
	// is basically a no-op.

	return newValue, nil
}

// TODO: Shouldn't (delete (map $obj to-upper).key) also work, i.e. not just
// symbols? Symbols are only important for the bang handler, which checks it
// independently.

// (delete VAR:Variable)
// (delete EXPR:PathExpression)
func deleteFunction(ctx types.Context, expr ast.Expression) (any, error) {
	symbol, ok := expr.(ast.Symbol)
	if !ok {
		return nil, fmt.Errorf("argument #0 is not a symbol, but %T", expr)
	}

	// catch symbols that are technically invalid
	if symbol.PathExpression == nil {
		return nil, fmt.Errorf("argument #0: must be path expression, got %s", symbol.ExpressionName())
	}

	// pre-evaluate the path
	pathExpr, err := eval.EvalPathExpression(ctx, symbol.PathExpression)
	if err != nil {
		return nil, fmt.Errorf("argument #0: invalid path expression: %w", err)
	}

	// get the current value
	var currentValue any

	if symbol.Variable != nil {
		varName := string(*symbol.Variable)

		// a non-existing variable is fine, this is how you define new variables in the first place
		currentValue, _ = ctx.GetVariable(varName)
	} else {
		currentValue = ctx.GetDocument().Data()
	}

	// we need to operate on a _copy_ of the value and then, if need be, rely on the BangHandler
	// to make the actual deletion happen and stick.
	currentValue, err = deepcopy.Clone(currentValue)
	if err != nil {
		return nil, fmt.Errorf("invalid current value: %w", err)
	}

	// delete the desired path in the value
	updatedValue, err := pathexpr.Delete(currentValue, pathexpr.FromEvaluatedPath(*pathExpr))
	if err != nil {
		return nil, fmt.Errorf("cannot delete %s in %T: %w", pathExpr, currentValue, err)
	}

	return updatedValue, nil
}

func deleteBangHandler(ctx types.Context, originalArgs []ast.Expression, value any) (types.Context, any, error) {
	if len(originalArgs) == 0 {
		return ctx, nil, errors.New("must have at least 1 symbol argument")
	}

	firstArg := originalArgs[0]
	sym, ok := firstArg.(ast.Symbol)
	if !ok {
		return ctx, nil, fmt.Errorf("must use Symbol as first argument, got %T", firstArg)
	}

	updatedValue := value

	// if the symbol has a path to traverse, do so
	if sym.PathExpression != nil {
		// pre-evaluate the path expression
		pathExpr, err := eval.EvalPathExpression(ctx, sym.PathExpression)
		if err != nil {
			return ctx, nil, fmt.Errorf("argument #0: invalid path expression: %w", err)
		}

		// get the current value of the symbol
		var currentValue any

		if sym.Variable != nil {
			varName := string(*sym.Variable)

			// a non-existing variable is fine, this is how you define new variables in the first place
			currentValue, _ = ctx.GetVariable(varName)
		} else {
			currentValue = ctx.GetDocument().Data()
		}

		// apply the path expression
		updatedValue, err = pathexpr.Delete(currentValue, pathexpr.FromEvaluatedPath(*pathExpr))
		if err != nil {
			return ctx, nil, fmt.Errorf("cannot set value in %T at %s: %w", currentValue, pathExpr, err)
		}
	}

	if sym.Variable != nil {
		varName := string(*sym.Variable)
		ctx = ctx.WithVariable(varName, updatedValue)
	} else {
		ctx.GetDocument().Set(updatedValue)
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
