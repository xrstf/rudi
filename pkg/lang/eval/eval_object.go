// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/coalescing"
	"go.xrstf.de/corel/pkg/lang/eval/types"
)

func evalObject(ctx types.Context, obj *ast.ObjectNode) (types.Context, interface{}, error) {
	innerCtx := ctx
	result := map[string]interface{}{}

	var (
		key   interface{}
		value interface{}
		err   error
	)

	for _, pair := range obj.Data {
		// as a convenience feature, we allow unquoted object keys, which are parsed as bare identifiers
		if pair.Key.IdentifierNode != nil {
			key = pair.Key.IdentifierNode.Name
		} else if pair.Key.ObjectNode != nil {
			return ctx, nil, fmt.Errorf("cannot handle object keys of type %T", pair.Key.ObjectNode)
		} else if pair.Key.VectorNode != nil {
			return ctx, nil, fmt.Errorf("cannot handle object keys of type %T", pair.Key.VectorNode)
		} else {
			// Just like with arrays, use a growing context during the object evaluation,
			// in case someone wants to define a variable here... for some reason.
			innerCtx, key, err = evalExpression(innerCtx, &pair.Key)
			if err != nil {
				return ctx, nil, fmt.Errorf("failed to evaluate object key %s: %w", pair.Key.String(), err)
			}
		}

		keyString, err := coalescing.ToString(key)
		if err != nil {
			return ctx, nil, fmt.Errorf("object key must be stringish, but got %T", key)
		}

		innerCtx, value, err = evalExpression(innerCtx, &pair.Value)
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to evaluate object value %s: %w", pair.Value.String(), err)
		}

		result[keyString] = value
	}

	return ctx, result, nil
}
