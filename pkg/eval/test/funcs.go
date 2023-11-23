// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

var (
	dummyFunctions = types.Functions{
		"eval": func(ctx types.Context, args []ast.Expression) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("expected 1 argument, got %d", len(args))
			}

			_, result, err := eval.EvalExpression(ctx, args[0])

			return result, err
		},
		// Funny enough, due to the way functions work in Rudi, "set" does not
		// actually set anything, it relies on the function magic behind the
		// scenes to handle the bang modifier.
		"set": func(ctx types.Context, args []ast.Expression) (any, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("expected 2 arguments, got %d", len(args))
			}

			_, value, err := eval.EvalExpression(ctx, args[1])
			if err != nil {
				return nil, err
			}

			return value, nil
		},
	}
)
