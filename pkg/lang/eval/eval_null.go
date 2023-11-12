// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/types"
)

func evalNull(ctx types.Context, n *ast.Null) (types.Context, any, error) {
	return ctx, *n, nil
}
