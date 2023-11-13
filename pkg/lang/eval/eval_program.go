// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

func evalProgram(ctx types.Context, p *ast.Program) (any, error) {
	innerCtx := ctx

	// This is all sorts of wonky and not really how the program execution should work.
	// But it compiles.
	var (
		result any
		err    error
	)

	for _, stmt := range p.Statements {
		innerCtx, result, err = evalStatement(innerCtx, &stmt)
		if err != nil {
			return nil, fmt.Errorf("failed to eval statement %s: %w", stmt.String(), err)
		}
	}

	return result, nil
}
