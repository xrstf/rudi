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
		"eval": types.NewFunction(func(ctx types.Context, args []ast.Expression) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("expected 1 argument, got %d", len(args))
			}

			_, value, err := eval.EvalExpression(ctx, args[0])

			return value, err
		}, "evaluates the given expression and returns its value"),

		// Funny enough, due to the way functions work in Rudi, "set" does not
		// actually set anything, it relies on the function magic behind the
		// scenes to handle the bang modifier.
		"set": types.NewFunction(func(ctx types.Context, args []ast.Expression) (any, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("expected 2 arguments, got %d", len(args))
			}

			_, value, err := eval.EvalExpression(ctx, args[1])

			return value, err
		}, "sets a variable or accesses the global document, most often used with the bang modifier"),
	}
)
