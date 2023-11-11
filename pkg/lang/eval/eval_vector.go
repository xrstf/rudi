package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalVector(vec *ast.Vector, rootObject *Object) (interface{}, error) {
	result := make([]interface{}, len(vec.Expressions))

	for i, expr := range vec.Expressions {
		data, err := evalExpression(&expr, rootObject)
		if err != nil {
			return nil, fmt.Errorf("failed to eval expression %s: %w", expr.String(), err)
		}

		result[i] = data
	}

	return result, nil
}
