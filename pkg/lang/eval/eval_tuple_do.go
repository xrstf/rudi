// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"errors"
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
)

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
