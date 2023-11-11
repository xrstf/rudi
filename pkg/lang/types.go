package lang

import (
	"fmt"
	"strings"
)

type Program struct {
	Expressions []Expression
}

func (p Program) String() string {
	exprs := make([]string, len(p.Expressions))
	for i, e := range p.Expressions {
		exprs[i] = e.String()
	}

	return strings.Join(exprs, "\n")
}

type Expression struct {
	SymbolNode *Symbol
	TupleNode  *Tuple
	VectorNode *Vector
	ObjectNode *Object
	NumberNode *Number
	StringNode *String
	BoolNode   *Bool
	NullNode   *Null
}

func (e Expression) String() string {
	switch {
	case e.SymbolNode != nil:
		return e.SymbolNode.String()
	case e.TupleNode != nil:
		return e.TupleNode.String()
	case e.VectorNode != nil:
		return e.VectorNode.String()
	case e.ObjectNode != nil:
		return e.ObjectNode.String()
	case e.NumberNode != nil:
		return e.NumberNode.String()
	case e.StringNode != nil:
		return e.StringNode.String()
	case e.BoolNode != nil:
		return e.BoolNode.String()
	case e.NullNode != nil:
		return e.NullNode.String()
	default:
		return "<unknown expression>"
	}
}

type Symbol struct {
	JSONPath   *JSONPath
	Variable   *Variable
	Identifier *Identifier
}

func (s Symbol) String() string {
	switch {
	case s.JSONPath != nil:
		return s.JSONPath.String()
	case s.Variable != nil:
		return s.Variable.String()
	case s.Identifier != nil:
		return s.Identifier.String()
	default:
		return "<unknown symbol>"
	}
}

type Tuple struct {
	Symbol      Symbol
	Expressions []Expression
}

func (t Tuple) String() string {
	exprs := make([]string, len(t.Expressions))
	for i, expr := range t.Expressions {
		exprs[i] = expr.String()
	}

	result := "(" + t.Symbol.String()
	if len(exprs) > 0 {
		result += " " + strings.Join(exprs, " ")
	}

	return result + ")"
}

type Vector struct {
	Expressions []Expression
}

func (v Vector) String() string {
	exprs := make([]string, len(v.Expressions))
	for i, expr := range v.Expressions {
		exprs[i] = expr.String()
	}
	return "[" + strings.Join(exprs, " ") + "]"
}

type Object struct {
	Data []KeyValuePair
}

func (o Object) String() string {
	pairs := make([]string, len(o.Data))
	for i, pair := range o.Data {
		pairs[i] = pair.String()
	}
	return "{" + strings.Join(pairs, " ") + "}"
}

type KeyValuePair struct {
	Key   Symbol
	Value Expression
}

func (kv KeyValuePair) String() string {
	return kv.Key.String() + " " + kv.Value.String()
}

type Variable struct {
	Name string
}

func (v Variable) String() string {
	return "$" + v.Name
}

type Identifier struct {
	Name string
}

func (i Identifier) String() string {
	return i.Name
}

type JSONPath []string

func (p JSONPath) String() string {
	return "." + strings.Join(p, ".")
}

type String struct {
	Value string
}

func (s String) String() string {
	return fmt.Sprintf("%q", s.Value)
}

type Number struct {
	Value float64
}

func (n Number) String() string {
	return fmt.Sprintf("%f", n.Value)
}

type Bool struct {
	Value bool
}

func (b Bool) String() string {
	if b.Value {
		return "true"
	} else {
		return "false"
	}
}

type Null struct{}

func (Null) String() string {
	return "null"
}
