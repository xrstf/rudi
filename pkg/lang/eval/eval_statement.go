// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/types"
)

func evalStatement(ctx types.Context, stmt *ast.Statement) (types.Context, interface{}, error) {
	newContext, result, err := evalExpression(ctx, &stmt.Expression)
	if err != nil {
		return ctx, nil, err
	}

	fmt.Printf("%s => %#v (%T)\n", stmt.String(), result, result)

	return newContext, result, nil
}