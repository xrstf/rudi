// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/types"
)

func evalNumber(ctx types.Context, n *ast.Number) (types.Context, interface{}, error) {
	return ctx, n.Value, nil
}
