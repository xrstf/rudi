// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"errors"
	"fmt"

	"go.xrstf.de/otto/pkg/eval/types"
	"go.xrstf.de/otto/pkg/lang/ast"
)

func EvalProgram(ctx types.Context, p *ast.Program) (types.Context, any, error) {
	if p == nil {
		return ctx, nil, errors.New("program is nil")
	}

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
