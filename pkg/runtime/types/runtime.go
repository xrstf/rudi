// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package types

import (
	"go.xrstf.de/rudi/pkg/lang/ast"
)

type Runtime interface {
	EvalNull(ctx Context, n ast.Null) (Context, any, error)
	EvalBool(ctx Context, b ast.Bool) (Context, any, error)
	EvalNumber(ctx Context, n ast.Number) (Context, any, error)
	EvalString(ctx Context, str ast.String) (Context, any, error)
	EvalSymbol(ctx Context, sym ast.Symbol) (Context, any, error)
	EvalVectorNode(ctx Context, vec ast.VectorNode) (Context, any, error)
	EvalObjectNode(ctx Context, obj ast.ObjectNode) (Context, any, error)
	EvalIdentifier(ctx Context, ident ast.Identifier) (Context, any, error)
	EvalExpression(ctx Context, expr ast.Expression) (Context, any, error)
	EvalTuple(ctx Context, tup ast.Tuple) (Context, any, error)
	EvalStatement(ctx Context, stmt ast.Statement) (Context, any, error)
	EvalProgram(ctx Context, p *ast.Program) (Context, any, error)

	CallFunction(ctx Context, fun ast.Identifier, args []ast.Expression) (Context, any, error)
}
