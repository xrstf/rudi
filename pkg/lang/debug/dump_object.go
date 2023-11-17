// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"io"
	"strings"

	"go.xrstf.de/otto/pkg/lang/ast"
)

func DumpObject(obj *ast.Object, out io.Writer, depth int) error {
	if depth == DoNotIndent || len(obj.Data) == 0 {
		return dumpObjectSingleline(obj, out, depth)
	}

	// check if we can in-line or if we need to put each element on its own line
	var buf strings.Builder
	for key, value := range obj.Data {
		if err := dumpAny(key, &buf, 0); err != nil {
			return err
		}

		if err := dumpAny(value, &buf, 0); err != nil {
			return err
		}
	}

	if buf.Len() > 50 || depth > 10 {
		return dumpObjectMultiline(obj, out, depth)
	} else {
		return dumpObjectSingleline(obj, out, depth)
	}
}

func dumpObjectSingleline(obj *ast.Object, out io.Writer, depth int) error {
	if err := writeString(out, "(object"); err != nil {
		return err
	}

	for key, value := range obj.Data {
		if err := writeString(out, " "); err != nil {
			return err
		}

		if err := dumpAny(key, out, depth); err != nil {
			return err
		}

		if err := writeString(out, " "); err != nil {
			return err
		}

		if err := dumpAny(value, out, depth); err != nil {
			return err
		}
	}

	return writeString(out, ")")
}

func dumpObjectMultiline(obj *ast.Object, out io.Writer, depth int) error {
	prefix := strings.Repeat(Indent, depth+1)

	if err := writeString(out, "(object"); err != nil {
		return err
	}

	for key, value := range obj.Data {
		if err := writeString(out, "\n"+prefix); err != nil {
			return err
		}

		if err := dumpAny(key, out, depth+1); err != nil {
			return err
		}

		if err := writeString(out, " "); err != nil {
			return err
		}

		if err := dumpAny(value, out, depth+1); err != nil {
			return err
		}
	}

	prefix = strings.Repeat(Indent, depth)
	return writeString(out, "\n"+prefix+")")
}

func DumpObjectNode(obj *ast.ObjectNode, out io.Writer, depth int) error {
	if depth == DoNotIndent || len(obj.Data) == 0 {
		return dumpObjectNodeSingleline(obj, out, depth)
	}

	// check if we can in-line or if we need to put each element on its own line
	var buf strings.Builder
	for _, pair := range obj.Data {
		if err := DumpExpression(pair.Key, &buf, 0); err != nil {
			return err
		}

		if err := DumpExpression(pair.Value, &buf, 0); err != nil {
			return err
		}
	}

	if buf.Len() > 50 || depth > 10 {
		return dumpObjectNodeMultiline(obj, out, depth)
	} else {
		return dumpObjectNodeSingleline(obj, out, depth)
	}
}

func dumpObjectNodeSingleline(obj *ast.ObjectNode, out io.Writer, depth int) error {
	if err := writeString(out, "(object"); err != nil {
		return err
	}

	for _, pair := range obj.Data {
		if err := writeString(out, " "); err != nil {
			return err
		}

		if err := DumpExpression(pair.Key, out, depth); err != nil {
			return err
		}

		if err := writeString(out, " "); err != nil {
			return err
		}

		if err := DumpExpression(pair.Value, out, depth); err != nil {
			return err
		}
	}

	if err := writeString(out, ")"); err != nil {
		return err
	}

	return dumpOptionalPathExpression(obj.PathExpression, out, depth)
}

func dumpObjectNodeMultiline(obj *ast.ObjectNode, out io.Writer, depth int) error {
	prefixInner := strings.Repeat(Indent, depth+1)

	if err := writeString(out, "(object"); err != nil {
		return err
	}

	for _, pair := range obj.Data {
		if err := writeString(out, "\n"+prefixInner); err != nil {
			return err
		}

		if err := DumpExpression(pair.Key, out, depth+1); err != nil {
			return err
		}

		if err := writeString(out, " "); err != nil {
			return err
		}

		if err := DumpExpression(pair.Value, out, depth+1); err != nil {
			return err
		}
	}

	prefixOuter := strings.Repeat(Indent, depth)

	if err := writeString(out, "\n"+prefixOuter+")"); err != nil {
		return err
	}

	return dumpOptionalPathExpression(obj.PathExpression, out, depth)
}
