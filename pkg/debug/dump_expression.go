// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"io"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

func DumpExpression(expr ast.Expression, out io.Writer, depth int) error {
	switch asserted := expr.(type) {
	case ast.Null:
		return DumpNull(out)
	case ast.Bool:
		return DumpBool(bool(asserted), out)
	case ast.String:
		return DumpString(string(asserted), out)
	case ast.Number:
		return DumpNumber(asserted.Value, out)
	case ast.ObjectNode:
		return DumpObjectNode(asserted, out, depth)
	case ast.VectorNode:
		return DumpVectorNode(asserted, out, depth)
	case ast.Symbol:
		return DumpSymbol(&asserted, out, depth)
	case ast.Tuple:
		return DumpTuple(&asserted, out, depth)
	case ast.Identifier:
		return DumpIdentifier(&asserted, out)
	}

	return fmt.Errorf("unknown expression %s (%s)", expr.ExpressionName(), expr.String())
}
