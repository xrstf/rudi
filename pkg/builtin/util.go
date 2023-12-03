// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func evalArgs(ctx types.Context, args []ast.Expression, argShift int) ([]any, error) {
	values := make([]any, len(args)-argShift)
	for i, arg := range args[argShift:] {
		_, evaluated, err := eval.EvalExpression(ctx, arg)
		if err != nil {
			return nil, fmt.Errorf("argument #%d: %w", i+argShift, err)
		}

		values[i] = evaluated
	}

	return values, nil
}

func checkIterable(ctx types.Context, val any) error {
	if _, err := ctx.Coalesce().ToVector(val); err != nil {
		if _, err := ctx.Coalesce().ToObject(val); err != nil {
			return fmt.Errorf("expected vector or object, got %T", val)
		}
	}

	return nil
}
