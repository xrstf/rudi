// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"errors"
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/builtin"
	"go.xrstf.de/corel/pkg/lang/eval/coalescing"
)

func evalTuple(ctx Context, tup *ast.Tuple) (Context, interface{}, error) {
	funcName := tup.Identifier.Name

	// hardcode root behaviour for those tuples where not all
	// expressions can be pre-computed (in case, for example,
	// the else-path of an if statement would have side effects)
	switch funcName {
	case "if":
		return evalIfTuple(ctx, tup)
	case "set":
		return evalSetTuple(ctx, tup)
	case "do":
		return evalDoTuple(ctx, tup)
	}

	function, ok := builtin.Functions[funcName]
	if !ok {
		return ctx, nil, fmt.Errorf("unknown function %s", funcName)
	}

	// evaluate all function arguments
	args := make([]interface{}, len(tup.Expressions))
	for i, expr := range tup.Expressions {
		// each function arg on its own cannot change the overall context, so discard it
		_, arg, err := evalExpression(ctx, &expr)
		if err != nil {
			return ctx, nil, fmt.Errorf("invalid argument %d: %w", i, err)
		}

		args[i] = arg
	}

	// call the function
	result, err := function(args)
	if err != nil {
		return ctx, nil, fmt.Errorf("function failed: %w", err)
	}

	return ctx, result, nil
}

func evalIfTuple(ctx Context, tup *ast.Tuple) (Context, interface{}, error) {
	if size := len(tup.Expressions); size != 2 && size != 3 {
		return ctx, nil, fmt.Errorf("invalid if tuple: expected 2 or 3 expressions, but got %d", size)
	}

	tupleCtx := ctx

	tupleCtx, condition, err := evalExpression(tupleCtx, &tup.Expressions[0])
	if err != nil {
		return ctx, nil, fmt.Errorf("failed to evaluate condition: %w", err)
	}

	success, err := coalescing.ToBool(condition)
	if err != nil {
		return ctx, nil, fmt.Errorf("condition did not return boolish value: %w", err)
	}

	if success {
		// discard context changes from the true path
		_, result, err := evalExpression(tupleCtx, &tup.Expressions[1])
		return ctx, result, err
	}

	// optional else part
	if len(tup.Expressions) > 2 {
		// discard context changes from the false path
		_, result, err := evalExpression(tupleCtx, &tup.Expressions[2])
		return ctx, result, err
	}

	return ctx, nil, nil
}

func evalSetTuple(ctx Context, tup *ast.Tuple) (Context, interface{}, error) {
	if size := len(tup.Expressions); size != 2 {
		return ctx, nil, errors.New("(set $varname EXPRESSION)")
	}

	varNameExpr := tup.Expressions[0]
	varName := ""

	if varNameExpr.SymbolNode == nil {
		return ctx, nil, errors.New("(set $varname EXPRESSION)")
	}

	symNode := varNameExpr.SymbolNode
	if symNode.Variable == nil {
		return ctx, nil, errors.New("(set $varname EXPRESSION)")
	}

	// forbid weird definitions like (set $var.foo (expr)) for now
	if symNode.PathExpression != nil {
		return ctx, nil, errors.New("(set $varname EXPRESSION)")
	}

	varName = symNode.Variable.Name

	// discard any context changes within the value expression
	_, value, err := evalExpression(ctx, &tup.Expressions[1])
	if err != nil {
		return ctx, nil, fmt.Errorf("failed to evaluate variable value: %w", err)
	}

	// make the variable's value the return value, so `(def $foo 12)` = 12
	return ctx.WithVariable(varName, value), value, nil
}

func evalDoTuple(ctx Context, tup *ast.Tuple) (Context, interface{}, error) {
	if size := len(tup.Expressions); size < 2 {
		return ctx, nil, errors.New("(do EXPRESSION+)")
	}

	innerCtx := ctx

	var (
		result interface{}
		err    error
	)

	for _, expr := range tup.Expressions {
		innerCtx, result, err = evalExpression(innerCtx, &expr)
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to eval expression %s: %w", expr.String(), err)
		}
	}

	// make the variable's value the return value, so `(def $foo 12)` = 12
	return ctx, result, nil
}
