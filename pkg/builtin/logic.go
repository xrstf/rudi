// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func andFunction(ctx types.Context, args []ast.Expression) (any, error) {
	for i, arg := range args {
		_, evaluated, err := eval.EvalExpression(ctx, arg)
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

func orFunction(ctx types.Context, args []ast.Expression) (any, error) {
	for i, arg := range args {
		_, evaluated, err := eval.EvalExpression(ctx, arg)
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

func notFunction(ctx types.Context, args []any) (any, error) {
	arg, err := ctx.Coalesce().ToBool(args[0])
	if err != nil {
		return nil, err
	}

	return !arg, nil
}
