// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package printer

import (
	"fmt"
	"io"
	"strings"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

type AST struct{}

func (r AST) WriteSingleline(v any, out io.Writer) error {
	return writeAny(v, r, out, DoNotIndent)
}

func (r AST) WriteMultiline(v any, out io.Writer) error {
	return writeAny(v, r, out, 0)
}

func (AST) WriteBool(b bool, out io.Writer) error {
	return writeString(out, fmt.Sprintf("(bool %v)", b))
}

func (AST) WriteIdentifier(ident *ast.Identifier, out io.Writer) error {
	return writeString(out, fmt.Sprintf("(identifier %s)", *ident))
}

func (AST) WriteNull(out io.Writer) error {
	return writeString(out, "(null)")
}

func (AST) WriteNumber(value any, out io.Writer) error {
	return writeString(out, fmt.Sprintf("(number (%T %v))", value, value))
}

func (AST) WriteString(str string, out io.Writer) error {
	return writeString(out, fmt.Sprintf("(string %q)", str))
}

func (r AST) WriteVectorSingleline(vec []any, pathExpr *ast.PathExpression, out io.Writer, depth int) error {
	if err := writeString(out, "(vector"); err != nil {
		return err
	}

	for _, val := range vec {
		if err := writeString(out, " "); err != nil {
			return err
		}

		if err := writeAny(val, r, out, depth); err != nil {
			return err
		}
	}

	if err := writeString(out, ")"); err != nil {
		return err
	}

	return r.writeOptionalPathExpression(pathExpr, out, depth)
}

func (r AST) WriteVectorMultiline(vec []any, pathExpr *ast.PathExpression, out io.Writer, depth int) error {
	prefixInner := strings.Repeat(Indent, depth+1)

	if err := writeString(out, "(vector"); err != nil {
		return err
	}

	for _, val := range vec {
		if err := writeString(out, "\n"+prefixInner); err != nil {
			return err
		}

		if err := writeAny(val, r, out, depth+1); err != nil {
			return err
		}
	}

	prefixOuter := strings.Repeat(Indent, depth)

	if err := writeString(out, "\n"+prefixOuter+")"); err != nil {
		return err
	}

	return r.writeOptionalPathExpression(pathExpr, out, depth)
}

func (r AST) WriteObjectSingleline(obj Object, pathExpr *ast.PathExpression, out io.Writer, depth int) error {
	if err := writeString(out, "(object"); err != nil {
		return err
	}

	for _, pair := range obj {
		if err := writeString(out, " "); err != nil {
			return err
		}

		if err := writeAny(pair.Key, r, out, depth); err != nil {
			return err
		}

		if err := writeString(out, " "); err != nil {
			return err
		}

		if err := writeAny(pair.Value, r, out, depth); err != nil {
			return err
		}
	}

	if err := writeString(out, ")"); err != nil {
		return err
	}

	return r.writeOptionalPathExpression(pathExpr, out, depth)
}

func (r AST) WriteObjectMultiline(obj Object, pathExpr *ast.PathExpression, out io.Writer, depth int) error {
	prefix := strings.Repeat(Indent, depth+1)

	if err := writeString(out, "(object"); err != nil {
		return err
	}

	for _, pair := range obj {
		if err := writeString(out, "\n"+prefix); err != nil {
			return err
		}

		if err := writeAny(pair.Key, r, out, depth+1); err != nil {
			return err
		}

		if err := writeString(out, " "); err != nil {
			return err
		}

		if err := writeAny(pair.Value, r, out, depth+1); err != nil {
			return err
		}
	}

	prefix = strings.Repeat(Indent, depth)
	return writeString(out, "\n"+prefix+")")
}

func (r AST) WriteTupleSingleline(tup *ast.Tuple, out io.Writer, depth int) error {
	if err := writeString(out, "(tuple"); err != nil {
		return err
	}

	for _, expr := range tup.Expressions {
		if err := writeString(out, " "); err != nil {
			return err
		}

		if err := writeAny(expr, r, out, depth); err != nil {
			return err
		}
	}

	if err := writeString(out, ")"); err != nil {
		return err
	}

	return r.writeOptionalPathExpression(tup.PathExpression, out, depth)
}

func (r AST) WriteTupleMultiline(tup *ast.Tuple, out io.Writer, depth int) error {
	prefixInner := strings.Repeat(Indent, depth+1)

	if err := writeString(out, "(tuple"); err != nil {
		return err
	}

	for _, expr := range tup.Expressions {
		if err := writeString(out, "\n"+prefixInner); err != nil {
			return err
		}

		if err := writeAny(expr, r, out, depth+1); err != nil {
			return err
		}
	}

	prefixOuter := strings.Repeat(Indent, depth)
	if err := writeString(out, "\n"+prefixOuter+")"); err != nil {
		return err
	}

	return r.writeOptionalPathExpression(tup.PathExpression, out, depth)
}

func (r AST) WriteSymbol(sym *ast.Symbol, out io.Writer, depth int) error {
	switch {
	case sym.Variable != nil:
		return r.writeVariable(sym.Variable, sym.PathExpression, out, depth)

	case sym.PathExpression != nil:
		if err := writeString(out, "(symbol "); err != nil {
			return err
		}

		if err := r.writePathExpression(sym.PathExpression, out, depth); err != nil {
			return err
		}

		return writeString(out, ")")
	}

	return fmt.Errorf("unknown symbol %T (%s)", sym, sym.String())
}

func (r AST) writeVariable(v *ast.Variable, path *ast.PathExpression, out io.Writer, depth int) error {
	varName := string(*v)

	if err := writeString(out, fmt.Sprintf("(symbol (var %s)", varName)); err != nil {
		return err
	}

	if path != nil {
		if err := writeString(out, " "); err != nil {
			return err
		}

		if err := r.writePathExpression(path, out, depth); err != nil {
			return err
		}
	}

	return writeString(out, ")")
}

func (r AST) writePathExpression(path *ast.PathExpression, out io.Writer, depth int) error {
	if err := writeString(out, "(path ["); err != nil {
		return err
	}

	for i, step := range path.Steps {
		err := writeAny(step, r, out, depth)
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

func (r AST) writeOptionalPathExpression(path *ast.PathExpression, out io.Writer, depth int) error {
	if path == nil {
		return nil
	}

	if err := writeString(out, "."); err != nil {
		return err
	}

	return r.writePathExpression(path, out, depth)
}
