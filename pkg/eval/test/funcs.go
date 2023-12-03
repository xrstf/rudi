// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/eval/util"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

var (
	dummyFunctions = types.Functions{
		"eval": util.NewLiteralFunction(func(ctx types.Context, args []any) (any, error) {
			return args[0], nil
		}, "").MinArgs(1).MaxArgs(1),

		// Funny enough, due to the way functions work in Rudi, "set" does not
		// actually set anything, it relies on the function magic behind the
		// scenes to handle the bang modifier.
		"set": util.NewRawFunction(func(ctx types.Context, args []ast.Expression) (any, error) {
			_, value, err := eval.EvalExpression(ctx, args[1])

			return value, err
		}, "").MinArgs(2).MaxArgs(2),
	}
)
