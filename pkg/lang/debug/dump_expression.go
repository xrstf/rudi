// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"io"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func dumpExpression(expr *ast.Expression, out io.Writer, depth int) error {
	switch {
	case expr.NullNode != nil:
		return dumpNull(expr.NullNode, out)
	case expr.BoolNode != nil:
		return dumpBool(expr.BoolNode, out)
	case expr.StringNode != nil:
		return dumpString(expr.StringNode, out)
	case expr.NumberNode != nil:
		return dumpNumber(expr.NumberNode, out)
	case expr.ObjectNode != nil:
		return dumpObject(expr.ObjectNode, out, depth)
	case expr.VectorNode != nil:
		return dumpVector(expr.VectorNode, out, depth)
	case expr.SymbolNode != nil:
		return dumpSymbol(expr.SymbolNode, out, depth)
	case expr.TupleNode != nil:
		return dumpTuple(expr.TupleNode, out, depth)
	case expr.IdentifierNode != nil:
		return dumpIdentifier(expr.IdentifierNode, out)
	}

	return fmt.Errorf("unknown expression %s (%s)", expr.NodeName(), expr.String())
}
