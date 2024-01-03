// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package types

import (
	"go.xrstf.de/rudi/pkg/lang/ast"
)

type Runtime interface {
	EvalNull(ctx Context, n ast.Null) (any, error)
	EvalBool(ctx Context, b ast.Bool) (any, error)
	EvalNumber(ctx Context, n ast.Number) (any, error)
	EvalString(ctx Context, str ast.String) (any, error)
	EvalSymbol(ctx Context, sym ast.Symbol) (any, error)
	EvalVectorNode(ctx Context, vec ast.VectorNode) (any, error)
	EvalObjectNode(ctx Context, obj ast.ObjectNode) (any, error)
	EvalIdentifier(ctx Context, ident ast.Identifier) (any, error)
	EvalExpression(ctx Context, expr ast.Expression) (any, error)
	EvalTuple(ctx Context, tup ast.Tuple) (any, error)
	EvalStatement(ctx Context, stmt ast.Statement) (any, error)
	EvalProgram(ctx Context, p *ast.Program) (any, error)

	CallFunction(ctx Context, fun ast.Identifier, args []ast.Expression) (any, error)
}
