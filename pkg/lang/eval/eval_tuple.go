// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"errors"
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/builtin"
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
	case "has":
		return evalHasTuple(ctx, argExprs)
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
