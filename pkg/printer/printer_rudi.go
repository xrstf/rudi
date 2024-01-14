// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package printer

import (
	"fmt"
	"io"
	"strconv"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

type rudiPrinter struct {
	basePrinter
}

func NewRudiPrinter(out io.Writer) Printer {
	return &rudiPrinter{
		basePrinter: newBasePrinter(out),
	}
}

func (p *rudiPrinter) Print(v any) error {
	return printAny(v, p)
}

func (p *rudiPrinter) Null() error {
	return p.write("null")
}

func (p *rudiPrinter) Bool(b bool) error {
	return p.write(strconv.FormatBool(b))
}

func (p *rudiPrinter) Number(value any) error {
	return p.write(fmt.Sprintf("%v", value))
}

func (p *rudiPrinter) String(str string) error {
	return p.write(fmt.Sprintf("%q", str))
}

func (p *rudiPrinter) Identifier(ident *ast.Identifier) error {
	name := ident.Name
	if ident.Bang {
		name += "!"
	}

	return p.write(name)
}

func (p *rudiPrinter) Vector(vec []any) error {
	return p.printVector(vec, nil)
}

func (p *rudiPrinter) VectorNode(vec *ast.VectorNode) error {
	data := make([]any, len(vec.Expressions))
	for i := range vec.Expressions {
		data[i] = vec.Expressions[i]
	}

	return p.printVector(data, vec.PathExpression)
}

func (p *rudiPrinter) printVector(vec []any, pathExpr *ast.PathExpression) error {
	if err := p.write("["); err != nil {
		return err
	}

	for i, val := range vec {
		if err := printAny(val, p); err != nil {
			return err
		}

		if i < len(vec)-1 {
			if err := p.write(" "); err != nil {
				return err
			}
		}
	}

	if err := p.write("]"); err != nil {
		return err
	}

	return p.printPathExpression(pathExpr)
}

func (p *rudiPrinter) Object(obj map[string]any) error {
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

func (p *rudiPrinter) ObjectNode(obj *ast.ObjectNode) error {
	out := make(Object, len(obj.Data))

	for i, pair := range obj.Data {
		out[i] = KeyValuePair{
			Key:   pair.Key,
			Value: pair.Value,
		}
	}

	return p.printObject(out, obj.PathExpression)
}

func (p *rudiPrinter) printObject(obj Object, pathExpr *ast.PathExpression) error {
	if err := p.write("{"); err != nil {
		return err
	}

	for i, pair := range obj {
		// turn basic string keys into identifiers (i.e. {"foo" "bar"} into {foo "bar"})
		key := pair.Key
		if str, ok := key.(ast.String); ok && ast.IdentifierNamePattern.MatchString(string(str)) {
			key = ast.Identifier{Name: string(str)}
		}

		if err := printAny(key, p); err != nil {
			return err
		}

		if err := p.write(" "); err != nil {
			return err
		}

		if err := printAny(pair.Value, p); err != nil {
			return err
		}

		if i < len(obj)-1 {
			if err := p.write(" "); err != nil {
				return err
			}
		}
	}

	if err := p.write("}"); err != nil {
		return err
	}

	return p.printPathExpression(pathExpr)
}

func (p *rudiPrinter) Tuple(tup *ast.Tuple) error {
	if err := p.write("("); err != nil {
		return err
	}

	for i, expr := range tup.Expressions {
		if err := printAny(expr, p); err != nil {
			return err
		}

		if i < len(tup.Expressions)-1 {
			if err := p.write(" "); err != nil {
				return err
			}
		}
	}

	if err := p.write(")"); err != nil {
		return err
	}

	return p.printPathExpression(tup.PathExpression)
}

func (p *rudiPrinter) Symbol(sym *ast.Symbol) error {
	switch {
	case sym.Variable != nil:
		return p.variable(sym.Variable, sym.PathExpression)

	case sym.PathExpression != nil:
		steps := sym.PathExpression.Steps
		if len(steps) == 0 {
			return fmt.Errorf("invalid symbol: path expression is empty")
		}

		// Path expressions on the global document that start with a vector step
		// must have an extra leading dot, to distinguish .[1] from [1] (which
		// would be "a vector with one element: 1" and not a path expression).
		if p.renderAsVectorStep(steps[0]) {
			if err := p.write("."); err != nil {
				return err
			}
		}

		return p.printPathExpression(sym.PathExpression)
	}

	return fmt.Errorf("unknown symbol %T (%s)", sym, sym.String())
}

func (p *rudiPrinter) variable(v *ast.Variable, path *ast.PathExpression) error {
	if err := p.write(fmt.Sprintf("$%s", string(*v))); err != nil {
		return err
	}

	if err := p.printPathExpression(path); err != nil {
		return err
	}

	return nil
}

func (p *rudiPrinter) Expression(expr ast.Expression) error {
	return printAny(expr, p)
}

func (p *rudiPrinter) Statement(stmt *ast.Statement) error {
	return printAny(stmt.Expression, p)
}

func (p *rudiPrinter) Program(prog *ast.Program) error {
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

func (p *rudiPrinter) printPathExpression(path *ast.PathExpression) error {
	if path == nil {
		return nil
	}

	for _, step := range path.Steps {
		if p.renderAsVectorStep(step) {
			if err := p.write("["); err != nil {
				return err
			}

			if step.Filter != nil {
				if err := p.write("?"); err != nil {
					return err
				}

				if err := printAny(step.Filter, p); err != nil {
					return err
				}
			} else {
				if err := printAny(step.Expression, p); err != nil {
					return err
				}
			}

			if err := p.write("]"); err != nil {
				return err
			}
		} else {
			if err := p.write("."); err != nil {
				return err
			}

			switch asserted := step.Expression.(type) {
			case ast.String:
				if err := p.write(string(asserted)); err != nil {
					return err
				}
			default:
				panic("Should not reach this point: renderAsVectorStep is out of sync.")
			}
		}
	}

	return nil
}

func (p *rudiPrinter) renderAsVectorStep(step ast.PathStep) bool {
	if step.Filter != nil {
		return true
	}

	switch asserted := step.Expression.(type) {
	case ast.String:
		return !ast.PathIdentifierPattern.MatchString(string(asserted))
	default:
		return true
	}
}
