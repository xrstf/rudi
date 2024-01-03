// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package logic

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/functions"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

var (
	Functions = types.Functions{
		"and": functions.NewBuilder(andFunction).WithDescription("returns true if all arguments are true").Build(),
		"or":  functions.NewBuilder(orFunction).WithDescription("returns true if any of the arguments is true").Build(),
		"not": functions.NewBuilder(notFunction).WithDescription("negates the given argument").Build(),
	}
)

func andFunction(ctx types.Context, args ...ast.Expression) (any, error) {
	for i, arg := range args {
		evaluated, err := ctx.Runtime().EvalExpression(ctx, arg)
		if err != nil {
			return nil, fmt.Errorf("argument #%d: %w", i, err)
		}

		part, err := ctx.Coalesce().ToBool(evaluated)
		if err != nil {
			return nil, fmt.Errorf("argument #%d: %w", i, err)
		}

		if !part {
			return false, nil
		}
	}

	return true, nil
}

func orFunction(ctx types.Context, args ...ast.Expression) (any, error) {
	for i, arg := range args {
		evaluated, err := ctx.Runtime().EvalExpression(ctx, arg)
		if err != nil {
			return nil, fmt.Errorf("argument #%d: %w", i, err)
		}

		part, err := ctx.Coalesce().ToBool(evaluated)
		if err != nil {
			return nil, fmt.Errorf("argument #%d: %w", i, err)
		}

		if part {
			return true, nil
		}
	}

	return false, nil
}

func notFunction(val bool) (any, error) {
	return !val, nil
}
