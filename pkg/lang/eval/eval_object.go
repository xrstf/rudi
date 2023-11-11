package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalObject(obj *ast.Object, rootObject *Object) (interface{}, error) {
	result := map[string]interface{}{}

	for _, pair := range obj.Data {
		key, err := evalSymbol(&pair.Key, rootObject)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate object key %s: %w", pair.Key.String(), err)
		}

		keyString, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("object key must be string, but got %T", key)
		}

		value, err := evalExpression(&pair.Value, rootObject)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate object value %s: %w", pair.Value.String(), err)
		}

		result[keyString] = value
	}

	return result, nil
}
