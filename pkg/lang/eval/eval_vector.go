// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/types"
)

func evalVectorNode(ctx types.Context, vec *ast.VectorNode) (types.Context, any, error) {
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
		innerCtx, data, err = evalExpression(innerCtx, &expr)
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to eval expression %s: %w", expr.String(), err)
		}

		result.Data[i] = data
	}

	return ctx, result, nil
}
