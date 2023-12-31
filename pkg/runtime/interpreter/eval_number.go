// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package interpreter

import (
	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

func (*interpreter) EvalNumber(ctx types.Context, n ast.Number) (types.Context, any, error) {
	return ctx, n.Value, nil
}
