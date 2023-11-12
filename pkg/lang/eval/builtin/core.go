// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"
	"fmt"

	"go.xrstf.de/corel/pkg/lang/eval/coalescing"
	"go.xrstf.de/corel/pkg/lang/eval/types"
)

// (if COND:Expr YES:Expr NO:Expr?)
func ifFunction(ctx types.Context, args []Argument) (any, error) {
	if size := len(args); size < 2 || size > 3 {
		return nil, fmt.Errorf("expected 2 or 3 arguments, got %d", size)
	}

	tupleCtx := ctx

	tupleCtx, condition, err := args[0].Eval(tupleCtx)
	if err != nil {
		return nil, fmt.Errorf("condition: %w", err)
	}

	success, err := coalescing.ToBool(condition)
	if err != nil {
		return nil, fmt.Errorf("condition is not boolish: %w", err)
	}

	if success {
		// discard context changes from the true path
		_, result, err := args[1].Eval(tupleCtx)
		return result, err
	}

	// optional else part
	if len(args) > 2 {
		// discard context changes from the false path
		_, result, err := args[2].Eval(tupleCtx)
		return result, err
	}

	return nil, nil
}

// (do STEP:Expr+)
func doFunction(ctx types.Context, args []Argument) (any, error) {
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
		tupleCtx, result, err = arg.Eval(tupleCtx)
		if err != nil {
			return nil, fmt.Errorf("argument #%d: %w", i, err)
		}
	}

	return result, nil
}

// (has PATH:PathExpression)
func hasFunction(ctx types.Context, args []Argument) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	arg := args[0].Expression()
	if arg.SymbolNode == nil || arg.SymbolNode.PathExpression == nil {
		return nil, fmt.Errorf("argument #0 is not a path expression, but %T", arg) // TODO: this generates just "Expression", make it nicer
	}

	_, value, err := args[0].Eval(ctx)
	if err != nil {
		return false, nil
	}

	return value != nil, nil
}

// (default TEST:Expression FALLBACK:any)
func defaultFunction(ctx types.Context, args []Argument) (any, error) {
	if size := len(args); size != 2 {
		return nil, fmt.Errorf("expected 2 arguments, got %d", size)
	}

	_, result, err := args[0].Eval(ctx)
	if err != nil {
		_, result, err = args[1].Eval(ctx)
		if err != nil {
			return nil, fmt.Errorf("argument #1: %w", err)
		}

		return false, nil
	}

	return result, nil
}

// (set VAR:Variable VALUE:any)
// (set EXPR:PathExpression VALUE:any) <- TODO
func setFunction(ctx types.Context, args []Argument) (types.Context, any, error) {
	if size := len(args); size != 2 {
		return ctx, nil, fmt.Errorf("expected 2 arguments, got %d", size)
	}

	varNameExpr := args[0].Expression()
	varName := ""

	if varNameExpr.SymbolNode == nil {
		return ctx, nil, fmt.Errorf("argument #0: must be path expression or variable, got %T", varNameExpr)
	}

	symNode := varNameExpr.SymbolNode
	if symNode.Variable == nil && symNode.PathExpression == nil {
		return ctx, nil, fmt.Errorf("argument #0: must be path expression or variable, got %T", varNameExpr)
	}

	// set a variable, which will result in a new context
	if symNode.Variable != nil {
		// forbid weird definitions like (set $var.foo (expr)) for now
		if symNode.PathExpression != nil {
			return ctx, nil, errors.New("argument #0: cannot use path expression when setting variable values")
		}

		varName = symNode.Variable.Name

		// discard any context changes within the value expression
		_, value, err := args[1].Eval(ctx)
		if err != nil {
			return ctx, nil, fmt.Errorf("argument #1: %w", err)
		}

		// make the variable's value the return value, so `(def $foo 12)` = 12
		return ctx.WithVariable(varName, value), value, nil
	}

	return ctx, nil, errors.New("setting a document path expression is not yet implemented")
}
