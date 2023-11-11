package eval

import (
	"errors"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalIdentifier(ctx Context, ident *ast.Identifier) (Context, interface{}, error) {
	return ctx, nil, errors.New("identifiers cannot be evaluated on themselves")
}
