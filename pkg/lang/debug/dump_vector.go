// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"io"
	"strings"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func dumpVector(vec *ast.VectorNode, out io.Writer, depth int) error {
	// check if we can in-line or if we need to put each element on its own line
	var buf strings.Builder
	for _, expr := range vec.Expressions {
		if err := dumpExpression(&expr, &buf, 0); err != nil {
			return err
		}
	}

	if buf.Len() > 50 || depth > 10 {
		return dumpVectorMultiline(vec, out, depth)
	} else {
		return dumpVectorSingleline(vec, out, depth)
	}
}

func dumpVectorSingleline(vec *ast.VectorNode, out io.Writer, depth int) error {
	if err := writeString(out, "["); err != nil {
		return err
	}

	for i, expr := range vec.Expressions {
		if err := dumpExpression(&expr, out, depth); err != nil {
			return err
		}

		if i < len(vec.Expressions)-1 {
			if err := writeString(out, " "); err != nil {
				return err
			}
		}
	}

	return writeString(out, "]")
}

func dumpVectorMultiline(vec *ast.VectorNode, out io.Writer, depth int) error {
	prefix := strings.Repeat(Indent, depth+1)

	if err := writeString(out, "["); err != nil {
		return err
	}

	for _, expr := range vec.Expressions {
		if err := writeString(out, "\n"+prefix); err != nil {
			return err
		}

		if err := dumpExpression(&expr, out, depth+1); err != nil {
			return err
		}
	}

	prefix = strings.Repeat(Indent, depth)
	return writeString(out, "\n"+prefix+"]")
}
