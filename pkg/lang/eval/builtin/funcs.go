// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

type Argument interface {
	Eval(ctx types.Context) (types.Context, any, error)
	String() string
	Node() ast.Node
}

func evalArgs(ctx types.Context, args []Argument, argShift int) ([]any, error) {
	values := make([]any, len(args)-argShift)
	for i, arg := range args[argShift:] {
		_, evaluated, err := arg.Eval(ctx)
		if err != nil {
			return nil, fmt.Errorf("argument #%d: %w", i+argShift, err)
		}

		values[i] = evaluated
	}

	return values, nil
}

type GenericFunc func(ctx types.Context, args []Argument) (types.Context, any, error)
type StatelessFunc func(ctx types.Context, args []Argument) (any, error)

func stateless(f StatelessFunc) GenericFunc {
	return func(ctx types.Context, args []Argument) (types.Context, any, error) {
		result, err := f(ctx, args)
		return ctx, result, err
	}
}

var Functions = map[string]GenericFunc{
	// core
	"if":      stateless(ifFunction),
	"do":      stateless(doFunction),
	"has":     stateless(hasFunction),
	"default": stateless(defaultFunction),
	"set":     setFunction,

	// math
	"+": stateless(sumFunction),
	"-": stateless(minusFunction),
	"*": stateless(multiplyFunction),
	"/": stateless(divideFunction),

	// strings
	"concat": stateless(concatFunction),
	"split":  stateless(splitFunction),

	// lists
	"len": stateless(lenFunction),

	// logic
	"and": stateless(andFunction),
	"or":  stateless(orFunction),
	"not": stateless(notFunction),

	// comparisons
	"eq": stateless(eqFunction),
}
