// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

var (
	dummyFunctions = types.Functions{
		"eval": func(ctx types.Context, args []ast.Expression) (types.Context, any, error) {
			if len(args) != 1 {
				return ctx, nil, fmt.Errorf("expected 1 argument, got %d", len(args))
			}

			return eval.EvalExpression(ctx, args[0])
		},
	}
)
