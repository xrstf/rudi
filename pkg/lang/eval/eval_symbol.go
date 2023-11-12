// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/coalescing"
	"go.xrstf.de/corel/pkg/lang/eval/types"
)

func evalSymbol(ctx types.Context, sym *ast.Symbol) (types.Context, any, error) {
	switch {
	case sym.Variable != nil:
		varName := sym.Variable.Name

		value, ok := ctx.GetVariable(varName)
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
		rootDoc := ctx.GetDocument()
		value, pathErr, traverseErr := traversePathExpression(ctx, rootDoc.Get(), sym.PathExpression)
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

func traversePathExpression(ctx types.Context, value any, path *ast.PathExpression) (result any, pathErr error, traverseErr error) {
	innerCtx := ctx

	for _, accessor := range path.Steps {
		var step any

		switch {
		case accessor.Identifier != nil:
			step = ast.String{Value: accessor.Identifier.Name}
		case accessor.Integer != nil:
			step = ast.Number{Value: *accessor.Integer}
		case accessor.StringNode != nil:
			step = ast.String{Value: accessor.StringNode.Value}
		case accessor.Variable != nil:
			value, ok := innerCtx.GetVariable(accessor.Variable.Name)
			if !ok {
				return nil, fmt.Errorf("unknown variable %s", accessor.Variable.Name), nil
			}
			step = value
		case accessor.Tuple != nil:
			var (
				value any
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

		if valueAsVector, ok := value.(ast.Vector); ok {
			stepInt, err := coalescing.ToInt64(step)
			if err != nil {
				return nil, nil, fmt.Errorf("cannot use %v (%T) as an array index: %w", step, step, err)
			}

			rawValue := valueAsVector.Data[stepInt]
			value, err = types.WrapNative(rawValue)
			if err != nil {
				return nil, nil, fmt.Errorf("cannot wrap %v (%T): %w", rawValue, rawValue, err)
			}

			continue
		}

		if valueAsObject, ok := value.(ast.Object); ok {
			stepString, err := coalescing.ToString(step)
			if err != nil {
				return nil, nil, fmt.Errorf("cannot use %v (%T) as an object key: %w", step, step, err)
			}

			rawValue := valueAsObject.Data[stepString]
			value, err = types.WrapNative(rawValue)
			if err != nil {
				return nil, nil, fmt.Errorf("cannot wrap %v (%T): %w", rawValue, rawValue, err)
			}

			continue
		}

		return nil, nil, fmt.Errorf("cannot descend into %T", value)
	}

	return value, nil, nil
}
