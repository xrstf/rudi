// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

var Functions = types.Functions{
	// core
	"if":      ifFunction,
	"do":      doFunction,
	"has?":    hasFunction,
	"default": defaultFunction,
	"try":     tryFunction,
	"set":     setFunction,
	"delete":  deleteFunction,
	"empty?":  isEmptyFunction,

	// math
	"+": sumFunction,
	"-": subFunction,
	"*": multiplyFunction,
	"/": divideFunction,

	// math aliases to make bang functions nicer (sum! vs +!)
	"add":  sumFunction,
	"sub":  subFunction,
	"mult": multiplyFunction,
	"div":  divideFunction,

	// strings
	// "len": lenFunction is defined for lists, but works for strings as well
	// "reverse" also works for strings
	"concat":      concatFunction,
	"split":       fromStringFunc(splitFunction, 2),
	"has-prefix?": fromStringFunc(hasPrefixFunction, 2),
	"has-suffix?": fromStringFunc(hasSuffixFunction, 2),
	"trim-prefix": fromStringFunc(trimPrefixFunction, 2),
	"trim-suffix": fromStringFunc(trimSuffixFunction, 2),
	"to-lower":    fromStringFunc(toLowerFunction, 1),
	"to-upper":    fromStringFunc(toUpperFunction, 1),

	// lists
	"len":     lenFunction,
	"append":  appendFunction,
	"prepend": prependFunction,
	"reverse": reverseFunction,
	"range":   rangeFunction,
	"map":     mapFunction,
	"filter":  filterFunction,

	// logic
	"and": andFunction,
	"or":  orFunction,
	"not": notFunction,

	// comparisons
	"eq?":   eqFunction,
	"like?": likeFunction,

	"lt?": makeNumberComparatorFunc(
		func(a, b int64) (ast.Bool, error) { return ast.Bool(a < b), nil },
		func(a, b float64) (ast.Bool, error) { return ast.Bool(a < b), nil },
	),
	"lte?": makeNumberComparatorFunc(
		func(a, b int64) (ast.Bool, error) { return ast.Bool(a <= b), nil },
		func(a, b float64) (ast.Bool, error) { return ast.Bool(a <= b), nil },
	),
	"gt?": makeNumberComparatorFunc(
		func(a, b int64) (ast.Bool, error) { return ast.Bool(a > b), nil },
		func(a, b float64) (ast.Bool, error) { return ast.Bool(a > b), nil },
	),
	"gte?": makeNumberComparatorFunc(
		func(a, b int64) (ast.Bool, error) { return ast.Bool(a >= b), nil },
		func(a, b float64) (ast.Bool, error) { return ast.Bool(a >= b), nil },
	),

	// types
	"type-of":   typeOfFunction,
	"to-string": toStringFunction,
	"to-int":    toIntFunction,
	"to-float":  toFloatFunction,
	"to-bool":   toBoolFunction,

	// hashes
	"sha1":   sha1Function,
	"sha256": sha256Function,
	"sha512": sha512Function,

	// encoding
	"to-base64":   toBase64Function,
	"from-base64": fromBase64Function,

	// dates & time
	"now": nowFunction,
}
