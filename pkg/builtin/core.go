// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/deepcopy"
	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/pathexpr"
)

// (if COND:Expr YES:Expr NO:Expr?)
func ifFunction(ctx types.Context, args []ast.Expression) (any, error) {
	_, condition, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, fmt.Errorf("condition: %w", err)
	}

	success, err := ctx.Coalesce().ToBool(condition)
	if err != nil {
		return nil, fmt.Errorf("condition: %w", err)
	}

	if success {
		// discard context changes from the true path
		_, result, err := eval.EvalExpression(ctx, args[1])
		return result, err
	}

	// optional else part
	if len(args) > 2 {
		// discard context changes from the false path
		_, result, err := eval.EvalExpression(ctx, args[2])
		return result, err
	}

	return nil, nil
}

// (do STEP:Expr+)
func doFunction(ctx types.Context, args []ast.Expression) (any, error) {
	tupleCtx := ctx

	var (
		result any
		err    error
	)

	// do not use evalArgs(), as we want to inherit the context between expressions
	for i, arg := range args {
		tupleCtx, result, err = eval.EvalExpression(tupleCtx, arg)
		if err != nil {
			return nil, fmt.Errorf("argument #%d: %w", i, err)
		}
	}

	return result, nil
}

// (has? SYM:SymbolWithPathExpression)
func hasFunction(ctx types.Context, args []ast.Expression) (any, error) {
	var (
		expr     ast.Expression
		pathExpr *ast.PathExpression
	)

	// separate base value expression from the path expression

	if symbol, ok := args[0].(ast.Symbol); ok {
		pathExpr = symbol.PathExpression

		if symbol.Variable != nil {
			symbol.PathExpression = nil
		} else {
			// for bare path expressions
			symbol.PathExpression = &ast.PathExpression{}
		}

		expr = symbol
	}

	if vectorNode, ok := args[0].(ast.VectorNode); ok {
		pathExpr = vectorNode.PathExpression
		vectorNode.PathExpression = nil
		expr = vectorNode
	}

	if objectNode, ok := args[0].(ast.ObjectNode); ok {
		pathExpr = objectNode.PathExpression
		objectNode.PathExpression = nil
		expr = objectNode
	}

	if tuple, ok := args[0].(ast.Tuple); ok {
		pathExpr = tuple.PathExpression
		tuple.PathExpression = nil
		expr = tuple
	}

	if expr == nil {
		return nil, fmt.Errorf("expected Symbol, Vector, Object or Tuple, got %T", args[0])
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
func defaultFunction(ctx types.Context, args []ast.Expression) (any, error) {
	_, result, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	// this function purposefully always uses humane coalescing
	boolified, err := coalescing.NewHumane().ToBool(result)
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	if boolified {
		return result, nil
	}

	_, result, err = eval.EvalExpression(ctx, args[1])
	if err != nil {
		return nil, fmt.Errorf("argument #1: %w", err)
	}

	return result, nil
}

// (try TEST:Expression FALLBACK:any?)
func tryFunction(ctx types.Context, args []ast.Expression) (any, error) {
	_, result, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		if len(args) == 1 {
			return nil, nil
		}

		_, result, err = eval.EvalExpression(ctx, args[1])
		if err != nil {
			return nil, fmt.Errorf("argument #1: %w", err)
		}
	}

	return result, nil
}

// (set VAR:Variable VALUE:any)
// (set EXPR:PathExpression VALUE:any)
func setFunction(ctx types.Context, args []ast.Expression) (any, error) {
	symbol, ok := args[0].(ast.Symbol)
	if !ok {
		return nil, fmt.Errorf("argument #0 is not a symbol, but %T", args[0])
	}

	// catch symbols that are technically invalid
	if symbol.Variable == nil && symbol.PathExpression == nil {
		return nil, fmt.Errorf("argument #0: must be path expression or variable, got %s", symbol.ExpressionName())
	}

	// discard any context changes within the newValue expression
	_, newValue, err := eval.EvalExpression(ctx, args[1])
	if err != nil {
		return nil, fmt.Errorf("argument #1: %w", err)
	}

	// Set relies entirely on the bang modifier handling to actually set values
	// in variables or the global document; without the bang modifier, (set)
	// is basically a no-op.

	return newValue, nil
}

// (delete VAR:Variable)
// (delete EXPR:PathExpression)
type deleteFunction struct{}

func (deleteFunction) Description() string {
	return "removes a key from an object or an item from a vector"
}

func (deleteFunction) Evaluate(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	symbol, ok := args[0].(ast.Symbol)
	if !ok {
		return nil, fmt.Errorf("argument #0 is not a symbol, but %T", args[0])
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

func (deleteFunction) BangHandler(ctx types.Context, sym ast.Symbol, value any) (types.Context, any, error) {
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

// (empty? VALUE:any)
func isEmptyFunction(ctx types.Context, args []any) (any, error) {
	// this function purposefully always uses humane coalescing
	boolified, err := coalescing.NewHumane().ToBool(args[0])
	if err != nil {
		return nil, err
	}

	return !boolified, nil
}

// (error MSG:string)
// (error FMT:string ARGS+)
func errorFunction(ctx types.Context, args []any) (any, error) {
	format, err := ctx.Coalesce().ToString(args[0])
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf(format, args[1:]...)
}

// (strictly EXPR+)
func strictlyFunction(ctx types.Context, args []ast.Expression) (any, error) {
	return coalescingChangerFunction(ctx, args, coalescing.NewStrict())
}

// (humanely EXPR+)
func humanelyFunction(ctx types.Context, args []ast.Expression) (any, error) {
	return coalescingChangerFunction(ctx, args, coalescing.NewHumane())
}

// (pedantically EXPR+)
func pedanticallyFunction(ctx types.Context, args []ast.Expression) (any, error) {
	return coalescingChangerFunction(ctx, args, coalescing.NewPedantic())
}

func coalescingChangerFunction(ctx types.Context, args []ast.Expression, c coalescing.Coalescer) (any, error) {
	_, result, err := eval.EvalExpression(ctx.WithCoalescer(c), args[0])
	return result, err
}
