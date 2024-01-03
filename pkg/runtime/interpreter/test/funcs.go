// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

var (
	dummyFunctions = types.Functions{
		"eval": types.NewFunction(func(ctx types.Context, args []ast.Expression) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("expected 1 argument, got %d", len(args))
			}

			return ctx.Runtime().EvalExpression(ctx, args[0])
		}, "evaluates the given expression and returns its value"),

		// Funny enough, due to the way functions work in Rudi, "set" does not
		// actually set anything, it relies on the function magic behind the
		// scenes to handle the bang modifier.
		"set": types.NewFunction(func(ctx types.Context, args []ast.Expression) (any, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("expected 2 arguments, got %d", len(args))
			}

			return ctx.Runtime().EvalExpression(ctx, args[1])
		}, "sets a variable or accesses the global document, most often used with the bang modifier"),
	}
)
