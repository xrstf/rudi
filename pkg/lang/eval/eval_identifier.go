// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"errors"

	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/types"
)

func evalIdentifier(ctx types.Context, ident *ast.Identifier) (types.Context, interface{}, error) {
	return ctx, nil, errors.New("unexpected identifier")
}
