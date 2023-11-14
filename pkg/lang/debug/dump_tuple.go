// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"io"
	"strings"

	"go.xrstf.de/otto/pkg/lang/ast"
)

func dumpTuple(tup *ast.Tuple, out io.Writer, depth int) error {
	if depth == doNotIndent {
		return dumpTupleSingleline(tup, out, depth)
	}

	// check if we can in-line or if we need to put each element on its own line
	var buf strings.Builder
	for _, expr := range tup.Expressions {
		if err := dumpNode(expr, &buf, 0); err != nil {
			return err
		}
	}

	if len(tup.Expressions) > 1 && buf.Len() > 50 || depth > 10 {
		return dumpTupleMultiline(tup, out, depth)
	} else {
		return dumpTupleSingleline(tup, out, depth)
	}
}

func dumpTupleSingleline(tup *ast.Tuple, out io.Writer, depth int) error {
	if err := writeString(out, "(tuple "); err != nil {
		return err
	}

	for i, expr := range tup.Expressions {
		if err := dumpNode(expr, out, depth); err != nil {
			return err
		}

		if i < len(tup.Expressions)-1 {
			if err := writeString(out, " "); err != nil {
				return err
			}
		}
	}

	return writeString(out, ")")
}

func dumpTupleMultiline(tup *ast.Tuple, out io.Writer, depth int) error {
	prefix := strings.Repeat(Indent, depth+1)

	if err := writeString(out, "(tuple"); err != nil {
		return err
	}

	if len(tup.Expressions) == 0 {
		return writeString(out, ")")
	}

	if err := writeString(out, " "); err != nil {
		return err
	}

	if err := dumpNode(tup.Expressions[0], out, depth+1); err != nil {
		return err
	}

	for _, expr := range tup.Expressions[1:] {
		if err := writeString(out, "\n"+prefix); err != nil {
			return err
		}

		if err := dumpNode(expr, out, depth+1); err != nil {
			return err
		}
	}

	prefix = strings.Repeat(Indent, depth)
	return writeString(out, "\n"+prefix+")")
}
