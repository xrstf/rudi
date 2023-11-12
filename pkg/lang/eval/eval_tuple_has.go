// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"errors"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalHasTuple(ctx Context, args []ast.Expression) (Context, interface{}, error) {
	if size := len(args); size != 1 {
		return ctx, nil, errors.New("(has PATH_EXPRESSION)")
	}

	arg := args[0]
	if arg.SymbolNode == nil || arg.SymbolNode.PathExpression == nil {
		return ctx, nil, errors.New("(has PATH_EXPRESSION)")
	}

	_, value, err := evalSymbol(ctx, arg.SymbolNode)
	if err != nil {
		return ctx, false, nil
	}

	return ctx, value != nil, nil
}
