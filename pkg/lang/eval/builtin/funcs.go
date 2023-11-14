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
	"empty?":  stateless(isEmptyFunction),

	// math
	"+": stateless(sumFunction),
	"-": stateless(minusFunction),
	"*": stateless(multiplyFunction),
	"/": stateless(divideFunction),

	// strings
	// "len": stateless(lenFunction) is defined for lists, but works for strings as well
	"concat":      stateless(concatFunction),
	"split":       stateless(splitFunction),
	"trim-prefix": stateless(trimPrefixFunction),
	"trim-suffix": stateless(trimSuffixFunction),
	"to-lower":    stateless(toLowerFunction),
	"to-upper":    stateless(toUpperFunction),

	// lists
	"len":     stateless(lenFunction),
	"append":  stateless(appendFunction),
	"prepend": stateless(prependFunction),

	// logic
	"and": stateless(andFunction),
	"or":  stateless(orFunction),
	"not": stateless(notFunction),

	// comparisons
	"eq": stateless(eqFunction),

	// types
	"type-of":   stateless(typeOfFunction),
	"to-string": stateless(toStringFunction),
	"to-int":    stateless(toIntFunction),
	"to-float":  stateless(toFloatFunction),
	"to-bool":   stateless(toBoolFunction),

	// hashes
	"sha1":   stateless(sha1Function),
	"sha256": stateless(sha256Function),
	"sha512": stateless(sha512Function),

	// encoding
	"to-base64":   stateless(toBase64Function),
	"from-base64": stateless(fromBase64Function),
}
