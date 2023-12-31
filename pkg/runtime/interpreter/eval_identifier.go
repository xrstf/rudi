// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package interpreter

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

func (*interpreter) EvalIdentifier(ctx types.Context, ident ast.Identifier) (types.Context, any, error) {
	return ctx, nil, fmt.Errorf("unexpected identifier: %v", ident)
}
