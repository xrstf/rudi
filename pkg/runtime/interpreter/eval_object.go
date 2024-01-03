// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package interpreter

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/pathexpr"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

func (i *interpreter) EvalObjectNode(ctx types.Context, obj ast.ObjectNode) (any, error) {
	result := map[string]any{}

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
				return nil, errors.New("cannot use bang modifier in object keys")
			}

			key = asserted.Name
		default:
			key, err = i.EvalExpression(ctx, pair.Key)
			if err != nil {
				return nil, fmt.Errorf("failed to evaluate object key %s: %w", pair.Key.String(), err)
			}
		}

		keyString, err := ctx.Coalesce().ToString(key)
		if err != nil {
			return nil, fmt.Errorf("object key: %w", err)
		}

		value, err = i.EvalExpression(ctx, pair.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate object value %s: %w", pair.Value.String(), err)
		}

		result[keyString] = value
	}

	deeper, err := pathexpr.Apply(ctx, result, obj.PathExpression)
	if err != nil {
		return nil, err
	}

	return deeper, nil
}
