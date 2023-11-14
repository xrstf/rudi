// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package ast

import (
	"fmt"
	"strings"
)

type Node interface {
	String() string
	NodeName() string
}

type Program struct {
	Statements []Statement
}

var _ Node = Program{}

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
	Tuple Tuple
}

var _ Node = Statement{}

func (s Statement) String() string {
	return s.Tuple.String()
}

func (Statement) NodeName() string {
	return "Statement"
}

type Symbol struct {
	PathExpression *PathExpression // can be combined with Variable
	Variable       *Variable
}

var _ Node = Symbol{}

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
	Expressions []Node
}

var _ Node = Tuple{}

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

// Vector is an evaluated vector.
type Vector struct {
	Data []any
}

func (Vector) NodeName() string {
	return "Vector"
}

func (v Vector) LiteralValue() any {
	return v.Data
}

// VectorNode represents the parsed code for constructing an vector.
// When an VectorNode is evaluated, it turns into an Vector.
type VectorNode struct {
	Expressions []Node
}

var _ Node = VectorNode{}

func (v VectorNode) String() string {
	exprs := make([]string, len(v.Expressions))
	for i, expr := range v.Expressions {
		exprs[i] = expr.String()
	}
	return "[" + strings.Join(exprs, " ") + "]"
}

func (VectorNode) NodeName() string {
	return "Vector"
}

// Object is an evaluated object.
type Object struct {
	Data map[string]any
}

func (Object) NodeName() string {
	return "Object"
}

func (o Object) LiteralValue() any {
	return o.Data
}

// ObjectNode represents the parsed code for constructing an object.
// When an ObjectNode is evaluated, it turns into an Object.
type ObjectNode struct {
	Data []KeyValuePair
}

var _ Node = ObjectNode{}

func (o ObjectNode) String() string {
	pairs := make([]string, len(o.Data))
	for i, pair := range o.Data {
		pairs[i] = pair.String()
	}
	return "{" + strings.Join(pairs, " ") + "}"
}

func (ObjectNode) NodeName() string {
	return "Object"
}

type KeyValuePair struct {
	Key   Node
	Value Node
}

func (kv KeyValuePair) String() string {
	return kv.Key.String() + " " + kv.Value.String()
}

func (KeyValuePair) NodeName() string {
	return "KeyValuePair"
}

type Variable string

var _ Node = Variable("")

func (v Variable) String() string {
	return "$" + string(v)
}

func (Variable) NodeName() string {
	return "Variable"
}

type Identifier string

var _ Node = Identifier("")

func (i Identifier) String() string {
	return string(i)
}

func (Identifier) NodeName() string {
	return "Identifier"
}

type String string

var _ Node = String("")

func (s String) String() string {
	return fmt.Sprintf("%q", string(s))
}

func (String) NodeName() string {
	return "String"
}

func (s String) LiteralValue() any {
	return string(s)
}

type Number struct {
	Value any
}

var _ Node = Number{}

func (n Number) ToInteger() (int64, bool) {
	switch asserted := n.Value.(type) {
	case int:
		return int64(asserted), true
	case int32:
		return int64(asserted), true
	case int64:
		return asserted, true
	default:
		return 0, false
	}
}

func (n Number) IsFloat() bool {
	_, ok := n.Value.(float64)
	return ok
}

func (n Number) ToFloat() float64 {
	switch asserted := n.Value.(type) {
	case int:
		return float64(asserted)
	case int32:
		return float64(asserted)
	case int64:
		return float64(asserted)
	case float32:
		return float64(asserted)
	case float64:
		return asserted
	default:
		panic(fmt.Sprintf("Number with non-numeric value %v (%T)", n.Value, n.Value))
	}
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

func (n Number) LiteralValue() any {
	return n.Value
}

type Bool bool

var _ Node = Bool(false)

func (b Bool) String() string {
	if b {
		return "true"
	} else {
		return "false"
	}
}

func (Bool) NodeName() string {
	return "Bool"
}

func (b Bool) LiteralValue() any {
	return bool(b)
}

type Null struct{}

var _ Node = Null{}

func (Null) String() string {
	return "null"
}

func (Null) NodeName() string {
	return "Null"
}

func (Null) LiteralValue() any {
	return nil
}

type PathExpression struct {
	Steps []Node
}

func (e *PathExpression) Prepend(step Node) {
	e.Steps = append([]Node{step}, e.Steps...)
}

// IsIdentity returns true if the entire pathExpression was just ".".
func (e PathExpression) IsIdentity() bool {
	return len(e.Steps) == 0
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

type EvaluatedPathExpression struct {
	Steps []EvaluatedPathStep
}

func (e *EvaluatedPathExpression) Prepend(step EvaluatedPathStep) {
	e.Steps = append([]EvaluatedPathStep{step}, e.Steps...)
}

// IsIdentity returns true if the entire EvaluatedPathExpression was just ".".
func (e EvaluatedPathExpression) IsIdentity() bool {
	return len(e.Steps) == 0
}

func (e EvaluatedPathExpression) String() string {
	result := ""
	for _, step := range e.Steps {
		result += step.String()
	}

	return result
}

func (EvaluatedPathExpression) NodeName() string {
	return "EvaluatedPathExpression"
}

// type Accessor struct {
// 	Expression Node
// }

// func (a Accessor) String() string {
// 	e := a.Expression

// 	switch {
// 	case e.SymbolNode != nil:
// 		return "[" + e.SymbolNode.String() + "]"
// 	case e.TupleNode != nil:
// 		return "[" + e.TupleNode.String() + "]"
// 	case e.NumberNode != nil:
// 		return "[" + e.NumberNode.String() + "]"
// 	case e.IdentifierNode != nil:
// 		return "." + e.IdentifierNode.String()
// 	case e.StringNode != nil:
// 		return "[" + e.StringNode.String() + "]"
// 	case e.BoolNode != nil:
// 		return "[" + e.BoolNode.String() + "]"
// 	case e.NullNode != nil:
// 		return "[" + e.NullNode.String() + "]"
// 	default:
// 		return "?<unknown accessor expression>"
// 	}
// }

// func (a Accessor) NodeName() string {
// 	return "Accessor(" + a.Expression.NodeName() + ")"
// }

type EvaluatedPathStep struct {
	StringValue  *string
	IntegerValue *int64
}

func (a EvaluatedPathStep) String() string {
	switch {
	case a.StringValue != nil:
		return fmt.Sprintf("[%q]", *a.StringValue)
	case a.IntegerValue != nil:
		return fmt.Sprintf("[%d]", *a.IntegerValue)
	default:
		return "<unknown evaluatedPathStep>"
	}
}

func (a EvaluatedPathStep) NodeName() string {
	name := ""

	switch {
	case a.StringValue != nil:
		name = "String"
	case a.IntegerValue != nil:
		name = "Integer"
	default:
		name = "?"
	}

	return "EvaluatedPathStep(" + name + ")"
}
