// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func EvalString(ctx types.Context, str ast.String) (types.Context, any, error) {
	return ctx, string(str), nil
}
