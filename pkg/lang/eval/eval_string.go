// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/types"
)

func evalString(ctx types.Context, str *ast.String) (types.Context, interface{}, error) {
	return ctx, str.Value, nil
}
