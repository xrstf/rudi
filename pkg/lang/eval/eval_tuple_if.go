// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/coalescing"
)

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
