// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"io"
	"strings"

	"go.xrstf.de/otto/pkg/lang/ast"
)

func DumpTuple(tup *ast.Tuple, out io.Writer, depth int) error {
	if depth == DoNotIndent || len(tup.Expressions) == 0 {
		return dumpTupleSingleline(tup, out, depth)
	}

	// check if we can in-line or if we need to put each element on its own line
	var buf strings.Builder
	for _, expr := range tup.Expressions {
		if err := DumpExpression(expr, &buf, 0); err != nil {
			return err
		}
	}

	if buf.Len() > 50 || depth > 10 {
		return dumpTupleMultiline(tup, out, depth)
	} else {
		return dumpTupleSingleline(tup, out, depth)
	}
}

func dumpTupleSingleline(tup *ast.Tuple, out io.Writer, depth int) error {
	if err := writeString(out, "(tuple"); err != nil {
		return err
	}

	for _, expr := range tup.Expressions {
		if err := writeString(out, " "); err != nil {
			return err
		}

		if err := DumpExpression(expr, out, depth); err != nil {
			return err
		}
	}

	if err := writeString(out, ")"); err != nil {
		return err
	}

	return dumpOptionalPathExpression(tup.PathExpression, out, depth)
}

func dumpTupleMultiline(tup *ast.Tuple, out io.Writer, depth int) error {
	prefixInner := strings.Repeat(Indent, depth+1)

	if err := writeString(out, "(tuple"); err != nil {
		return err
	}

	for _, expr := range tup.Expressions {
		if err := writeString(out, "\n"+prefixInner); err != nil {
			return err
		}

		if err := DumpExpression(expr, out, depth+1); err != nil {
			return err
		}
	}

	prefixOuter := strings.Repeat(Indent, depth)
	if err := writeString(out, "\n"+prefixOuter+")"); err != nil {
		return err
	}

	return dumpOptionalPathExpression(tup.PathExpression, out, depth)
}
