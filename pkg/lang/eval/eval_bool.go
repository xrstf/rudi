package eval

import (
	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalBool(b *ast.Bool, rootObject *Object) (interface{}, error) {
	return b.Value, nil
}
