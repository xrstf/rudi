// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package interpreter

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/pathexpr"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

func (i *interpreter) EvalVectorNode(ctx types.Context, vec ast.VectorNode) (types.Context, any, error) {
	innerCtx := ctx
	result := make([]any, len(vec.Expressions))

	var (
		data any
		err  error
	)

	for ii, expr := range vec.Expressions {
		// Keep overwriting the current context, so that e.g. variables
		// defined in one vector element can be used in all following
		// elements (no idea why you would define vars in vectors tho).
		innerCtx, data, err = i.EvalExpression(innerCtx, expr)
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to eval expression %s: %w", expr.String(), err)
		}

		result[ii] = data
	}

	deeper, err := pathexpr.Apply(ctx, result, vec.PathExpression)
	if err != nil {
		return ctx, nil, err
	}

	return ctx, deeper, nil
}
