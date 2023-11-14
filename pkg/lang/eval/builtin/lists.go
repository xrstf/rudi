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
		return ast.Number{Value: len(str)}, nil
	}

	if vector, ok := list.(ast.Vector); ok {
		return ast.Number{Value: len(vector.Data)}, nil
	}

	if obj, ok := list.(ast.Object); ok {
		return ast.Number{Value: len(obj.Data)}, nil
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

	result := ast.Vector{}
	copy(result.Data, vector.Data)

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

	result := ast.Vector{}
	copy(result.Data, vector.Data)

	result.Data = append(result.Data, evaluated...)

	return result, nil
}
