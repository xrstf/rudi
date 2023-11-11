package eval

import (
	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalNumber(ctx Context, n *ast.Number) (Context, interface{}, error) {
	return ctx, n.Value, nil
}
