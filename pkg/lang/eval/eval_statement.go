// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

func EvalStatement(ctx types.Context, stmt ast.Statement) (types.Context, any, error) {
	newContext, result, err := EvalExpression(ctx, stmt.Expression)
	if err != nil {
		return ctx, nil, err
	}

	return newContext, result, nil
}
