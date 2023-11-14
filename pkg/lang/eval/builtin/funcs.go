// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

func evalArgs(ctx types.Context, args []ast.Expression, argShift int) ([]any, error) {
	values := make([]any, len(args)-argShift)
	for i, arg := range args[argShift:] {
		_, evaluated, err := eval.EvalExpression(ctx, arg)
		if err != nil {
			return nil, fmt.Errorf("argument #%d: %w", i+argShift, err)
		}

		values[i] = evaluated
	}

	return values, nil
}

type StatelessFunc func(ctx types.Context, args []ast.Expression) (any, error)

func stateless(f StatelessFunc) types.Function {
	return func(ctx types.Context, args []ast.Expression) (types.Context, any, error) {
		result, err := f(ctx, args)
		return ctx, result, err
	}
}

var Functions = types.Functions{
	// core
	"if":      stateless(ifFunction),
	"do":      stateless(doFunction),
	"has":     stateless(hasFunction),
	"default": stateless(defaultFunction),
	"try":     stateless(tryFunction),
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
