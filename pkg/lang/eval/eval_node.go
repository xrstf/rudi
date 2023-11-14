// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

func EvalExpression(ctx types.Context, node ast.Expression) (types.Context, any, error) {
	switch asserted := node.(type) {
	case ast.Null:
		return EvalNull(ctx, asserted)
	case ast.Bool:
		return EvalBool(ctx, asserted)
	case ast.String:
		return EvalString(ctx, asserted)
	case ast.Number:
		return EvalNumber(ctx, asserted)
	case ast.ObjectNode:
		return EvalObjectNode(ctx, asserted)
	case ast.VectorNode:
		return EvalVectorNode(ctx, asserted)
	case ast.Symbol:
		return EvalSymbol(ctx, asserted)
	case ast.Tuple:
		return EvalTuple(ctx, asserted)
	case ast.Identifier:
		return EvalIdentifier(ctx, asserted)
	}

	return ctx, nil, fmt.Errorf("unknown expression %s (%s)", node.ExpressionName(), node.String())
}
