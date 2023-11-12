// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalExpression(ctx Context, expr *ast.Expression) (Context, interface{}, error) {
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
		return evalObject(ctx, expr.ObjectNode)
	case expr.VectorNode != nil:
		return evalVector(ctx, expr.VectorNode)
	case expr.SymbolNode != nil:
		return evalSymbol(ctx, expr.SymbolNode)
	case expr.TupleNode != nil:
		return evalTuple(ctx, expr.TupleNode)
	}

	return ctx, nil, fmt.Errorf("unknown expression %T (%s)", expr, expr.String())
}
