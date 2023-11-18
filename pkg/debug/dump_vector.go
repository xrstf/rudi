// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"io"
	"strings"

	"go.xrstf.de/otto/pkg/lang/ast"
)

func DumpVector(vec *ast.Vector, out io.Writer, depth int) error {
	if depth == DoNotIndent || len(vec.Data) == 0 {
		return dumpVectorSingleline(vec, out, depth)
	}

	// check if we can in-line or if we need to put each element on its own line
	var buf strings.Builder
	for _, val := range vec.Data {
		if err := dumpAny(val, &buf, 0); err != nil {
			return err
		}
	}

	if buf.Len() > 50 || depth > 10 {
		return dumpVectorMultiline(vec, out, depth)
	} else {
		return dumpVectorSingleline(vec, out, depth)
	}
}

func dumpVectorSingleline(vec *ast.Vector, out io.Writer, depth int) error {
	if err := writeString(out, "(vector"); err != nil {
		return err
	}

	for _, val := range vec.Data {
		if err := writeString(out, " "); err != nil {
			return err
		}

		if err := dumpAny(val, out, depth); err != nil {
			return err
		}
	}

	return writeString(out, ")")
}

func dumpVectorMultiline(vec *ast.Vector, out io.Writer, depth int) error {
	prefix := strings.Repeat(Indent, depth+1)

	if err := writeString(out, "(vector"); err != nil {
		return err
	}

	for _, val := range vec.Data {
		if err := writeString(out, "\n"+prefix); err != nil {
			return err
		}

		if err := dumpAny(val, out, depth+1); err != nil {
			return err
		}
	}

	prefix = strings.Repeat(Indent, depth)
	return writeString(out, "\n"+prefix+")")
}

func DumpVectorNode(vec *ast.VectorNode, out io.Writer, depth int) error {
	if depth == DoNotIndent || len(vec.Expressions) == 0 {
		return dumpVectorNodeSingleline(vec, out, depth)
	}

	// check if we can in-line or if we need to put each element on its own line
	var buf strings.Builder
	for _, expr := range vec.Expressions {
		if err := DumpExpression(expr, &buf, 0); err != nil {
			return err
		}
	}

	if buf.Len() > 50 || depth > 10 {
		return dumpVectorNodeMultiline(vec, out, depth)
	} else {
		return dumpVectorNodeSingleline(vec, out, depth)
	}
}

func dumpVectorNodeSingleline(vec *ast.VectorNode, out io.Writer, depth int) error {
	if err := writeString(out, "(vector"); err != nil {
		return err
	}

	for _, expr := range vec.Expressions {
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

	return dumpOptionalPathExpression(vec.PathExpression, out, depth)
}

func dumpVectorNodeMultiline(vec *ast.VectorNode, out io.Writer, depth int) error {
	prefixInner := strings.Repeat(Indent, depth+1)

	if err := writeString(out, "(vector"); err != nil {
		return err
	}

	for _, expr := range vec.Expressions {
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

	return dumpOptionalPathExpression(vec.PathExpression, out, depth)
}
