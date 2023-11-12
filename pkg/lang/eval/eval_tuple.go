// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"errors"
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/builtin"
	"go.xrstf.de/corel/pkg/lang/eval/types"
)

func evalTuple(ctx types.Context, tup *ast.Tuple) (types.Context, interface{}, error) {
	if len(tup.Expressions) == 0 {
		return ctx, nil, errors.New("invalid tuple: tuple cannot be empty")
	}

	funcExpr := tup.Expressions[0]
	if funcExpr.IdentifierNode == nil {
		return ctx, nil, errors.New("invalid tuple: first expression must be an identifier")
	}

	funcName := funcExpr.IdentifierNode.Name
	argExprs := tup.Expressions[1:]

	// wrap all args to allow the builtin package to use code from this
	// package without causing a circular dependency
	funcArgs := make([]builtin.Argument, len(argExprs))
	for i := range argExprs {
		funcArgs[i] = &argument{
			expression: &argExprs[i],
		}
	}

	function, ok := builtin.Functions[funcName]
	if !ok {
		return ctx, nil, fmt.Errorf("unknown function %s", funcName)
	}

	// call the function
	newContext, result, err := function(ctx, funcArgs)
	if err != nil {
		return ctx, nil, fmt.Errorf("%s: %w", funcName, err)
	}

	return newContext, result, nil
}

type argument struct {
	expression *ast.Expression
}

var _ builtin.Argument = &argument{}

func (a *argument) String() string {
	return a.expression.String()
}

func (a *argument) Eval(ctx types.Context) (types.Context, any, error) {
	return evalExpression(ctx, a.expression)
}

func (a *argument) Expression() *ast.Expression {
	return a.expression
}
