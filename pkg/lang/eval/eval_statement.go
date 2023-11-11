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

	fmt.Printf("%s => %v\n", stmt.String(), result)

	return newContext, result, nil
}
