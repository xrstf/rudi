// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/equality"
	"go.xrstf.de/otto/pkg/lang/eval"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

func eqFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if len(args) != 2 {
		return nil, errors.New("(eq LEFT RIGHT)")
	}

	_, leftData, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	_, rightData, err := eval.EvalExpression(ctx, args[1])
	if err != nil {
		return nil, fmt.Errorf("argument #1: %w", err)
	}

	equal, err := equality.StrictEqual(leftData, rightData)
	if err != nil {
		return nil, fmt.Errorf("cannot compare %T with %T: %w", leftData, rightData, err)
	}

	return ast.Bool(equal), nil
}
