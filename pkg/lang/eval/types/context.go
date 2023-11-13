// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package types

import (
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
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

func (d *Document) Get() any {
	return d.data
}

func (d *Document) Set(wrappedData any) {
	d.data = wrappedData
}

type Context struct {
	document  Document
	variables Variables
}

func NewContext(doc Document, variables Variables) Context {
	return Context{
		document:  doc,
		variables: variables,
	}
}

func (c Context) GetDocument() Document {
	return c.document
}

func (c Context) GetVariable(name string) (any, bool) {
	return c.variables.Get(name)
}

func (c Context) WithVariable(name string, val any) Context {
	return Context{
		document:  c.document,
		variables: c.variables.With(name, val),
	}
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

func WrapNative(val any) (any, error) {
	switch asserted := val.(type) {
	case nil:
		return ast.Null{}, nil
	case ast.Null:
		return val, nil
	case string:
		return ast.String(asserted), nil
	case ast.String:
		return val, nil
	case bool:
		return ast.Bool(asserted), nil
	case ast.Bool:
		return val, nil
	case int:
		return ast.Number{Value: int64(asserted)}, nil
	case int32:
		return ast.Number{Value: int64(asserted)}, nil
	case int64:
		return ast.Number{Value: int64(asserted)}, nil
	case float32:
		return ast.Number{Value: float64(asserted)}, nil
	case float64:
		return ast.Number{Value: float64(asserted)}, nil
	case ast.Number:
		return val, nil
	case []any:
		return ast.Vector{Data: asserted}, nil
	case ast.Vector:
		return val, nil
	case map[string]any:
		return ast.Object{Data: asserted}, nil
	case ast.Object:
		return val, nil
	default:
		return nil, fmt.Errorf("cannot wrap %v (%T)", val, val)
	}
}

func Must(val any, _ error) any {
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
	case int64:
		return asserted, nil
	case float64:
		return asserted, nil
	case ast.Vector:
		return asserted.Data, nil
	case *ast.Vector:
		return asserted.Data, nil
	case []any:
		return asserted, nil
	case ast.Object:
		return asserted.Data, nil
	case *ast.Object:
		return asserted.Data, nil
	case map[string]any:
		return asserted, nil
	default:
		return nil, fmt.Errorf("cannot unwrap %v (%T)", val, val)
	}
}
