// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func EvalStatement(ctx types.Context, stmt ast.Statement) (types.Context, any, error) {
	return EvalExpression(ctx, stmt.Expression)
}
