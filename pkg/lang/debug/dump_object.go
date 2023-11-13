// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"io"
	"strings"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func dumpObject(obj *ast.ObjectNode, out io.Writer, depth int) error {
	// check if we can in-line or if we need to put each element on its own line
	var buf strings.Builder
	for _, pair := range obj.Data {
		if err := dumpExpression(&pair.Key, &buf, 0); err != nil {
			return err
		}

		if err := dumpExpression(&pair.Value, &buf, 0); err != nil {
			return err
		}
	}

	if buf.Len() > 50 || depth > 10 {
		return dumpObjectMultiline(obj, out, depth)
	} else {
		return dumpObjectSingleline(obj, out, depth)
	}
}

func dumpObjectSingleline(obj *ast.ObjectNode, out io.Writer, depth int) error {
	if err := writeString(out, "{"); err != nil {
		return err
	}

	for i, pair := range obj.Data {
		if err := dumpExpression(&pair.Key, out, depth); err != nil {
			return err
		}

		if err := writeString(out, " "); err != nil {
			return err
		}

		if err := dumpExpression(&pair.Value, out, depth); err != nil {
			return err
		}

		if i < len(obj.Data)-1 {
			if err := writeString(out, " "); err != nil {
				return err
			}
		}
	}

	return writeString(out, "}")
}

func dumpObjectMultiline(obj *ast.ObjectNode, out io.Writer, depth int) error {
	prefix := strings.Repeat(Indent, depth+1)

	if err := writeString(out, "{"); err != nil {
		return err
	}

	for _, pair := range obj.Data {
		if err := writeString(out, "\n"+prefix); err != nil {
			return err
		}

		if err := dumpExpression(&pair.Key, out, depth+1); err != nil {
			return err
		}

		if err := writeString(out, " "); err != nil {
			return err
		}

		if err := dumpExpression(&pair.Value, out, depth+1); err != nil {
			return err
		}
	}

	prefix = strings.Repeat(Indent, depth)
	return writeString(out, "\n"+prefix+"}")
}
