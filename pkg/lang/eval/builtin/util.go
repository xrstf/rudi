// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval"
	"go.xrstf.de/otto/pkg/lang/eval/types"
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
