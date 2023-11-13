// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"errors"
	"fmt"
	"io"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func dumpSymbol(sym *ast.Symbol, out io.Writer, depth int) error {
	switch {
	case sym.Variable != nil:
		return dumpVariable(sym.Variable, sym.PathExpression, out, depth)

	case sym.PathExpression != nil:
		return dumpPathExpression(sym.PathExpression, out, depth)
	}

	return fmt.Errorf("unknown symbol %T (%s)", sym, sym.String())
}

func dumpVariable(v *ast.Variable, path *ast.PathExpression, out io.Writer, depth int) error {
	varName := string(*v)

	if err := writeString(out, fmt.Sprintf("(var %s", varName)); err != nil {
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
		var err error

		switch {
		case step.Identifier != nil:
			err = dumpIdentifier(step.Identifier, out)
		case step.StringNode != nil:
			err = dumpString(step.StringNode, out)
		case step.Variable != nil:
			err = dumpVariable(step.Variable, nil, out, depth)
		case step.Tuple != nil:
			err = dumpTupleSingleline(step.Tuple, out, depth)
		case step.Integer != nil:
			err = dumpInteger(step.Integer, out)
		default:
			err = errors.New("unexpected element in path expression")
		}

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
