package eval

import (
	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalString(str *ast.String, rootObject *Object) (interface{}, error) {
	return str.Value, nil
}
