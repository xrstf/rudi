// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package ast

import (
	"fmt"
	"strings"
)

type Expression interface {
	String() string
	ExpressionName() string
}

type Literal interface {
	Expression

	LiteralValue() any
}

// A program is either a series of statements or a single, non-tuple expression.
type Program struct {
	Statements []Statement
}

var _ Expression = Program{}

func (p Program) String() string {
	statements := make([]string, len(p.Statements))
	for i, s := range p.Statements {
		statements[i] = s.String()
	}

	return strings.Join(statements, " ")
}

func (p Program) ExpressionName() string {
	var name string

	switch {
	case len(p.Statements) > 0:
		name = "Statements"
	default:
		name = "?"
	}

	return "Program(" + name + ")"
}

type Statement struct {
	Expression Expression
}

var _ Expression = Statement{}

func (s Statement) String() string {
	if s.Expression == nil {
		return "<invalid Statement>"
	}

	return s.Expression.String()
}

func (Statement) ExpressionName() string {
	return "Statement"
}

type Symbol struct {
	Variable       *Variable
	PathExpression *PathExpression
}

var _ Expression = Symbol{}

func (s Symbol) IsDot() bool {
	return s.Variable == nil && s.PathExpression != nil && s.PathExpression.IsIdentity()
}

func (s Symbol) String() string {
	if s.IsDot() {
		return "."
	}

	var path string
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
		return "<invalid Symbol>"
	}
}

func (s Symbol) ExpressionName() string {
	var name string

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
	Expressions    []Expression
	PathExpression *PathExpression
}

var _ Expression = Tuple{}

func (t Tuple) String() string {
	path := ""
	if t.PathExpression != nil {
		path = t.PathExpression.String()
	}

	exprs := make([]string, len(t.Expressions))
	for i, expr := range t.Expressions {
		exprs[i] = expr.String()
	}

	return "(" + strings.Join(exprs, " ") + ")" + path
}

func (Tuple) ExpressionName() string {
	return "Tuple"
}

// Vector is an evaluated vector.
type Vector struct {
	Data []any
}

var _ Expression = Vector{}
var _ Literal = Vector{}

func (v Vector) String() string {
	exprs := make([]string, len(v.Data))
	for i, expr := range v.Data {
		exprs[i] = fmt.Sprintf("%s", expr)
	}
	return "[" + strings.Join(exprs, " ") + "]"
}

func (Vector) ExpressionName() string {
	return "Vector"
}

func (v Vector) LiteralValue() any {
	return v.Data
}

func (v Vector) Clone() Vector {
	result := Vector{
		Data: make([]any, len(v.Data)),
	}
	copy(result.Data, v.Data)

	return result
}

// VectorNode represents the parsed code for constructing an vector.
// When an VectorNode is evaluated, it turns into an Vector.
type VectorNode struct {
	Expressions    []Expression
	PathExpression *PathExpression
}

var _ Expression = VectorNode{}

func (v VectorNode) String() string {
	path := ""
	if v.PathExpression != nil {
		path = v.PathExpression.String()
	}

	exprs := make([]string, len(v.Expressions))
	for i, expr := range v.Expressions {
		exprs[i] = expr.String()
	}

	return "[" + strings.Join(exprs, " ") + "]" + path
}

func (VectorNode) ExpressionName() string {
	return "Vector"
}

// Object is an evaluated object.
type Object struct {
	Data map[string]any
}

var _ Expression = Object{}
var _ Literal = Object{}

func (o Object) String() string {
	exprs := make([]string, 0)
	for key, value := range o.Data {
		exprs = append(exprs, fmt.Sprintf("%q: %v", key, value))
	}
	return "{" + strings.Join(exprs, ", ") + "}"
}

func (Object) ExpressionName() string {
	return "Object"
}

func (o Object) LiteralValue() any {
	return o.Data
}

// ObjectNode represents the parsed code for constructing an object.
// When an ObjectNode is evaluated, it turns into an Object.
type ObjectNode struct {
	Data           []KeyValuePair
	PathExpression *PathExpression
}

