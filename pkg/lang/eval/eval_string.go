// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalString(ctx Context, str *ast.String) (Context, interface{}, error) {
	return ctx, str.Value, nil
}
