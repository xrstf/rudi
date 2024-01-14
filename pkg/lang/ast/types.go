// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package ast

import (
	"fmt"
	"regexp"
	"strings"
)

// These variables must manually be kept in-sync with the Rudi grammar.
// Hoping that https://github.com/mna/pigeon/issues/141 will provide us with a better way.
// Compared to the original grammar, these regex are anchored to make matching easier.

var (
	VariableNamePattern   = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	PathIdentifierPattern = VariableNamePattern
	IdentifierNamePattern = regexp.MustCompile(`^[a-zA-Z_+/*_%?-][a-zA-Z0-9_+/*_%?!-]*$`)
)

type Expression interface {
	String() string
	ExpressionName() string
}

type Pathed interface {
	Expression

	// GetPathExpression returns the optional path expression. Just because a type can hold
	// a path expression does not mean one is always set.
	GetPathExpression() *PathExpression

	// Pathless returns a shallow copy with the path expression omitted, e.g.
	// turning a "(foo).bar" into "(foo)".
	Pathless() Pathed
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
	return "Program"
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
var _ Pathed = Symbol{}

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

func (s Symbol) GetPathExpression() *PathExpression {
	return s.PathExpression
}

func (s Symbol) Pathless() Pathed {
	if s.Variable != nil {
		return Symbol{
			Variable: s.Variable,
		}
	}

	// for bare path expressions
	return Symbol{
		PathExpression: &PathExpression{},
	}
}

type Tuple struct {
	Expressions    []Expression
	PathExpression *PathExpression
}

var _ Expression = Tuple{}
var _ Pathed = Tuple{}

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

func (t Tuple) GetPathExpression() *PathExpression {
	return t.PathExpression
}

func (t Tuple) Pathless() Pathed {
	return Tuple{
		Expressions: t.Expressions,
	}
}

// VectorNode represents the parsed code for constructing an vector.
// When an VectorNode is evaluated, it turns into an Vector.
type VectorNode struct {
	Expressions    []Expression
	PathExpression *PathExpression
}

var _ Expression = VectorNode{}
var _ Pathed = VectorNode{}

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

func (v VectorNode) GetPathExpression() *PathExpression {
	return v.PathExpression
}

func (v VectorNode) Pathless() Pathed {
	return VectorNode{
		Expressions: v.Expressions,
	}
}

// ObjectNode represents the parsed code for constructing an object.
// When an ObjectNode is evaluated, it turns into an Object.
type ObjectNode struct {
	Data           []KeyValuePair
	PathExpression *PathExpression
}

var _ Expression = ObjectNode{}
var _ Pathed = ObjectNode{}

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

func (o ObjectNode) GetPathExpression() *PathExpression {
	return o.PathExpression
}

func (o ObjectNode) Pathless() Pathed {
	return ObjectNode{
		Data: o.Data,
	}
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
	Steps []PathStep
}

func (e *PathExpression) Prepend(step PathStep) {
	e.Steps = append([]PathStep{step}, e.Steps...)
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

func (PathExpression) ExpressionName() string {
	return "PathExpression"
}

// PathStep is a single step in a path expression. It can be either a literal
// expression (e.g. [1] or .foo or ["foo"]), or an expression that yields a
// literal (e.g. [(add 1 2)]) or a filter expression, which is evaluated within
// the context of the current element (e.g. [?(eq? .name "foo")]).
type PathStep struct {
	Expression Expression
	Filter     Expression
}

func (s PathStep) String() string {
	if s.Filter != nil {
		return "[?" + s.Filter.String() + "]"
	}

	if s.Expression != nil {
		if s, ok := s.Expression.(String); ok && PathIdentifierPattern.MatchString(string(s)) {
			return "." + string(s)
		}

		return "[" + s.Expression.String() + "]"
	}

	return "<?>"
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
