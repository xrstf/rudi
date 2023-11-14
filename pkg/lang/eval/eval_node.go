// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

func evalNode(ctx types.Context, node ast.Node) (types.Context, any, error) {
	switch asserted := node.(type) {
	case ast.Null:
		return evalNull(ctx, &asserted)
	case ast.Bool:
		return evalBool(ctx, &asserted)
	case ast.String:
		return evalString(ctx, &asserted)
	case ast.Number:
		return evalNumber(ctx, &asserted)
	case ast.ObjectNode:
		return evalObjectNode(ctx, &asserted)
	case ast.VectorNode:
		return evalVectorNode(ctx, &asserted)
	case ast.Symbol:
		return evalSymbol(ctx, &asserted)
	case ast.Tuple:
		return evalTuple(ctx, &asserted)
	case ast.Identifier:
		return evalIdentifier(ctx, &asserted)
	}

	return ctx, nil, fmt.Errorf("unknown node %T (%s)", node.NodeName(), node.String())
}
