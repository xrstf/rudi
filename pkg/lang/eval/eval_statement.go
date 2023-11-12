// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalStatement(ctx Context, stmt *ast.Statement) (Context, interface{}, error) {
	newContext, result, err := evalExpression(ctx, &stmt.Expression)
	if err != nil {
		return ctx, nil, err
	}

	fmt.Printf("%s => %#v (%T)\n", stmt.String(), result, result)

	return newContext, result, nil
}
