// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

func evalExpression(ctx types.Context, expr *ast.Expression) (types.Context, any, error) {
	switch {
	case expr.NullNode != nil:
		return evalNull(ctx, expr.NullNode)
	case expr.BoolNode != nil:
		return evalBool(ctx, expr.BoolNode)
	case expr.StringNode != nil:
		return evalString(ctx, expr.StringNode)
	case expr.NumberNode != nil:
		return evalNumber(ctx, expr.NumberNode)
	case expr.ObjectNode != nil:
		return evalObjectNode(ctx, expr.ObjectNode)
	case expr.VectorNode != nil:
		return evalVectorNode(ctx, expr.VectorNode)
	case expr.SymbolNode != nil:
		return evalSymbol(ctx, expr.SymbolNode)
	case expr.TupleNode != nil:
		return evalTuple(ctx, expr.TupleNode)
	case expr.IdentifierNode != nil:
		return evalIdentifier(ctx, expr.IdentifierNode)
	}

	return ctx, nil, fmt.Errorf("unknown expression %T (%s)", expr.NodeName(), expr.String())
}
