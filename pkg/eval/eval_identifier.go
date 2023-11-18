// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/otto/pkg/eval/types"
	"go.xrstf.de/otto/pkg/lang/ast"
)

func EvalIdentifier(ctx types.Context, ident ast.Identifier) (types.Context, any, error) {
	return ctx, nil, fmt.Errorf("unexpected identifier: %v", ident)
}
