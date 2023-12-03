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

func (s String) Equal(other String) bool {
	return string(s) == string(other)
}

func (s String) String() string {
	return fmt.Sprintf("%q", string(s))
}

func (String) ExpressionName() string {
	return "String"
}

type Number struct {
	Value any
}

var _ Expression = Number{}

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

func (n Number) ToFloat() (float64, bool) {
	switch asserted := n.Value.(type) {
	case float32:
		return float64(asserted), true
	case float64:
		return asserted, true
	default:
		return 0, false
	}
}

func (n Number) MustToFloat() float64 {
	if i, ok := n.ToInteger(); ok {
		return float64(i)
	}

	if f, ok := n.ToFloat(); ok {
		return f
	}

	panic(fmt.Sprintf("invalid number value %#v (%T)", n.Value, n.Value))
}

func (n Number) String() string {
	if f, ok := n.ToFloat(); ok {
		return fmt.Sprintf("%f", f)
	}

	return fmt.Sprintf("%d", n.Value)
}

func (Number) ExpressionName() string {
	return "Number"
}

type Bool bool

var _ Expression = Bool(false)

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

type Null struct{}

var _ Expression = Null{}

func (Null) Equal(other Null) bool {
	return true
}

func (Null) String() string {
	return "null"
}

func (Null) ExpressionName() string {
	return "Null"
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

// Shims are used to turn any Go value into a Rudi expression. This is done when
// constructing new expressions and tuples at runtime. A Rudi program itself can
// never contain Shim nodes.
type Shim struct {
	Value any
}

var _ Expression = Shim{}

func (s Shim) Equal(other Shim) bool {
	return s.Value == other.Value
}

func (Shim) String() string {
	return "Shim"
}

func (Shim) ExpressionName() string {
	return "Shim"
}
