// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func EvalSymbol(ctx types.Context, sym ast.Symbol) (types.Context, any, error) {
	// pre-evaluate the path expression
	pathExpr := ast.EvaluatedPathExpression{}

	if sym.PathExpression != nil {
		evaluated, err := EvalPathExpression(ctx, sym.PathExpression)
		if err != nil {
			return ctx, nil, fmt.Errorf("invalid path expression: %w", err)
		}
		pathExpr = *evaluated
	}

	return EvalSymbolWithEvaluatedPath(ctx, sym, pathExpr)
}

func EvalSymbolWithEvaluatedPath(ctx types.Context, sym ast.Symbol, path ast.EvaluatedPathExpression) (types.Context, any, error) {
	rootDoc := ctx.GetDocument()
	rootValue := rootDoc.Get()

	// sanity check
	if sym.Variable == nil && sym.PathExpression == nil {
		return ctx, nil, errors.New("invalid symbol")
	}

	// . always returns the root document
	if sym.IsDot() {
		wrapped, err := types.WrapNative(rootValue)
		if err != nil {
			return ctx, nil, err
		}

		return ctx, wrapped, nil
	}

	// if this symbol is a variable, replace the root value with the variable's value
	if sym.Variable != nil {
		var ok bool

		varName := string(*sym.Variable)

		rootValue, ok = ctx.GetVariable(varName)
		if !ok {
			return ctx, nil, fmt.Errorf("unknown variable %s", varName)
		}
	}

	deeper, err := TraverseEvaluatedPathExpression(ctx, rootValue, path)
	if err != nil {
		return ctx, nil, fmt.Errorf("cannot evaluate %s: %w", sym.String(), err)
	}

	return ctx, deeper, nil
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
			evaluated = ast.String(string(asserted))
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

func TraverseEvaluatedPathExpression(ctx types.Context, value any, path ast.EvaluatedPathExpression) (any, error) {
	if len(path.Steps) == 0 {
		return types.WrapNative(value)
	}

	for _, step := range path.Steps {
		unwrappedValue, err := types.UnwrapType(value)
		if err != nil {
			return nil, fmt.Errorf("cannot descend with %s into %T", step.String(), value)
		}

		if valueAsVector, ok := unwrappedValue.([]any); ok {
			if step.IntegerValue == nil {
				return nil, fmt.Errorf("cannot use %v as an array index", step.String())
			}

			index := int(*step.IntegerValue)
			if index < 0 || index >= len(valueAsVector) {
				return nil, fmt.Errorf("index %d out of bounds", index)
			}

			rawValue := valueAsVector[index]
			value, err = types.WrapNative(rawValue)
			if err != nil {
				return nil, err
			}

			continue
		}

		if valueAsObject, ok := unwrappedValue.(map[string]any); ok {
			if step.StringValue == nil {
				return nil, fmt.Errorf("cannot use %v as an object key", step.String())
			}

			rawValue, exists := valueAsObject[*step.StringValue]
			if !exists {
				return nil, fmt.Errorf("no such key: %q", *step.StringValue)
			}

			value, err = types.WrapNative(rawValue)
			if err != nil {
				return nil, err
			}

			continue
		}

		return nil, fmt.Errorf("cannot descend with %s into %T", step.String(), value)
	}

	return value, nil
}
