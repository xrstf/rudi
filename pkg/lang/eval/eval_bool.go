// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/types"
)

func evalBool(ctx types.Context, b *ast.Bool) (types.Context, any, error) {
	return ctx, *b, nil
}
