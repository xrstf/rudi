// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"errors"
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
)

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
