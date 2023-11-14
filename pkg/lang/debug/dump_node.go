// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"io"

	"go.xrstf.de/otto/pkg/lang/ast"
)

func dumpNode(node ast.Expression, out io.Writer, depth int) error {
	switch asserted := node.(type) {
	case ast.Null:
		return dumpNull(&asserted, out)
	case ast.Bool:
		return dumpBool(&asserted, out)
	case ast.String:
		return dumpString(&asserted, out)
	case ast.Number:
		return dumpNumber(&asserted, out)
	case ast.ObjectNode:
		return dumpObject(&asserted, out, depth)
	case ast.VectorNode:
		return dumpVector(&asserted, out, depth)
	case ast.Symbol:
		return dumpSymbol(&asserted, out, depth)
	case ast.Tuple:
		return dumpTuple(&asserted, out, depth)
	case ast.Identifier:
		return dumpIdentifier(&asserted, out)
	}

	return fmt.Errorf("unknown node %s (%s)", node.ExpressionName(), node.String())
}
