package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalSymbol(ctx Context, sym *ast.Symbol) (Context, interface{}, error) {
	switch {
	case sym.Identifier != nil:
		return ctx, sym.Identifier, nil

	case sym.Variable != nil:
		varName := sym.Variable.Name

		value, ok := ctx.Variables[varName]
		if !ok {
			return ctx, nil, fmt.Errorf("unknown variable %s", varName)
		}

		return ctx, value, nil
	}

	return ctx, nil, fmt.Errorf("unknown symbol %T (%s)", sym, sym.String())
}
