// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package printer

import (
	"fmt"
	"io"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

type astPrinter struct {
	basePrinter
}

func NewAstPrinter(out io.Writer) Printer {
	return &astPrinter{
		basePrinter: newBasePrinter(out),
	}
}

func (p *astPrinter) Print(v any) error {
	return printAny(v, p)
}

func (p *astPrinter) Null() error {
	return p.write("(null)")
}

func (p *astPrinter) Bool(b bool) error {
	return p.write(fmt.Sprintf("(bool %v)", b))
}

func (p *astPrinter) Number(value any) error {
	return p.write(fmt.Sprintf("(number (%T %v))", value, value))
}

func (p *astPrinter) String(str string) error {
	return p.write(fmt.Sprintf("(string %q)", str))
}

func (p *astPrinter) Identifier(ident *ast.Identifier) error {
	var bang string
	if ident.Bang {
		bang = " (bang)"
	}

	return p.write(fmt.Sprintf("(identifier %s%s)", ident.Name, bang))
}

func (p *astPrinter) Vector(vec []any) error {
	return p.printVector(vec, nil)
}

func (p *astPrinter) VectorNode(vec *ast.VectorNode) error {
	data := make([]any, len(vec.Expressions))
	for i := range vec.Expressions {
		data[i] = vec.Expressions[i]
	}

	return p.printVector(data, vec.PathExpression)
}

func (p *astPrinter) printVector(vec []any, pathExpr *ast.PathExpression) error {
	if err := p.write("(vector"); err != nil {
		return err
	}

	for _, val := range vec {
		if err := p.write(" "); err != nil {
			return err
		}

		if err := printAny(val, p); err != nil {
			return err
		}
	}

	if err := p.write(")"); err != nil {
		return err
	}

	return p.writeOptionalPathExpression(pathExpr)
}

func (p *astPrinter) Object(obj map[string]any) error {
	out := make(Object, len(obj))

	i := 0
	for k, v := range obj {
		out[i] = KeyValuePair{
			Key:   k,
			Value: v,
		}
		i++
	}

	return p.printObject(out, nil)
}

func (p *astPrinter) ObjectNode(obj *ast.ObjectNode) error {
	out := make(Object, len(obj.Data))

	for i, pair := range obj.Data {
		out[i] = KeyValuePair{
			Key:   pair.Key,
			Value: pair.Value,
		}
	}

	return p.printObject(out, obj.PathExpression)
}

func (p *astPrinter) printObject(obj Object, pathExpr *ast.PathExpression) error {
	if err := p.write("(object"); err != nil {
		return err
	}

	for _, pair := range obj {
		if err := p.write(" ("); err != nil {
			return err
		}

		if err := printAny(pair.Key, p); err != nil {
			return err
		}

		if err := p.write(" "); err != nil {
			return err
		}

		if err := printAny(pair.Value, p); err != nil {
			return err
		}

		if err := p.write(")"); err != nil {
			return err
		}
	}

	if err := p.write(")"); err != nil {
		return err
	}

	return p.writeOptionalPathExpression(pathExpr)
}

func (p *astPrinter) Tuple(tup *ast.Tuple) error {
	if err := p.write("(tuple"); err != nil {
		return err
	}

	for _, expr := range tup.Expressions {
		if err := p.write(" "); err != nil {
			return err
		}

		if err := printAny(expr, p); err != nil {
			return err
		}
	}

	if err := p.write(")"); err != nil {
		return err
	}

	return p.writeOptionalPathExpression(tup.PathExpression)
}

func (p *astPrinter) Symbol(sym *ast.Symbol) error {
	switch {
	case sym.Variable != nil:
		return p.variable(sym.Variable, sym.PathExpression)

	case sym.PathExpression != nil:
		if err := p.write("(symbol "); err != nil {
			return err
		}

		if err := p.writePathExpression(sym.PathExpression); err != nil {
			return err
		}

		return p.write(")")
	}

	return fmt.Errorf("unknown symbol %T (%s)", sym, sym.String())
}

func (p *astPrinter) variable(v *ast.Variable, path *ast.PathExpression) error {
	varName := string(*v)

	if err := p.write(fmt.Sprintf("(symbol (var %s)", varName)); err != nil {
		return err
	}

	if path != nil {
		if err := p.write(" "); err != nil {
			return err
		}

		if err := p.writePathExpression(path); err != nil {
			return err
		}
	}

	return p.write(")")
}

func (p *astPrinter) Expression(expr ast.Expression) error {
	return printAny(expr, p)
}

func (p *astPrinter) Statement(stmt *ast.Statement) error {
	return printAny(stmt.Expression, p)
}

func (p *astPrinter) Program(prog *ast.Program) error {
	for i, stmt := range prog.Statements {
		if err := p.Statement(&stmt); err != nil {
			return err
		}

		if i < len(prog.Statements)-1 {
			if err := p.write(" "); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *astPrinter) writePathExpression(path *ast.PathExpression) error {
	if err := p.write("(path ["); err != nil {
		return err
	}

	for i, step := range path.Steps {
		if step.Filter != nil {
			if err := p.write("(filter "); err != nil {
				return err
			}

			if err := printAny(step.Filter, p); err != nil {
				return err
			}

			if err := p.write(")"); err != nil {
				return err
			}
		} else {
			if err := p.write("(expr "); err != nil {
				return err
			}

			if err := printAny(step.Expression, p); err != nil {
				return err
			}

			if err := p.write(")"); err != nil {
				return err
			}
		}

		if i < len(path.Steps)-1 {
			if err := p.write(" "); err != nil {
				return err
			}
		}
	}

	return p.write("])")
}

func (p *astPrinter) writeOptionalPathExpression(path *ast.PathExpression) error {
	if path == nil {
		return nil
	}

	if err := p.write("."); err != nil {
		return err
	}

	return p.writePathExpression(path)
}
