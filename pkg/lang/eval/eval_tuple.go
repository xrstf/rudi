package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/builtin"
)

func evalTuple(tup *ast.Tuple, rootObject *Object) (interface{}, error) {
	funcName := tup.Identifier.Name

	// hardcode root behaviour for those tuples where not all
	// expressions can be pre-computed (in case, for example,
	// the else-path of an if statement would have side effects)
	switch funcName {
	case "if":
		return evalIfTuple(tup, rootObject)
	}

	function, ok := builtin.Functions[funcName]
	if !ok {
		return nil, fmt.Errorf("unknown function %s", funcName)
	}

	args := make([]interface{}, len(tup.Expressions))

	for i, expr := range tup.Expressions {
		arg, err := evalExpression(&expr, rootObject)
		if err != nil {
			return nil, fmt.Errorf("invalid argument %d: %w", i, err)
		}

		args[i] = arg
	}

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

	success, err := coalesceBool(condition)
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

func coalesceBool(val interface{}) (bool, error) {
	switch v := val.(type) {
	case bool:
		return v, nil
	case int64:
		return v != 0, nil
	case float64:
		return v != 0, nil
	case string:
		return len(v) > 0, nil
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("cannot coalesce %T into bool", val)
	}
}
