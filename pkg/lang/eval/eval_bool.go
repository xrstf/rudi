package eval

import (
	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalBool(ctx Context, b *ast.Bool) (Context, interface{}, error) {
	return ctx, b.Value, nil
}
