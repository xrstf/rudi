// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package ast

import (
	"fmt"
	"strings"
)

type Program struct {
	Statements []Statement
}

func (p Program) String() string {
	statements := make([]string, len(p.Statements))
	for i, s := range p.Statements {
		statements[i] = s.String()
	}

	return strings.Join(statements, "\n")
}

func (Program) NodeName() string {
	return "Program"
}

type Statement struct {
	Expression Expression
}

func (s Statement) String() string {
	return s.Expression.String()
}

func (Statement) NodeName() string {
	return "Statement"
}

type Expression struct {
	SymbolNode     *Symbol
	TupleNode      *Tuple
	VectorNode     *Vector
	ObjectNode     *Object
	NumberNode     *Number
	IdentifierNode *Identifier
	StringNode     *String
	BoolNode       *Bool
	NullNode       *Null
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
	case e.IdentifierNode != nil:
		return e.IdentifierNode.String()
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

func (e Expression) NodeName() string {
	name := ""

	switch {
	case e.SymbolNode != nil:
		name = "Symbol"
	case e.TupleNode != nil:
		name = "Tuple"
	case e.VectorNode != nil:
		name = "Vector"
	case e.ObjectNode != nil:
		name = "Object"
	case e.NumberNode != nil:
		name = "Number"
	case e.IdentifierNode != nil:
		name = "Identifier"
	case e.StringNode != nil:
		name = "String"
	case e.BoolNode != nil:
		name = "Bool"
	case e.NullNode != nil:
		name = "Null"
	default:
		name = "?"
	}

	return "Expression(" + name + ")"
}

type Symbol struct {
	PathExpression *PathExpression // can be combined with Variable
	Variable       *Variable
}

func (s Symbol) String() string {
	path := ""
	if s.PathExpression != nil {
		path = s.PathExpression.String()
	}

	switch {
	case s.Variable != nil:
		return s.Variable.String() + path
	case s.PathExpression != nil:
		// bare path expressions have a leading dot to distinguish them from array constructors
		if strings.HasPrefix(path, "[") {
			path = "." + path
		}

		return path
	default:
		return "<unknown symbol>"
	}
}

func (s Symbol) NodeName() string {
	name := ""

	switch {
	case s.Variable != nil:
		name = "Variable"
	case s.PathExpression != nil:
		name = "PathExpression"
	default:
		name = "?"
	}

	return "Symbol(" + name + ")"
}

type Tuple struct {
	Expressions []Expression
}

func (t Tuple) String() string {
	exprs := make([]string, len(t.Expressions))
	for i, expr := range t.Expressions {
		exprs[i] = expr.String()
	}

	return "(" + strings.Join(exprs, " ") + ")"
}

func (Tuple) NodeName() string {
	return "Tuple"
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

func (Vector) NodeName() string {
	return "Vector"
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

func (Object) NodeName() string {
	return "Object"
}

type KeyValuePair struct {
	Key   Expression
	Value Expression
}

func (kv KeyValuePair) String() string {
	return kv.Key.String() + " " + kv.Value.String()
}

func (KeyValuePair) NodeName() string {
	return "KeyValuePair"
}

type Variable struct {
	Name string
}

func (v Variable) String() string {
	return "$" + v.Name
}

func (Variable) NodeName() string {
	return "Variable"
}

type Identifier struct {
	Name string
}

func (i Identifier) String() string {
	return i.Name
}

func (Identifier) NodeName() string {
	return "Identifier"
}

type String struct {
	Value string
}

func (s String) String() string {
	return fmt.Sprintf("%q", s.Value)
}

func (String) NodeName() string {
	return "String"
}

type Number struct {
	Value interface{}
}

func (n Number) IsInteger() bool {
	_, ok := n.Value.(int64)
	return ok
}

func (n Number) IsFloat() bool {
	_, ok := n.Value.(float64)
	return ok
}

func (n Number) String() string {
	if n.IsFloat() {
		return fmt.Sprintf("%f", n.Value)
	}

	return fmt.Sprintf("%d", n.Value)
}

func (Number) NodeName() string {
	return "Number"
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

func (Bool) NodeName() string {
	return "Bool"
}

type Null struct{}

func (Null) String() string {
	return "null"
}

func (Null) NodeName() string {
	return "Null"
}

type PathExpression struct {
	Steps []Accessor
}

func (e *PathExpression) Prepend(step Accessor) {
	e.Steps = append([]Accessor{step}, e.Steps...)
}

func (e PathExpression) String() string {
	result := ""
	for _, step := range e.Steps {
		result += step.String()
	}

	return result
}

func (PathExpression) NodeName() string {
	return "PathExpression"
}

type Accessor struct {
	Identifier *Identifier
	StringNode *String
	Variable   *Variable
	Tuple      *Tuple
	Integer    *int64
}

func (a Accessor) String() string {
	switch {
	case a.Identifier != nil:
		return fmt.Sprintf(".%s", a.Identifier.String())
	case a.Variable != nil:
		return fmt.Sprintf(".%s", a.Variable.String())
	case a.StringNode != nil:
		return fmt.Sprintf("[%s]", a.StringNode.String())
	case a.Tuple != nil:
		return fmt.Sprintf("[%s]", a.Tuple.String())
	case a.Integer != nil:
		return fmt.Sprintf("[%d]", *a.Integer)
	default:
		return "<unknown accessor>"
	}
}

func (a Accessor) NodeName() string {
	name := ""

	switch {
	case a.Identifier != nil:
		name = "Identifier"
	case a.Variable != nil:
		name = "Variable"
	case a.StringNode != nil:
		name = "StringNode"
	case a.Tuple != nil:
		name = "Tuple"
	case a.Integer != nil:
		name = "Integer"
	default:
		name = "?"
	}

	return "Accessor(" + name + ")"
}
