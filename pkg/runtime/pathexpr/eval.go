// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package pathexpr

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
	genericpathexpr "go.xrstf.de/rudi/pkg/pathexpr"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

func Apply(ctx types.Context, value any, path *ast.PathExpression) (any, error) {
	if path == nil {
		return value, nil
	}

	evaluated, err := Eval(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("invalid path expression: %w", err)
	}

	return Traverse(value, *evaluated)
}

func Traverse(value any, path ast.EvaluatedPathExpression) (any, error) {
	return genericpathexpr.Get(value, genericpathexpr.FromEvaluatedPath(path))
}

func Eval(ctx types.Context, path *ast.PathExpression) (*ast.EvaluatedPathExpression, error) {
	result := &ast.EvaluatedPathExpression{
		Steps: []ast.EvaluatedPathStep{},
	}

	// The parsed path might just be "."; in this case it would still have 1 step in it,
	// because my peg syntax is wonky, but here we skip that step and just return an empty
	// result instead.
	if path.IsIdentity() {
		return result, nil
	}

	innerCtx := ctx
	runtime := ctx.Runtime()

	for _, step := range path.Steps {
		var (
			evaluated any
			err       error
		)

		// keep accumulating context changes, so you _could_ in theory do
		// $var[(set $bla 2)][(add $bla 2)] <-- would be $var[2][4]
		switch asserted := step.(type) {
		case ast.Identifier:
			evaluated = asserted.Name
		default:
			innerCtx, evaluated, err = runtime.EvalExpression(innerCtx, step)
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
	case string:
		return &ast.EvaluatedPathStep{StringValue: &asserted}, nil
	case int:
		return &ast.EvaluatedPathStep{IntegerValue: ptrTo(int64(asserted))}, nil
	case int32:
		return &ast.EvaluatedPathStep{IntegerValue: ptrTo(int64(asserted))}, nil
	case int64:
		return &ast.EvaluatedPathStep{IntegerValue: &asserted}, nil
	default:
		return nil, fmt.Errorf("cannot use %T in path expression", asserted)
	}
}
