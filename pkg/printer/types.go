// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package printer

import (
	"io"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

var Indent = "  "

const DoNotIndent = -1

type Object []KeyValuePair

type KeyValuePair struct {
	Key   any
	Value any
}

type Renderer interface {
	WriteSingleline(v any, out io.Writer) error
	WriteMultiline(v any, out io.Writer) error
	WriteNull(out io.Writer) error
	WriteBool(b bool, out io.Writer) error
	WriteNumber(value any, out io.Writer) error
	WriteString(str string, out io.Writer) error
	WriteIdentifier(ident *ast.Identifier, out io.Writer) error
	WriteVectorSingleline(vec []any, pathExpr *ast.PathExpression, out io.Writer, depth int) error
	WriteVectorMultiline(vec []any, pathExpr *ast.PathExpression, out io.Writer, depth int) error
	WriteObjectSingleline(obj Object, pathExpr *ast.PathExpression, out io.Writer, depth int) error
	WriteObjectMultiline(obj Object, pathExpr *ast.PathExpression, out io.Writer, depth int) error
	WriteTupleSingleline(tup *ast.Tuple, out io.Writer, depth int) error
	WriteTupleMultiline(tup *ast.Tuple, out io.Writer, depth int) error
	WriteSymbol(sym *ast.Symbol, out io.Writer, depth int) error
}
