// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package printer

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

type Object []KeyValuePair

type KeyValuePair struct {
	Key   any
	Value any
}

type Printer interface {
	Print(v any) error
	Null() error
	Bool(b bool) error
	Number(value any) error
	String(str string) error
	Vector(vec []any) error
	VectorNode(vec *ast.VectorNode) error
	Object(obj map[string]any) error
	ObjectNode(obj *ast.ObjectNode) error
	Identifier(ident *ast.Identifier) error
	Symbol(sym *ast.Symbol) error
	Tuple(tup *ast.Tuple) error
	Expression(expr ast.Expression) error
	Statement(tup *ast.Statement) error
	Program(tup *ast.Program) error
}

func printAny(val any, r Printer) error {
	switch asserted := val.(type) {
	case nil:
		return r.Null()
	case ast.Null:
		return r.Null()
	case bool:
		return r.Bool(asserted)
	case ast.Bool:
		return r.Bool(bool(asserted))
	case int:
		return r.Number(asserted)
	case int32:
		return r.Number(asserted)
	case int64:
		return r.Number(asserted)
	case float32:
		return r.Number(asserted)
	case float64:
		return r.Number(asserted)
	case ast.Number:
		return r.Number(asserted.Value)
	case string:
		return r.String(asserted)
	case ast.String:
		return r.String(string(asserted))
	case []any:
		return r.Vector(asserted)
	case ast.VectorNode:
		return r.VectorNode(&asserted)
	case map[string]any:
		return r.Object(asserted)
	case ast.ObjectNode:
		return r.ObjectNode(&asserted)
	case ast.Symbol:
		return r.Symbol(&asserted)
	case ast.Tuple:
		return r.Tuple(&asserted)
	case ast.Identifier:
		return r.Identifier(&asserted)
	case ast.Statement:
		return r.Statement(&asserted)
	case ast.Program:
		return r.Program(&asserted)
	}

	return fmt.Errorf("cannot dump values of type %T", val)
}
