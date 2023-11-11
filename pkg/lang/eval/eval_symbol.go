package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalSymbol(sym *ast.Symbol, rootObject *Object) (interface{}, error) {
	switch {
	case sym.Identifier != nil:
		return sym.Identifier, nil

	case sym.Variable != nil:
		varName := sym.Variable.Name

		value, ok := runtimeVariables[varName]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", varName)
		}

		return value, nil
	}

	return nil, fmt.Errorf("unknown symbol %T (%s)", sym, sym.String())
}
