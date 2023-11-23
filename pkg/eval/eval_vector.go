// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

// evaluated vectors are technically considered expressions.
func EvalVector(ctx types.Context, vec ast.Vector) (types.Context, any, error) {
	return ctx, vec, nil
}

func EvalVectorNode(ctx types.Context, vec ast.VectorNode) (types.Context, any, error) {
	innerCtx := ctx
	result := ast.Vector{
		Data: make([]any, len(vec.Expressions)),
	}

	var (
		data any
		err  error
	)

	for i, expr := range vec.Expressions {
		// Keep overwriting the current context, so that e.g. variables
		// defined in one vector element can be used in all following
		// elements (no idea why you would define vars in vectors tho).
		innerCtx, data, err = EvalExpression(innerCtx, expr)
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to eval expression %s: %w", expr.String(), err)
		}

		result.Data[i] = data
	}

	if vec.PathExpression != nil {
		deeper, err := TraversePathExpression(ctx, result, vec.PathExpression)
		if err != nil {
			return ctx, nil, err
		}

		return ctx, deeper, nil
	}

	return ctx, result, nil
}
