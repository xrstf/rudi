// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval"
	"go.xrstf.de/otto/pkg/lang/eval/coalescing"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

// (if COND:Expr YES:Expr NO:Expr?)
func ifFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size < 2 || size > 3 {
		return nil, fmt.Errorf("expected 2 or 3 arguments, got %d", size)
	}

	tupleCtx := ctx

	tupleCtx, condition, err := eval.EvalExpression(tupleCtx, args[0])
	if err != nil {
		return nil, fmt.Errorf("condition: %w", err)
	}

	success, ok := condition.(ast.Bool)
	if !ok {
		return nil, fmt.Errorf("condition is not bool, but %T", err)
	}

	if success {
		// discard context changes from the true path
		_, result, err := eval.EvalExpression(tupleCtx, args[1])
		return result, err
	}

	// optional else part
	if len(args) > 2 {
		// discard context changes from the false path
		_, result, err := eval.EvalExpression(tupleCtx, args[2])
		return result, err
	}

	return ast.Null{}, nil
}

// (do STEP:Expr+)
func doFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size < 1 {
		return nil, fmt.Errorf("expected 1+ arguments, got %d", size)
	}

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

// (has? PATH:PathExpression)
func hasFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	symbol, ok := args[0].(ast.Symbol)
	if !ok {
		return nil, fmt.Errorf("argument #0 is not a path expression, but %s", args[0].ExpressionName())
	}

	if symbol.PathExpression == nil {
		return nil, errors.New("argument #0 has no path expression")
	}

	_, value, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return false, nil
	}

	return ast.Bool(value != nil), nil
}

// (default TEST:Expression FALLBACK:any)
func defaultFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 2 {
		return nil, fmt.Errorf("expected 2 arguments, got %d", size)
	}

	_, result, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	isEmpty, err := coalescing.IsEmpty(result)
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	if !isEmpty {
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
	if size := len(args); size < 1 || size > 2 {
		return nil, fmt.Errorf("expected 1 or 2 arguments, got %d", size)
	}

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
// (set EXPR:PathExpression VALUE:any) <- TODO
func setFunction(ctx types.Context, args []ast.Expression) (types.Context, any, error) {
	if size := len(args); size != 2 {
		return ctx, nil, fmt.Errorf("expected 2 arguments, got %d", size)
	}

	symbol, ok := args[0].(ast.Symbol)
	if !ok {
		return ctx, nil, fmt.Errorf("argument #0 is not symbol, but %s", args[0].ExpressionName())
	}

	// catch symbols that are technically invalid
	if symbol.Variable == nil && symbol.PathExpression == nil {
		return ctx, nil, fmt.Errorf("argument #0: must be path expression or variable, got %s", symbol.ExpressionName())
	}

	varName := ""

	// discard any context changes within the newValue expression
	_, newValue, err := eval.EvalExpression(ctx, args[1])
	if err != nil {
		return ctx, nil, fmt.Errorf("argument #1: %w", err)
	}

	// set a variable, which will result in a new context
	if symbol.Variable != nil {
		// forbid weird definitions like (set $var.foo (expr)) for now
		if symbol.PathExpression != nil {
			return ctx, nil, errors.New("argument #0: cannot use path expression when setting variable values")
		}

		varName = string(*symbol.Variable)

		// make the variable's value the return value, so `(def $foo 12)` = 12
		return ctx.WithVariable(varName, newValue), newValue, nil
	}

	// set new value at path expression
	doc := ctx.GetDocument()
	setValueAtPath(ctx, doc.Get(), symbol.PathExpression.Steps, newValue)

	return ctx, nil, errors.New("setting a document path expression is not yet implemented")
}

func setValueAtPath(ctx types.Context, document any, steps []ast.Expression, newValue any) (any, error) {
	if len(steps) == 0 {
		return nil, nil
	}

	return nil, nil

	// firstStep := steps[0]
	// remainingPath := steps[1:]

	// // short-circuit for expressions like (set . 42)
	// if firstStep.IsIdentity() {
	// 	return newValue, nil
	// }

	// innerCtx := ctx

	// // evaluate the current step
	// switch {
	// case firstStep.Identifier != nil:
	// 	step = ast.String(string(*firstStep.Identifier))
	// case firstStep.StringNode != nil:
	// 	step = ast.String(string(*firstStep.StringNode))
	// case firstStep.Integer != nil:
	// 	step = ast.Number{Value: *firstStep.Integer}
	// case firstStep.Variable != nil:
	// 	name := string(*firstStep.Variable)

	// 	value, ok := innerCtx.GetVariable(name)
	// 	if !ok {
	// 		return nil, fmt.Errorf("unknown variable %s (%T)", name, name), nil
	// 	}
	// 	step = value
	// case firstStep.Tuple != nil:
	// 	var (
	// 		value any
	// 		err   error
	// 	)

	// 	// keep accumulating context changes, so you _could_ in theory do
	// 	// $var[(set $bla 2)][(add $bla 2)] <-- would be $var[2][4]
	// 	innerCtx, value, err = evalTuple(innerCtx, firstStep.Tuple)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("invalid accessor: %w", err), nil
	// 	}

	// 	step = value
	// }
}

// (empty? VALUE:any)
func isEmptyFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 or 2 arguments, got %d", size)
	}

	_, result, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	switch asserted := result.(type) {
	case ast.Bool:
		return bool(asserted) == false, nil
	case ast.Number:
		return asserted.ToFloat() == 0, nil
	case ast.String:
		return len(string(asserted)) == 0, nil
	case ast.Null:
		return true, nil
	case ast.Vector:
		return len(asserted.Data) == 0, nil
	case ast.Object:
		return len(asserted.Data) == 0, nil
	case ast.Identifier:
		return nil, fmt.Errorf("unexpected identifier %s", asserted)
	default:
		return nil, fmt.Errorf("unexpected argument %v (%T)", result, result)
	}
}
