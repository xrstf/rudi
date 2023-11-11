package eval

import (
	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalNumber(number *ast.Number, rootObject *Object) (interface{}, error) {
	return number.Value, nil
}
