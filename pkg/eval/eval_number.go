// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func EvalNumber(ctx types.Context, n ast.Number) (types.Context, any, error) {
	return ctx, n.Value, nil
}
