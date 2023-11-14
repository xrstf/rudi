// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

func evalSymbol(ctx types.Context, sym *ast.Symbol) (types.Context, any, error) {
	switch {
	case sym.Variable != nil:
		varName := string(*sym.Variable)

		value, ok := ctx.GetVariable(varName)
		if !ok {
			return ctx, nil, fmt.Errorf("unknown variable %s (%T)", varName, varName)
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

func ptrTo[T any](s T) *T {
	return &s
}

func convertToAccessor(evaluated any) (*ast.EvaluatedPathStep, error) {
	switch asserted := evaluated.(type) {
	case ast.String:
		return &ast.EvaluatedPathStep{StringValue: ptrTo(string(asserted))}, nil
	case ast.Identifier:
		return &ast.EvaluatedPathStep{StringValue: ptrTo(string(asserted))}, nil
	case ast.Number:
		intVal, ok := asserted.ToInteger()
		if !ok {
			return nil, fmt.Errorf("cannot use floats as indices: %v", asserted.ToFloat())
		}

		return &ast.EvaluatedPathStep{IntegerValue: &intVal}, nil
	default:
		return nil, fmt.Errorf("cannot use %T in path expression", asserted)
	}
}

func traversePathExpression(ctx types.Context, value any, path *ast.PathExpression) (result any, pathErr error, traverseErr error) {
	evaluatedPath, err := evalPathExpression(ctx, path)
	if err != nil {
		return nil, err, nil
	}

	result, err = traverseEvaluatedPathExpression(ctx, value, evaluatedPath)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}

func evalPathExpression(ctx types.Context, path *ast.PathExpression) (*ast.EvaluatedPathExpression, error) {
	innerCtx := ctx
	result := &ast.EvaluatedPathExpression{
		Steps: []ast.EvaluatedPathStep{},
	}

	// The parsed path might just be "."; in this case it would still have 1 step in it,
	// because my peg syntax is wonky, but here we skip that step and just return an empty
	// result instead.
	if path.IsIdentity() {
		return result, nil
	}

	for _, step := range path.Steps {
		var (
			evaluated any
			err       error
		)

		// keep accumulating context changes, so you _could_ in theory do
		// $var[(set $bla 2)][(add $bla 2)] <-- would be $var[2][4]
		innerCtx, evaluated, err = evalNode(innerCtx, step)
		if err != nil {
			return nil, fmt.Errorf("invalid accessor: %w", err)
		}

		evaledAccessor, err := convertToAccessor(evaluated)
		if err != nil {
			return nil, err
		}

		result.Steps = append(result.Steps, *evaledAccessor)
	}

	return result, nil
}

func traverseEvaluatedPathExpression(ctx types.Context, value any, path *ast.EvaluatedPathExpression) (any, error) {
	var err error

	for _, step := range path.Steps {
		if valueAsVector, ok := value.(ast.Vector); ok {
			if step.IntegerValue == nil {
				return nil, fmt.Errorf("cannot use %v as an array index", step.String())
			}

			rawValue := valueAsVector.Data[*step.IntegerValue]
			value, err = types.WrapNative(rawValue)
			if err != nil {
				return nil, fmt.Errorf("cannot wrap %v (%T): %w", rawValue, rawValue, err)
			}

			continue
		}

		if valueAsObject, ok := value.(ast.Object); ok {
			if step.StringValue == nil {
				return nil, fmt.Errorf("cannot use %v as an object key", step.String())
			}

			rawValue := valueAsObject.Data[*step.StringValue]
			value, err = types.WrapNative(rawValue)
			if err != nil {
				return nil, fmt.Errorf("cannot wrap %v (%T): %w", rawValue, rawValue, err)
			}

			continue
		}

		return nil, fmt.Errorf("cannot descend with %s into %T", step.String(), value)
	}

	return value, nil
}
