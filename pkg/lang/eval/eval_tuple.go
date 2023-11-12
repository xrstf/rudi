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
	if len(tup.Expressions) == 0 {
		return ctx, nil, errors.New("invalid tuple: tuple cannot be empty")
	}

	funcExpr := tup.Expressions[0]
	if funcExpr.IdentifierNode == nil {
		return ctx, nil, errors.New("invalid tuple: first expression must be an identifier")
	}

	funcName := funcExpr.IdentifierNode.Name
	argExprs := tup.Expressions[1:]

	// hardcode root behaviour for those tuples where not all
	// expressions can be pre-computed (in case, for example,
	// the else-path of an if statement would have side effects)
	switch funcName {
	case "if":
		return evalIfTuple(ctx, argExprs)
	case "set":
		return evalSetTuple(ctx, argExprs)
	case "default":
		return evalDefaultTuple(ctx, argExprs)
	case "do":
		return evalDoTuple(ctx, argExprs)
	}

	function, ok := builtin.Functions[funcName]
	if !ok {
		return ctx, nil, fmt.Errorf("unknown function %s", funcName)
	}

	// evaluate all function arguments
	args := make([]interface{}, len(argExprs))
	for i, expr := range argExprs {
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

func evalIfTuple(ctx Context, args []ast.Expression) (Context, interface{}, error) {
	if size := len(args); size != 2 && size != 3 {
		return ctx, nil, fmt.Errorf("invalid if tuple: expected 2 or 3 expressions, but got %d", size)
	}

	tupleCtx := ctx

	tupleCtx, condition, err := evalExpression(tupleCtx, &args[0])
	if err != nil {
		return ctx, nil, fmt.Errorf("failed to evaluate condition: %w", err)
	}

	success, err := coalescing.ToBool(condition)
	if err != nil {
		return ctx, nil, fmt.Errorf("condition did not return boolish value: %w", err)
	}

	if success {
		// discard context changes from the true path
		_, result, err := evalExpression(tupleCtx, &args[1])
		return ctx, result, err
	}

	// optional else part
	if len(args) > 2 {
		// discard context changes from the false path
		_, result, err := evalExpression(tupleCtx, &args[2])
		return ctx, result, err
	}

	return ctx, nil, nil
}

var setTupleSyntaxError = errors.New("(set $varname EXPRESSION) or (set PATH_EXPRESSION EXPRESSION)")

func evalSetTuple(ctx Context, args []ast.Expression) (Context, interface{}, error) {
	if size := len(args); size != 2 {
		return ctx, nil, setTupleSyntaxError
	}

	varNameExpr := args[0]
	varName := ""

	if varNameExpr.SymbolNode == nil {
		return ctx, nil, setTupleSyntaxError
	}

	symNode := varNameExpr.SymbolNode
	if symNode.Variable == nil && symNode.PathExpression == nil {
		return ctx, nil, setTupleSyntaxError
	}

	// set a variable, which will result in a new context
	if symNode.Variable != nil {
		// forbid weird definitions like (set $var.foo (expr)) for now
		if symNode.PathExpression != nil {
			return ctx, nil, setTupleSyntaxError
		}

		varName = symNode.Variable.Name

		// discard any context changes within the value expression
		_, value, err := evalExpression(ctx, &args[1])
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to evaluate variable value: %w", err)
		}

		// make the variable's value the return value, so `(def $foo 12)` = 12
		return ctx.WithVariable(varName, value), value, nil
	}

	return ctx, nil, errors.New("setting a document path expression is not yet implemented")
}

func evalDoTuple(ctx Context, args []ast.Expression) (Context, interface{}, error) {
	if size := len(args); size < 2 {
		return ctx, nil, errors.New("(do EXPRESSION+)")
	}

	innerCtx := ctx

	var (
		result interface{}
		err    error
	)

	for _, expr := range args {
		innerCtx, result, err = evalExpression(innerCtx, &expr)
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to eval expression %s: %w", expr.String(), err)
		}
	}

	return ctx, result, nil
}

func evalDefaultTuple(ctx Context, args []ast.Expression) (Context, interface{}, error) {
	if size := len(args); size != 2 {
		return ctx, nil, errors.New("(default EXPRESSION DEFAULT)")
	}

	_, result, err := evalExpression(ctx, &args[0])
	if err != nil {
		defaultExpr := args[1]

		_, result, err = evalExpression(ctx, &defaultExpr)
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to eval expression %s: %w", defaultExpr.String(), err)
		}
	}

	return ctx, result, nil
}
