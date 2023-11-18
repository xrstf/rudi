// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"go.xrstf.de/otto/pkg/eval/types"
	"go.xrstf.de/otto/pkg/lang/ast"
)

func EvalBool(ctx types.Context, b ast.Bool) (types.Context, any, error) {
	return ctx, b, nil
}
