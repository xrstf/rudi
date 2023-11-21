// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"io"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

func DumpSymbol(sym *ast.Symbol, out io.Writer, depth int) error {
	switch {
	case sym.Variable != nil:
		return dumpVariable(sym.Variable, sym.PathExpression, out, depth)

	case sym.PathExpression != nil:
		if err := writeString(out, "(symbol "); err != nil {
			return err
		}

		if err := dumpPathExpression(sym.PathExpression, out, depth); err != nil {
			return err
		}

		return writeString(out, ")")
	}

	return fmt.Errorf("unknown symbol %T (%s)", sym, sym.String())
}

func dumpVariable(v *ast.Variable, path *ast.PathExpression, out io.Writer, depth int) error {
	varName := string(*v)

	if err := writeString(out, fmt.Sprintf("(symbol (var %s)", varName)); err != nil {
		return err
	}

	if path != nil {
		if err := writeString(out, " "); err != nil {
			return err
		}

		if err := dumpPathExpression(path, out, depth); err != nil {
			return err
		}
	}

	return writeString(out, ")")
}

func dumpPathExpression(path *ast.PathExpression, out io.Writer, depth int) error {
	if err := writeString(out, "(path ["); err != nil {
		return err
	}

	for i, step := range path.Steps {
		err := DumpExpression(step, out, depth)
		if err != nil {
			return err
		}

		if i < len(path.Steps)-1 {
			if err := writeString(out, " "); err != nil {
				return err
			}
		}
	}

	return writeString(out, "])")
}

func dumpOptionalPathExpression(path *ast.PathExpression, out io.Writer, depth int) error {
	if path == nil {
		return nil
	}

	if err := writeString(out, "."); err != nil {
		return err
	}

	return dumpPathExpression(path, out, depth)
}
