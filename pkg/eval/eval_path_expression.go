// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/pathexpr"
)

func EvalPathExpression(ctx types.Context, path *ast.PathExpression) (*ast.EvaluatedPathExpression, error) {
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
		switch asserted := step.(type) {
		case ast.Identifier:
			evaluated = asserted
		default:
			innerCtx, evaluated, err = EvalExpression(innerCtx, step)
			if err != nil {
				return nil, fmt.Errorf("invalid accessor: %w", err)
			}
		}

		evaledAccessor, err := convertToAccessor(evaluated)
		if err != nil {
			return nil, err
		}

		result.Steps = append(result.Steps, *evaledAccessor)
	}

	return result, nil
}

func ptrTo[T any](s T) *T {
	return &s
}

func convertToAccessor(evaluated any) (*ast.EvaluatedPathStep, error) {
	switch asserted := evaluated.(type) {
	case ast.String:
		return &ast.EvaluatedPathStep{StringValue: ptrTo(string(asserted))}, nil
	case ast.Identifier:
		if asserted.Bang {
			return nil, errors.New("cannot use bang modifier in path expression")
		}

		return &ast.EvaluatedPathStep{StringValue: &asserted.Name}, nil
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

func TraverseEvaluatedPathExpression(ctx types.Context, value any, path ast.EvaluatedPathExpression) (any, error) {
	result, err := pathexpr.Get(value, pathexpr.FromEvaluatedPath(path))
	if err != nil {
		return nil, err
	}

	return types.WrapNative(result)
}
