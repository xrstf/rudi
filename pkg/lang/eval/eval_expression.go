package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalExpression(expr *ast.Expression, rootObject *Object) (interface{}, error) {
	switch {
	case expr.NullNode != nil:
		return evalNull(expr.NullNode, rootObject)
	case expr.BoolNode != nil:
		return evalBool(expr.BoolNode, rootObject)
	case expr.StringNode != nil:
		return evalString(expr.StringNode, rootObject)
	case expr.NumberNode != nil:
		return evalNumber(expr.NumberNode, rootObject)
	case expr.ObjectNode != nil:
		return evalObject(expr.ObjectNode, rootObject)
	case expr.VectorNode != nil:
		return evalVector(expr.VectorNode, rootObject)
	case expr.SymbolNode != nil:
		return evalSymbol(expr.SymbolNode, rootObject)
	case expr.TupleNode != nil:
		return evalTuple(expr.TupleNode, rootObject)
	}

	return nil, fmt.Errorf("unknown expression %T (%s)", expr, expr.String())
}
