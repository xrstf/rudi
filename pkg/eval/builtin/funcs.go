// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

var Functions = types.Functions{
	// core
	"if":      types.BasicFunction(ifFunction, "evaluate one of two expressions based on a condition"),
	"do":      types.BasicFunction(doFunction, ""),
	"has?":    types.BasicFunction(hasFunction, ""),
	"default": types.BasicFunction(defaultFunction, ""),
	"try":     types.BasicFunction(tryFunction, ""),
	"set":     types.BasicFunction(setFunction, ""),
	"delete":  deleteFunction{},
	"empty?":  types.BasicFunction(isEmptyFunction, ""),

	// math
	"+": types.BasicFunction(sumFunction, ""),
	"-": types.BasicFunction(subFunction, ""),
	"*": types.BasicFunction(multiplyFunction, ""),
	"/": types.BasicFunction(divideFunction, ""),

	// math aliases to make bang functions nicer (sum! vs +!)
	"add":  types.BasicFunction(sumFunction, ""),
	"sub":  types.BasicFunction(subFunction, ""),
	"mult": types.BasicFunction(multiplyFunction, ""),
	"div":  types.BasicFunction(divideFunction, ""),

	// strings
	// "len": lenFunction is defined for lists, but works for strings as well
	// "reverse" also works for strings
	"concat":      types.BasicFunction(concatFunction, ""),
	"split":       fromStringFunc(splitFunction, 2),
	"has-prefix?": fromStringFunc(hasPrefixFunction, 2),
	"has-suffix?": fromStringFunc(hasSuffixFunction, 2),
	"trim-prefix": fromStringFunc(trimPrefixFunction, 2),
	"trim-suffix": fromStringFunc(trimSuffixFunction, 2),
	"to-lower":    fromStringFunc(toLowerFunction, 1),
	"to-upper":    fromStringFunc(toUpperFunction, 1),
	"trim":        fromStringFunc(trimFunction, 1),

	// lists
	"len":       types.BasicFunction(lenFunction, ""),
	"append":    types.BasicFunction(appendFunction, ""),
	"prepend":   types.BasicFunction(prependFunction, ""),
	"reverse":   types.BasicFunction(reverseFunction, ""),
	"range":     types.BasicFunction(rangeFunction, ""),
	"map":       types.BasicFunction(mapFunction, ""),
	"filter":    types.BasicFunction(filterFunction, ""),
	"contains?": types.BasicFunction(containsFunction, ""),

	// logic
	"and": types.BasicFunction(andFunction, ""),
	"or":  types.BasicFunction(orFunction, ""),
	"not": types.BasicFunction(notFunction, ""),

	// comparisons
	"eq?":   types.BasicFunction(eqFunction, ""),
	"like?": types.BasicFunction(likeFunction, ""),

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
	"type-of":   types.BasicFunction(typeOfFunction, ""),
	"to-string": types.BasicFunction(toStringFunction, ""),
	"to-int":    types.BasicFunction(toIntFunction, ""),
	"to-float":  types.BasicFunction(toFloatFunction, ""),
	"to-bool":   types.BasicFunction(toBoolFunction, ""),

	// hashes
	"sha1":   types.BasicFunction(sha1Function, ""),
	"sha256": types.BasicFunction(sha256Function, ""),
	"sha512": types.BasicFunction(sha512Function, ""),

	// encoding
	"to-base64":   types.BasicFunction(toBase64Function, ""),
	"from-base64": types.BasicFunction(fromBase64Function, ""),

	// dates & time
	"now": types.BasicFunction(nowFunction, ""),
}
