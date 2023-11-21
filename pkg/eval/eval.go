// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func Run(ctx types.Context, p *ast.Program) (types.Context, any, error) {
	newCtx, result, err := EvalProgram(ctx, p)
	if err != nil {
		return ctx, nil, err
	}

	unwrapped, err := types.UnwrapType(result)
	if err != nil {
		return ctx, nil, err
	}

	return newCtx, unwrapped, nil
}
