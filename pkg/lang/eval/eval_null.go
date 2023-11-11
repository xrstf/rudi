package eval

import (
	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalNull(ctx Context, n *ast.Null) (Context, interface{}, error) {
	return ctx, nil, nil
}
