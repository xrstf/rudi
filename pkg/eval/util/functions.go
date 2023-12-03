// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package util

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

// RawFunction is a function that receives its raw, unevaluated child expressions as arguments.
// This is the lowest level a function can be, allowing to selectively evaluate the arguments to
// control side effects.
type RawFunction func(ctx types.Context, args []ast.Expression) (any, error)

// LiteralFunction is a function that receives all of its arguments already evaluated, but not yet
// coalesced into specific types.
type LiteralFunction func(ctx types.Context, args []any) (any, error)

// Function is a Rudi function with additional options to make custom things less redundant.
type Function interface {
	types.Function

	MinArgs(min int) Function
	MaxArgs(max int) Function
}

type genericFunc struct {
	fun     RawFunction
	minArgs int
	maxArgs int
	desc    string
}

var _ Function = &genericFunc{}

func NewRawFunction(f RawFunction, description string) Function {
	return &genericFunc{
		fun:     f,
		desc:    description,
		minArgs: -1,
		maxArgs: -1,
	}
}

func NewLiteralFunction(fun LiteralFunction, description string) Function {
	return NewRawFunction(func(ctx types.Context, args []ast.Expression) (any, error) {
		values := make([]any, len(args))
		for i, arg := range args {
			_, evaluated, err := eval.EvalExpression(ctx, arg)
			if err != nil {
				return nil, fmt.Errorf("argument #%d: %w", i, err)
			}

			values[i] = evaluated
		}

		return fun(ctx, values)
	}, description)
}

func (f *genericFunc) Evaluate(ctx types.Context, args []ast.Expression) (any, error) {
	if err := f.checkSignature(args); err != nil {
		return nil, err
	}

	return f.fun(ctx, args)
}

func (f *genericFunc) Description() string {
	return f.desc
}

func (f *genericFunc) MaxArgs(max int) Function {
	if max < 0 {
		f.maxArgs = -1
	} else {
		f.maxArgs = max
	}

	return f
}

func (f *genericFunc) MinArgs(min int) Function {
	if min < 0 {
		f.minArgs = -1
	} else {
		f.minArgs = min
	}

	return f
}

func arguments(num int) string {
	if num == 1 {
		return "argument"
	}

	return "arguments"
}

func (f *genericFunc) checkSignature(args []ast.Expression) error {
	if f.minArgs < 0 && f.maxArgs < 0 {
		return nil
	}

	num := len(args)

	if f.minArgs == f.maxArgs {
		if num != f.minArgs {
			return fmt.Errorf("expected %d %s, got %d", f.minArgs, arguments(f.minArgs), num)
		}

		return nil
	}

	if f.minArgs < 0 {
		if num > f.maxArgs {
			return fmt.Errorf("expected up to %d %s, got %d", f.maxArgs, arguments(f.maxArgs), num)
		}
	} else if f.maxArgs < 0 {
		if num < f.minArgs {
			return fmt.Errorf("expected at most %d %s, got %d", f.minArgs, arguments(f.minArgs), num)
		}
	} else if num < f.minArgs || num > f.maxArgs {
		return fmt.Errorf("expected %d to %d %s, got %d", f.minArgs, f.maxArgs, arguments(f.maxArgs), num)
	}

	return nil
}
