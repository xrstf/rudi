// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/coalescing"
)

func evalSymbol(ctx Context, sym *ast.Symbol) (Context, interface{}, error) {
	switch {
	case sym.Variable != nil:
		varName := sym.Variable.Name

		value, ok := ctx.variables[varName]
		if !ok {
			return ctx, nil, fmt.Errorf("unknown variable %s", varName)
		}

		// descend into the variable value
		if expr := sym.PathExpression; expr != nil {
			deeper, pathErr, traverseErr := traversePathExpression(ctx, value, expr)
			if pathErr != nil {
				return ctx, nil, fmt.Errorf("invalid path expression %s: %w", expr.String(), pathErr)
			}
			if traverseErr != nil {
				return ctx, nil, fmt.Errorf("path %s not found: %w", expr.String(), traverseErr)
			}
			value = deeper
		}

		return ctx, value, nil

	case sym.PathExpression != nil:
		value, pathErr, traverseErr := traversePathExpression(ctx, ctx.document.Data, sym.PathExpression)
		if pathErr != nil {
			return ctx, nil, fmt.Errorf("invalid path expression %s: %w", sym.PathExpression.String(), pathErr)
		}
		if traverseErr != nil {
			return ctx, nil, fmt.Errorf("path %s not found: %w", sym.PathExpression.String(), traverseErr)
		}

		return ctx, value, nil
	}

	return ctx, nil, fmt.Errorf("unknown symbol %T (%s)", sym, sym.String())
}

func traversePathExpression(ctx Context, value interface{}, path *ast.PathExpression) (result interface{}, pathErr error, traverseErr error) {
	innerCtx := ctx

	for _, accessor := range path.Steps {
		var step interface{}

		switch {
		case accessor.Identifier != nil:
			step = accessor.Identifier.Name
		case accessor.Integer != nil:
			step = *accessor.Integer
		case accessor.StringNode != nil:
			step = accessor.StringNode.Value
		case accessor.Variable != nil:
			value, ok := innerCtx.variables[accessor.Variable.Name]
			if !ok {
				return nil, fmt.Errorf("unknown variable %s", accessor.Variable.Name), nil
			}
			step = value
		case accessor.Tuple != nil:
			var (
				value interface{}
				err   error
			)

			// keep accumulating context changes, so you _could_ in theory do
			// $var[(set $bla 2)][(add $bla 2)] <-- would be $var[2][4]
			innerCtx, value, err = evalTuple(innerCtx, accessor.Tuple)
			if err != nil {
				return nil, fmt.Errorf("invalid accessor: %w", err), nil
			}

			step = value
		}

		if valueAsSlice, ok := value.([]interface{}); ok {
			stepInt, err := coalescing.ToInt64(step)
			if err != nil {
				return nil, nil, fmt.Errorf("cannot use %v (%T) as an array index", step, step)
			}

			value = valueAsSlice[stepInt]
			continue
		}

		if valueAsMap, ok := value.(map[string]interface{}); ok {
			stepString, err := coalescing.ToString(step)
			if err != nil {
				return nil, nil, fmt.Errorf("cannot use %v (%T) as an object key", step, step)
			}

			value = valueAsMap[stepString]
			continue
		}

		return nil, nil, fmt.Errorf("cannot descend into %T", value)
	}

	return value, nil, nil
}
