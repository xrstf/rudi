// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

func EvalStatement(ctx types.Context, stmt ast.Statement) (types.Context, any, error) {
	newContext, result, err := EvalTuple(ctx, stmt.Tuple)
	if err != nil {
		return ctx, nil, err
	}

	fmt.Printf("%s => %#v (%T)\n", stmt.String(), result, result)

	return newContext, result, nil
}
