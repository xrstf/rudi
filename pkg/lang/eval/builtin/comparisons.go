// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/equality"
	"go.xrstf.de/otto/pkg/lang/eval"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

func eqFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 2 {
		return nil, fmt.Errorf("expected exactly 2 arguments, got %d", size)
	}

	_, leftData, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	leftValue, ok := leftData.(ast.Literal)
	if !ok {
		return nil, fmt.Errorf("argument #0 is not a literal, but %T", leftData)
	}

	_, rightData, err := eval.EvalExpression(ctx, args[1])
	if err != nil {
		return nil, fmt.Errorf("argument #1: %w", err)
	}

	rightValue, ok := rightData.(ast.Literal)
	if !ok {
		return nil, fmt.Errorf("argument #1 is not a literal, but %T", rightData)
	}

	equal, err := equality.StrictEqual(leftValue, rightValue)
	if err != nil {
		return nil, err
	}

	return ast.Bool(equal), nil
}

func likeFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 2 {
		return nil, fmt.Errorf("expected exactly 2 arguments, got %d", size)
	}

	_, leftData, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	leftValue, ok := leftData.(ast.Literal)
	if !ok {
		return nil, fmt.Errorf("argument #0 is not a literal, but %T", leftData)
	}

	_, rightData, err := eval.EvalExpression(ctx, args[1])
	if err != nil {
		return nil, fmt.Errorf("argument #1: %w", err)
	}

	rightValue, ok := rightData.(ast.Literal)
	if !ok {
		return nil, fmt.Errorf("argument #1 is not a literal, but %T", rightData)
	}

	equal, err := equality.EqualEnough(leftValue, rightValue)
	if err != nil {
		return nil, err
	}

	return ast.Bool(equal), nil
}
