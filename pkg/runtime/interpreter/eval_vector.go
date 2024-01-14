// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package interpreter

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/pathexpr"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

func (i *interpreter) EvalVectorNode(ctx types.Context, vec ast.VectorNode) (any, error) {
	result := make([]any, len(vec.Expressions))

	var (
		data any
		err  error
	)

	for ii, expr := range vec.Expressions {
		data, err = i.EvalExpression(ctx, expr)
		if err != nil {
			return nil, fmt.Errorf("failed to eval expression %s: %w", expr.String(), err)
		}

		result[ii] = data
	}

	deeper, err := pathexpr.Traverse(ctx, result, vec.PathExpression)
	if err != nil {
		return nil, err
	}

	return deeper, nil
}
