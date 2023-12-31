// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package interpreter

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

// Run implements types.Runtime.
func (i *interpreter) EvalProgram(ctx types.Context, p *ast.Program) (types.Context, any, error) {
	if p == nil {
		return ctx, nil, errors.New("program is nil")
	}

	if len(p.Statements) == 0 {
		return ctx, nil, nil
	}

	innerCtx := ctx

	var (
		result any
		err    error
	)

	for _, stmt := range p.Statements {
		innerCtx, result, err = i.EvalStatement(innerCtx, stmt)
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to eval statement %s: %w", stmt.String(), err)
		}
	}

	return innerCtx, result, nil
}
