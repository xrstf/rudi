// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package printer

import (
	"fmt"
	"io"
	"strings"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

func writeAny(val any, r Renderer, out io.Writer, depth int) error {
	switch asserted := val.(type) {
	case nil:
		return r.WriteNull(out)
	case ast.Null:
		return r.WriteNull(out)
	case bool:
		return r.WriteBool(asserted, out)
	case ast.Bool:
		return r.WriteBool(bool(asserted), out)
	case int:
		return r.WriteNumber(asserted, out)
	case int32:
		return r.WriteNumber(asserted, out)
	case int64:
		return r.WriteNumber(asserted, out)
	case float32:
		return r.WriteNumber(asserted, out)
	case float64:
		return r.WriteNumber(asserted, out)
	case ast.Number:
		return r.WriteNumber(asserted.Value, out)
	case string:
		return r.WriteString(asserted, out)
	case ast.String:
		return r.WriteString(string(asserted), out)
	case []any:
		return dumpVector(asserted, nil, r, out, depth)
	case ast.VectorNode:
		return dumpVectorNode(asserted, r, out, depth)
	case map[string]any:
		return dumpObject(asserted, r, out, depth)
	case ast.ObjectNode:
		return dumpObjectNode(asserted, r, out, depth)
	case ast.Symbol:
		return r.WriteSymbol(&asserted, out, depth)
	case ast.Tuple:
		return writeTuple(&asserted, r, out, depth)
	case ast.Identifier:
		return r.WriteIdentifier(&asserted, out)
	case ast.Statement:
		return writeStatement(&asserted, r, out, depth)
	case ast.Program:
		return writeProgram(&asserted, r, out, depth)
	}

	return fmt.Errorf("cannot dump values of type %T", val)
}

func writeTuple(tup *ast.Tuple, r Renderer, out io.Writer, depth int) error {
	if depth == DoNotIndent || len(tup.Expressions) == 0 {
		return r.WriteTupleSingleline(tup, out, depth)
	}

	// check if we can in-line or if we need to put each element on its own line
	var buf strings.Builder
	for _, expr := range tup.Expressions {
		if err := writeAny(expr, r, &buf, 0); err != nil {
			return err
		}
	}

	if buf.Len() > 50 || depth > 10 {
		return r.WriteTupleMultiline(tup, out, depth)
	} else {
		return r.WriteTupleSingleline(tup, out, depth)
	}
}

func writeStatement(stmt *ast.Statement, r Renderer, out io.Writer, depth int) error {
	return writeAny(stmt.Expression, r, out, depth)
}

func writeProgram(p *ast.Program, r Renderer, out io.Writer, depth int) error {
	for _, stmt := range p.Statements {
		if err := writeStatement(&stmt, r, out, depth); err != nil {
			return fmt.Errorf("failed to dump statement: %w", err)
		}

		separator := "\n"
		if depth == DoNotIndent {
			separator = " "
		}

		if err := writeString(out, separator); err != nil {
			return err
		}
	}

	return nil
}

func dumpVector(vec []any, pathExpr *ast.PathExpression, r Renderer, out io.Writer, depth int) error {
	if depth == DoNotIndent || len(vec) == 0 {
		return r.WriteVectorSingleline(vec, pathExpr, out, depth)
	}

	// check if we can in-line or if we need to put each element on its own line
	var buf strings.Builder
	for _, val := range vec {
		if err := writeAny(val, r, &buf, 0); err != nil {
			return err
		}
	}

	if buf.Len() > 50 || depth > 10 {
		return r.WriteVectorMultiline(vec, pathExpr, out, depth)
	} else {
		return r.WriteVectorSingleline(vec, pathExpr, out, depth)
	}
}

func dumpVectorNode(vec ast.VectorNode, r Renderer, out io.Writer, depth int) error {
	data := make([]any, len(vec.Expressions))
	for i, expr := range vec.Expressions {
		data[i] = expr
	}

	return dumpVector(data, vec.PathExpression, r, out, depth)
}

func dumpGenericObject(obj Object, pathExpr *ast.PathExpression, r Renderer, out io.Writer, depth int) error {
	if depth == DoNotIndent || len(obj) == 0 {
		return r.WriteObjectSingleline(obj, pathExpr, out, depth)
	}

	// check if we can in-line or if we need to put each element on its own line
	var buf strings.Builder
	for _, pair := range obj {
		if err := writeAny(pair.Key, r, &buf, 0); err != nil {
			return err
		}

		if err := writeAny(pair.Value, r, &buf, 0); err != nil {
			return err
		}
	}

	if buf.Len() > 50 || depth > 10 {
		return r.WriteObjectMultiline(obj, pathExpr, out, depth)
	} else {
		return r.WriteObjectSingleline(obj, pathExpr, out, depth)
	}
}

func mapToObject(m map[string]any) Object {
	out := make(Object, len(m))

	i := 0
	for k, v := range m {
		out[i] = KeyValuePair{
			Key:   k,
			Value: v,
		}
		i++
	}

	return out
}

func dumpObject(obj map[string]any, r Renderer, out io.Writer, depth int) error {
	return dumpGenericObject(mapToObject(obj), nil, r, out, depth)
}

func objectNodeToObject(o ast.ObjectNode) Object {
	out := make(Object, len(o.Data))

	for i, pair := range o.Data {
		out[i] = KeyValuePair{
			Key:   pair.Key,
			Value: pair.Value,
		}
	}

	return out
}

func dumpObjectNode(obj ast.ObjectNode, r Renderer, out io.Writer, depth int) error {
	return dumpGenericObject(objectNodeToObject(obj), obj.PathExpression, r, out, depth)
}

func writeString(out io.Writer, str string) error {
	_, err := out.Write([]byte(str))
	return err
}
