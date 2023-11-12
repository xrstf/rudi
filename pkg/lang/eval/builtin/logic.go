// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/eval/coalescing"
	"go.xrstf.de/corel/pkg/lang/eval/types"
)

func andFunction(ctx types.Context, args []Argument) (any, error) {
	if size := len(args); size < 1 {
		return nil, fmt.Errorf("expected 1+ arguments, got %d", size)
	}

	evaluated, err := evalArgs(ctx, args, 0)
	if err != nil {
		return nil, err
	}

	result := true
	for i, arg := range evaluated {
		part, err := coalescing.ToBool(arg)
		if err != nil {
			return nil, fmt.Errorf("argument #%d not boolish: %w", i, err)
		}

		result = result && part
	}

	return result, nil
}

func orFunction(ctx types.Context, args []Argument) (any, error) {
	if size := len(args); size < 1 {
		return nil, fmt.Errorf("expected 1+ arguments, got %d", size)
	}

	evaluated, err := evalArgs(ctx, args, 0)
	if err != nil {
		return nil, err
	}

	result := false
	for i, arg := range evaluated {
		part, err := coalescing.ToBool(arg)
		if err != nil {
			return nil, fmt.Errorf("argument #%d not boolish: %w", i, err)
		}

		result = result || part
	}

	return result, nil
}

func notFunction(ctx types.Context, args []Argument) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, evaluated, err := args[0].Eval(ctx)
	if err != nil {
		return nil, err
	}

	arg, err := coalescing.ToBool(evaluated)
	if err != nil {
		return nil, fmt.Errorf("argument is not boolish: %w", err)
	}

	return !arg, nil
}
