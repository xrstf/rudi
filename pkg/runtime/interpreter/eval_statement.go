// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package interpreter

import (
	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

func (i *interpreter) EvalStatement(ctx types.Context, stmt ast.Statement) (types.Context, any, error) {
	return i.EvalExpression(ctx, stmt.Expression)
}
