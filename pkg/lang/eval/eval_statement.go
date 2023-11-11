package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalStatement(stmt *ast.Statement, rootObject *Object) (*Object, error) {
	result, err := evalExpression(&stmt.Expression, rootObject)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s => %v\n", stmt.String(), result)

	return nil, nil
}
