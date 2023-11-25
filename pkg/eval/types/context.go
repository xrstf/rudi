// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package types

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

type Document struct {
	data any
}

func NewDocument(data any) (Document, error) {
	wrapped, err := WrapNative(data)
	if err != nil {
		return Document{}, fmt.Errorf("invalid document data: %w", err)
	}

	return Document{
		data: wrapped,
	}, nil
}

func (d *Document) Data() any {
	return d.data
}

func (d *Document) Set(wrappedData any) {
	d.data = wrappedData
}

type Context struct {
	document  *Document
	funcs     Functions
	variables Variables
}

func NewContext(doc Document, variables Variables, funcs Functions) Context {
	if funcs == nil {
		funcs = NewFunctions()
	}

	if variables == nil {
		variables = NewVariables()
	}

	return Context{
		document:  &doc,
		funcs:     funcs,
		variables: variables,
	}
}

func (c Context) GetDocument() *Document {
	return c.document
}

func (c Context) GetVariable(name string) (any, bool) {
	return c.variables.Get(name)
}

func (c Context) GetFunction(name string) (Function, bool) {
	return c.funcs.Get(name)
}

func (c Context) WithVariable(name string, val any) Context {
	return Context{
		document:  c.document,
		funcs:     c.funcs,
		variables: c.variables.With(name, val),
	}
}

type Function interface {
	Evaluate(ctx Context, args []ast.Expression) (any, error)
	Description() string
}

type TupleFunction func(ctx Context, args []ast.Expression) (any, error)

type basicFunc struct {
	f    TupleFunction
	desc string
}

func BasicFunction(f TupleFunction, description string) Function {
	return basicFunc{
		f:    f,
		desc: description,
	}
}

var _ Function = basicFunc{}

func (f basicFunc) Evaluate(ctx Context, args []ast.Expression) (any, error) {
	return f.f(ctx, args)
}

func (f basicFunc) Description() string {
	return f.desc
}

type Functions map[string]Function

func NewFunctions() Functions {
	return Functions{}
}

func (f Functions) Get(name string) (Function, bool) {
	variable, exists := f[name]
	return variable, exists
}

// Set sets/replaces the function in the current set (in-place).
// The function returns the same Functions to allow fluent access.
func (f Functions) Set(name string, fun Function) Functions {
	f[name] = fun
	return f
}

func (f Functions) DeepCopy() Functions {
	result := NewFunctions()
	for key, val := range f {
		result[key] = val
	}
	return result
}

type Variables map[string]any

func NewVariables() Variables {
	return Variables{}
}

func (v Variables) Get(name string) (any, bool) {
	variable, exists := v[name]
	return variable, exists
}

// Set sets/replaces the variable value in the current set (in-place).
// The function returns the same variables to allow fluent access.
func (v Variables) Set(name string, val any) Variables {
	v[name] = val
	return v
}

// With returns a copy of the variables, with the new variable being added to it.
func (v Variables) With(name string, val any) Variables {
	return v.DeepCopy().Set(name, val)
}

func (v Variables) DeepCopy() Variables {
	result := NewVariables()
	for key, val := range v {
		result[key] = val
	}
	return result
}

func WrapNative(val any) (ast.Literal, error) {
	switch asserted := val.(type) {
	case nil:
		return ast.Null{}, nil
	case ast.Null:
		return asserted, nil
	case string:
		return ast.String(asserted), nil
	case ast.String:
		return asserted, nil
	case bool:
		return ast.Bool(asserted), nil
	case ast.Bool:
		return asserted, nil
	case int:
		return ast.Number{Value: int64(asserted)}, nil
	case int32:
		return ast.Number{Value: int64(asserted)}, nil
	case int64:
		return ast.Number{Value: asserted}, nil
	case float32:
		return ast.Number{Value: float64(asserted)}, nil
	case float64:
		return ast.Number{Value: asserted}, nil
	case ast.Number:
		return asserted, nil
	case []any:
		return ast.Vector{Data: asserted}, nil
	case ast.Vector:
		return asserted, nil
	case map[string]any:
		return ast.Object{Data: asserted}, nil
	case ast.Object:
		return asserted, nil
	default:
		return nil, fmt.Errorf("cannot wrap %v (%T)", val, val)
	}
}

func Must[T any](val T, _ error) T {
	return val
}

func UnwrapType(val any) (any, error) {
	switch asserted := val.(type) {
	case ast.Null:
		return nil, nil
	case *ast.Null:
		return nil, nil
	case nil:
		return nil, nil
	case ast.Bool:
		return bool(asserted), nil
	case *ast.Bool:
		return bool(*asserted), nil
	case bool:
		return asserted, nil
	case ast.String:
		return string(asserted), nil
	case *ast.String:
		return string(*asserted), nil
	case string:
		return asserted, nil
	case ast.Number:
		return asserted.Value, nil
	case *ast.Number:
		return asserted.Value, nil
	case int:
		return int64(asserted), nil
	case int32:
		return int64(asserted), nil
	case int64:
		return asserted, nil
	case float32:
		return float64(asserted), nil
	case float64:
		return asserted, nil
	case ast.Vector:
		return unwrapVector(&asserted)
	case *ast.Vector:
		return unwrapVector(asserted)
	case []any:
		return unwrapVector(&ast.Vector{Data: asserted})
	case ast.Object:
		return unwrapObject(&asserted)
	case *ast.Object:
		return unwrapObject(asserted)
	case map[string]any:
		return unwrapObject(&ast.Object{Data: asserted})
	default:
		return nil, fmt.Errorf("cannot unwrap %v (%T)", val, val)
	}
}

func unwrapVector(v *ast.Vector) ([]any, error) {
	result := make([]any, len(v.Data))
	for i, item := range v.Data {
		unwrappedItem, err := UnwrapType(item)
		if err != nil {
			return nil, err
		}

		result[i] = unwrappedItem
	}

	return result, nil
}

func unwrapObject(o *ast.Object) (map[string]any, error) {
	result := map[string]any{}
	for key, value := range o.Data {
		unwrappedValue, err := UnwrapType(value)
		if err != nil {
			return nil, err
		}

		result[key] = unwrappedValue
	}

	return result, nil
}
