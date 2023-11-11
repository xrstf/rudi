package eval

import (
	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalNull(n *ast.Null, rootObject *Object) (interface{}, error) {
	return nil, nil
}
