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
	VectorNode     *VectorNode
	ObjectNode     *ObjectNode
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
	Expressions []Expression
}

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
	Key   Expression
	Value Expression
}

func (kv KeyValuePair) String() string {
	return kv.Key.String() + " " + kv.Value.String()
}

func (KeyValuePair) NodeName() string {
	return "KeyValuePair"
}

type Variable string

func (v Variable) String() string {
	return "$" + string(v)
}

func (Variable) NodeName() string {
	return "Variable"
}

type Identifier string

func (i Identifier) String() string {
	return string(i)
}

func (Identifier) NodeName() string {
	return "Identifier"
}

type String string

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
	Steps []Accessor
}

func (e *PathExpression) Prepend(step Accessor) {
	e.Steps = append([]Accessor{step}, e.Steps...)
}

// IsIdentity returns true if the entire pathExpression was just ".".
func (e PathExpression) IsIdentity() bool {
	return len(e.Steps) == 1 && e.Steps[0].IsIdentity()
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
	Steps []EvaluatedAccessor
}

func (e *EvaluatedPathExpression) Prepend(step EvaluatedAccessor) {
	e.Steps = append([]EvaluatedAccessor{step}, e.Steps...)
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

type Accessor struct {
	Identifier *Identifier
	StringNode *String
	Variable   *Variable
	Tuple      *Tuple
	Integer    *int64
}

// IsIdentity returns true if the accessor is for the current document (i.e. the entire pathExpression
// was just ".").
func (a Accessor) IsIdentity() bool {
	return true &&
		a.Identifier == nil &&
		a.StringNode == nil &&
		a.Variable == nil &&
		a.Tuple == nil &&
		a.Integer == nil
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
		return "."
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

type EvaluatedAccessor struct {
	StringValue  *string
	IntegerValue *int64
}

func (a EvaluatedAccessor) String() string {
	switch {
	case a.StringValue != nil:
		return fmt.Sprintf("[%q]", *a.StringValue)
	case a.IntegerValue != nil:
		return fmt.Sprintf("[%d]", *a.IntegerValue)
	default:
		return "<unknown evaluatedAccessor>"
	}
}

func (a EvaluatedAccessor) NodeName() string {
	name := ""

	switch {
	case a.StringValue != nil:
		name = "String"
	case a.IntegerValue != nil:
		name = "Integer"
	default:
		name = "?"
	}

	return "EvaluatedAccessor(" + name + ")"
}
