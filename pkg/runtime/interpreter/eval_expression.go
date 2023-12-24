// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package interpreter

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

func (i *interpreter) EvalExpression(ctx types.Context, expr ast.Expression) (types.Context, any, error) {
	switch asserted := expr.(type) {
	case ast.Null:
		return i.EvalNull(ctx, asserted)
	case ast.Bool:
		return i.EvalBool(ctx, asserted)
	case ast.String:
		return i.EvalString(ctx, asserted)
	case ast.Number:
		return i.EvalNumber(ctx, asserted)
	case ast.ObjectNode:
		return i.EvalObjectNode(ctx, asserted)
	case ast.VectorNode:
		return i.EvalVectorNode(ctx, asserted)
	case ast.Symbol:
		return i.EvalSymbol(ctx, asserted)
	case ast.Tuple:
		return i.EvalTuple(ctx, asserted)
	case ast.Identifier:
		return i.EvalIdentifier(ctx, asserted)
	case ast.Shim:
		return i.EvalShim(ctx, asserted)
	}

	return ctx, nil, fmt.Errorf("unknown expression %s (%T)", expr.String(), expr)
}
