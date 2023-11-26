// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

// evaluated objects are technically considered expressions.
func EvalObject(ctx types.Context, obj ast.Object) (types.Context, any, error) {
	return ctx, obj, nil
}

func EvalObjectNode(ctx types.Context, obj ast.ObjectNode) (types.Context, any, error) {
	innerCtx := ctx
	result := ast.Object{
		Data: map[string]any{},
	}

	var (
		key   any
		value any
		err   error
	)

	for _, pair := range obj.Data {
		switch asserted := pair.Key.(type) {
		// as a convenience feature, we allow unquoted object keys, which are parsed as bare identifiers
		case ast.Identifier:
			if asserted.Bang {
				return ctx, nil, errors.New("cannot use bang modifier in object keys")
			}

			key = asserted.Name
		default:
			// Just like with arrays, use a growing context during the object evaluation,
			// in case someone wants to define a variable here... for some reason.
			innerCtx, key, err = EvalExpression(innerCtx, pair.Key)
			if err != nil {
				return ctx, nil, fmt.Errorf("failed to evaluate object key %s: %w", pair.Key.String(), err)
			}
		}

		keyString, err := ctx.Coalesce().ToString(key)
		if err != nil {
			return ctx, nil, fmt.Errorf("object key: %w", err)
		}

		innerCtx, value, err = EvalExpression(innerCtx, pair.Value)
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to evaluate object value %s: %w", pair.Value.String(), err)
		}

		result.Data[keyString] = value
	}

	if obj.PathExpression != nil {
		deeper, err := TraversePathExpression(ctx, result, obj.PathExpression)
		if err != nil {
			return ctx, nil, err
		}

		return ctx, deeper, nil
	}

	return ctx, result, nil
}
