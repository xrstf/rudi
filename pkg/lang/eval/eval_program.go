// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

func EvalProgram(ctx types.Context, p ast.Program) (types.Context, any, error) {
	innerCtx := ctx

	var (
		result any
		err    error
	)

	for _, stmt := range p.Statements {
		innerCtx, result, err = EvalStatement(innerCtx, stmt)
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to eval statement %s: %w", stmt.String(), err)
		}
	}

	return innerCtx, result, nil
}
