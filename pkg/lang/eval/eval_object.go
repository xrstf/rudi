// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalObject(ctx Context, obj *ast.Object) (Context, interface{}, error) {
	innerCtx := ctx
	result := map[string]interface{}{}

	var (
		key   interface{}
		value interface{}
		err   error
	)

	for _, pair := range obj.Data {
		// Just like with arrays, use a growing context during the object evaluation,
		// in case someone wants to define a variable here... for some reason.
		innerCtx, key, err = evalSymbol(innerCtx, &pair.Key)
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to evaluate object key %s: %w", pair.Key.String(), err)
		}

		keyString, ok := key.(string)
		if !ok {
			ident, ok := key.(*ast.Identifier)
			if !ok {
				return ctx, nil, fmt.Errorf("object key must be string or identifier, but got %T", key)
			}

			keyString = ident.Name
		}

		innerCtx, value, err = evalExpression(innerCtx, &pair.Value)
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to evaluate object value %s: %w", pair.Value.String(), err)
		}

		result[keyString] = value
	}

	return ctx, result, nil
}
