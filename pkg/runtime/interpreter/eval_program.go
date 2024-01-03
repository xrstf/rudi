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
func (i *interpreter) EvalProgram(ctx types.Context, p *ast.Program) (any, error) {
	if p == nil {
		return nil, errors.New("program is nil")
	}

	if len(p.Statements) == 0 {
		return nil, nil
	}

	scope := ctx.NewScope()

	var (
		result any
		err    error
	)

	for _, stmt := range p.Statements {
		result, err = i.EvalStatement(scope, stmt)
		if err != nil {
			return nil, fmt.Errorf("failed to eval statement %s: %w", stmt.String(), err)
		}
	}

	return result, nil
}
