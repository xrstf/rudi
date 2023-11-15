// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

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
	"has?":    stateless(hasFunction),
	"default": stateless(defaultFunction),
	"try":     stateless(tryFunction),
	"range":   stateless(rangeFunction),
	"set":     setFunction,
	"empty?":  stateless(isEmptyFunction),

	// math
	"+": stateless(sumFunction),
	"-": stateless(minusFunction),
	"*": stateless(multiplyFunction),
	"/": stateless(divideFunction),

	// strings
	// "len": stateless(lenFunction) is defined for lists, but works for strings as well
	// "reverse" also works for strings
	"concat":      stateless(concatFunction),
	"split":       stateless(fromStringFunc(splitFunction, 2)),
	"has-prefix?": stateless(fromStringFunc(hasPrefixFunction, 2)),
	"has-suffix?": stateless(fromStringFunc(hasSuffixFunction, 2)),
	"trim-prefix": stateless(fromStringFunc(trimPrefixFunction, 2)),
	"trim-suffix": stateless(fromStringFunc(trimSuffixFunction, 2)),
	"to-lower":    stateless(fromStringFunc(toLowerFunction, 1)),
	"to-upper":    stateless(fromStringFunc(toUpperFunction, 1)),

	// lists
	"len":     stateless(lenFunction),
	"append":  stateless(appendFunction),
	"prepend": stateless(prependFunction),
	"reverse": stateless(reverseFunction),

	// logic
	"and": stateless(andFunction),
	"or":  stateless(orFunction),
	"not": stateless(notFunction),

	// comparisons
	"eq?":   stateless(eqFunction),
	"like?": stateless(likeFunction),

	"lt?": stateless(makeNumberComparatorFunc(
		func(a, b int64) (any, error) { return a < b, nil },
		func(a, b float64) (any, error) { return a < b, nil },
	)),
	"lte?": stateless(makeNumberComparatorFunc(
		func(a, b int64) (any, error) { return a <= b, nil },
		func(a, b float64) (any, error) { return a <= b, nil },
	)),
	"gt?": stateless(makeNumberComparatorFunc(
		func(a, b int64) (any, error) { return a > b, nil },
		func(a, b float64) (any, error) { return a > b, nil },
	)),
	"gte?": stateless(makeNumberComparatorFunc(
		func(a, b int64) (any, error) { return a >= b, nil },
		func(a, b float64) (any, error) { return a >= b, nil },
	)),

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
