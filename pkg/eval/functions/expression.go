// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package functions

import (
	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

// cachedExpression are used to avoid having to re-compute the same expressions when attempting
// to match against multiple function signatures. Since the functions themselves use coalescing to
// cast data and we allow `func(expr ast.Expression)`, we must have a way to only execute
// expressions once and re-use the result for other matchings.
type cachedExpression struct {
	result    any
	evaluated bool
	expr      ast.Expression
}

func convertArgs(args []ast.Expression) []cachedExpression {
	result := make([]cachedExpression, len(args))
	for i := range args {
		result[i] = cachedExpression{expr: args[i]}
	}
	return result
}

func (e *cachedExpression) Eval(ctx types.Context) (any, error) {
	if !e.evaluated {
		_, result, err := eval.EvalExpression(ctx, e.expr)
		if err != nil {
			return nil, err
		}

		e.result = result
		e.evaluated = true
	}

	return e.result, nil
}
