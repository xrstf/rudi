// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func andFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size < 1 {
		return nil, fmt.Errorf("expected 1+ arguments, got %d", size)
	}

	evaluated, err := evalArgs(ctx, args, 0)
	if err != nil {
		return nil, err
	}

	result := true
	for i, arg := range evaluated {
		part, ok := arg.(ast.Bool)
		if !ok {
			return nil, fmt.Errorf("argument #%d is not a bool, but %T", i, arg)
		}

		result = result && bool(part)
	}

	return ast.Bool(result), nil
}

func orFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size < 1 {
		return nil, fmt.Errorf("expected 1+ arguments, got %d", size)
	}

	evaluated, err := evalArgs(ctx, args, 0)
	if err != nil {
		return nil, err
	}

	result := false
	for i, arg := range evaluated {
		part, ok := arg.(ast.Bool)
		if !ok {
			return nil, fmt.Errorf("argument #%d is not a bool, but %T", i, arg)
		}

		result = result || bool(part)
	}

	return ast.Bool(result), nil
}

func notFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, evaluated, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	arg, ok := evaluated.(ast.Bool)
	if !ok {
		return nil, fmt.Errorf("argument is not bool, but %T", evaluated)
	}

	return ast.Bool(!arg), nil
}
