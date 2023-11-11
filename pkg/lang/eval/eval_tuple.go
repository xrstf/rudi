package eval

import (
	"errors"
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/builtin"
	"go.xrstf.de/corel/pkg/lang/eval/coalescing"
)

func evalTuple(tup *ast.Tuple, rootObject *Object) (interface{}, error) {
	funcName := tup.Identifier.Name

	// hardcode root behaviour for those tuples where not all
	// expressions can be pre-computed (in case, for example,
	// the else-path of an if statement would have side effects)
	switch funcName {
	case "if":
		return evalIfTuple(tup, rootObject)
	case "def":
		return evalDefTuple(tup, rootObject)
	}

	function, ok := builtin.Functions[funcName]
	if !ok {
		return nil, fmt.Errorf("unknown function %s", funcName)
	}

	// evaluate all function arguments
	args := make([]interface{}, len(tup.Expressions))
	for i, expr := range tup.Expressions {
		arg, err := evalExpression(&expr, rootObject)
		if err != nil {
			return nil, fmt.Errorf("invalid argument %d: %w", i, err)
		}

		args[i] = arg
	}

	// call the function
	result, err := function(args)
	if err != nil {
		return nil, fmt.Errorf("function failed: %w", err)
	}

	return result, nil
}

func evalIfTuple(tup *ast.Tuple, rootObject *Object) (interface{}, error) {
	if size := len(tup.Expressions); size != 2 && size != 3 {
		return nil, fmt.Errorf("invalid if tuple: expected 2 or 3 expressions, but got %d", size)
	}

	condition, err := evalExpression(&tup.Expressions[0], rootObject)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate condition: %w", err)
	}

	success, err := coalescing.ToBool(condition)
	if err != nil {
		return nil, fmt.Errorf("condition did not return boolish value: %w", err)
	}

	if success {
		return evalExpression(&tup.Expressions[1], rootObject)
	}

	// optional else part
	if len(tup.Expressions) > 2 {
		return evalExpression(&tup.Expressions[2], rootObject)
	}

	return rootObject, nil
}

// this should be scoped, but for testing we just use one global map
var runtimeVariables = map[string]interface{}{}

func evalDefTuple(tup *ast.Tuple, rootObject *Object) (interface{}, error) {
	if size := len(tup.Expressions); size != 2 {
		return nil, fmt.Errorf("invalid set tuple: expected exactly 2 expressions, but got %d", size)
	}

	varNameExpr := tup.Expressions[0]
	varName := ""

	if varNameExpr.SymbolNode == nil {
		return nil, errors.New("(def $varname expression)")
	}

	symNode := varNameExpr.SymbolNode
	if symNode.Variable == nil {
		return nil, errors.New("(def $varname expression)")
	}

	// forbid weird definitions like (def $var.foo (expr)) for now
	if symNode.PathExpression != nil {
		return nil, errors.New("(def $varname expression)")
	}

	varName = symNode.Variable.Name

	value, err := evalExpression(&tup.Expressions[1], rootObject)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate variable value: %w", err)
	}

	runtimeVariables[varName] = value

	return rootObject, nil
}
