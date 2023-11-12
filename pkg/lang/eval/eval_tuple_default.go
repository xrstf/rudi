// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"errors"
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
)

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