var _ Expression = ObjectNode{}

func (o ObjectNode) String() string {
	pairs := make([]string, len(o.Data))
	for i, pair := range o.Data {
		pairs[i] = pair.String()
	}
	return "{" + strings.Join(pairs, " ") + "}"
}

func (ObjectNode) ExpressionName() string {
	return "Object"
}

type KeyValuePair struct {
	Key   Expression
	Value Expression
}

func (kv KeyValuePair) String() string {
	return kv.Key.String() + " " + kv.Value.String()
}

func (KeyValuePair) ExpressionName() string {
	return "KeyValuePair"
}

type Variable string

var _ Expression = Variable("")

func (v Variable) String() string {
	return "$" + string(v)
}

func (Variable) ExpressionName() string {
	return "Variable"
}

type Identifier struct {
	Name string
	Bang bool
}

var _ Expression = Identifier{}

func (i Identifier) Equal(other Identifier) bool {
	return i.Name == other.Name
}

func (i Identifier) String() string {
	result := i.Name
	if i.Bang {
		result += "!"
	}

	return result
}

func (i Identifier) ExpressionName() string {
	result := "Identifier"
	if i.Bang {
		result += "!"
	}

	return result
}

type String string

var _ Expression = String("")
var _ Literal = String("")

func (s String) Equal(other String) bool {
	return string(s) == string(other)
}

func (s String) String() string {
	return fmt.Sprintf("%q", string(s))
}

func (String) ExpressionName() string {
	return "String"
}

func (s String) LiteralValue() any {
	return string(s)
}

type Number struct {
	Value any
}

var _ Expression = Number{}
var _ Literal = Number{}

func (n Number) Equal(other Number) bool {
	// handle technically invalid numbers
	if other.Value == nil || n.Value == nil {
		return (other.Value == nil) == (n.Value == nil)
	}

	selfInt, selfOk := n.ToInteger()
	otherInt, otherOk := other.ToInteger()

	// not the same type
	if selfOk != otherOk {
		return false
	}

	if selfOk {
		return otherInt == selfInt
	}

	return n.ToFloat() == other.ToFloat()
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

func (Number) ExpressionName() string {
	return "Number"
}

func (n Number) LiteralValue() any {
	return n.Value
}

type Bool bool

var _ Expression = Bool(false)
var _ Literal = Bool(false)

func (b Bool) Equal(other Bool) bool {
	return bool(b) == bool(other)
}

func (b Bool) String() string {
	if b {
		return "true"
	} else {
		return "false"
	}
}

func (Bool) ExpressionName() string {
	return "Bool"
}

func (b Bool) LiteralValue() any {
	return bool(b)
}

type Null struct{}

var _ Expression = Null{}
var _ Literal = Null{}

func (Null) Equal(other Null) bool {
	return true
}

func (Null) String() string {
	return "null"
}

func (Null) ExpressionName() string {
	return "Null"
}

func (Null) LiteralValue() any {
	return nil
}

type PathExpression struct {
	Steps []Expression
}

func (e *PathExpression) Prepend(step Expression) {
	e.Steps = append([]Expression{step}, e.Steps...)
}

// IsIdentity returns true if the entire pathExpression was just ".".
func (e PathExpression) IsIdentity() bool {
	return len(e.Steps) == 0
}

func (e PathExpression) String() string {
	result := ""
	for _, step := range e.Steps {
		if ident, ok := step.(Identifier); ok {
			result += "." + ident.String()
		} else {
			result += "[" + step.String() + "]"
		}
	}

	return result
}

func (PathExpression) ExpressionName() string {
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

func (EvaluatedPathExpression) ExpressionName() string {
	return "PathExpression"
}

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
		return "<unknown PathStep>"
	}
}

func (a EvaluatedPathStep) ExpressionName() string {
	var name string

	switch {
	case a.StringValue != nil:
		name = "String"
	case a.IntegerValue != nil:
		name = "Number"
	default:
		name = "?"
	}

	return "PathStep(" + name + ")"
}
