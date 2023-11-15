// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

func lenFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, list, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	if str, ok := list.(ast.String); ok {
		return ast.Number{Value: int64(len(str))}, nil
	}

	if vector, ok := list.(ast.Vector); ok {
		return ast.Number{Value: int64(len(vector.Data))}, nil
	}

	if obj, ok := list.(ast.Object); ok {
		return ast.Number{Value: int64(len(obj.Data))}, nil
	}

	return nil, errors.New("argument is neither a string, vector nor object")
}

func appendFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size < 2 {
		return nil, fmt.Errorf("expected 2+ arguments, got %d", size)
	}

	_, list, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	vector, ok := list.(ast.Vector)
	if !ok {
		return nil, fmt.Errorf("argument #0 is not a vector, but %T", list)
	}

	evaluated, err := evalArgs(ctx, args, 1)
	if err != nil {
		return nil, err
	}

	result := vector.Clone()
	result.Data = append(result.Data, evaluated...)

	return result, nil
}

func prependFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size < 2 {
		return nil, fmt.Errorf("expected 2+ arguments, got %d", size)
	}

	_, list, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	vector, ok := list.(ast.Vector)
	if !ok {
		return nil, fmt.Errorf("argument #0 is not a vector, but %T", list)
	}

	evaluated, err := evalArgs(ctx, args, 1)
	if err != nil {
		return nil, err
	}

	wrapped, err := types.WrapNative(evaluated)
	if err != nil {
		panic("failed to wrap a []any, this should never happen")
	}

	evaluatedVector, ok := wrapped.(ast.Vector)
	if !ok {
		return nil, fmt.Errorf("argument #0 is not a vector, but %T", list)
	}

	evaluatedVector.Data = append(evaluatedVector.Data, vector.Data...)

	return evaluatedVector, nil
}

func reverseFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, list, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	strVal, ok := list.(ast.String)
	if ok {
		// thank you https://stackoverflow.com/a/10030772
		result := []rune(strVal)
		for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
			result[i], result[j] = result[j], result[i]
		}

		return ast.String(result), nil
	}

	vector, ok := list.(ast.Vector)
	if ok {
		result := vector.Clone()
		for i, j := 0, len(result.Data)-1; i < j; i, j = i+1, j-1 {
			result.Data[i], result.Data[j] = result.Data[j], result.Data[i]
		}

		return result, nil
	}

	return nil, fmt.Errorf("argument is neither a vector nor a string, but %T", list)
}
