// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"
	"time"

	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func nowFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 arguments, got %d", size)
	}

	_, format, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	formatString, err := ctx.Coalesce().ToString(format)
	if err != nil {
		return nil, err
	}

	formatted := time.Now().Format(string(formatString))

	return ast.String(formatted), nil
}
