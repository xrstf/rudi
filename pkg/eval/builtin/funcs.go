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
	"do":      types.BasicFunction(doFunction, "eval a sequence of statements where only one expression is valid"),
	"has?":    types.BasicFunction(hasFunction, "returns true if the given symbol's path expression points to an existing value"),
	"default": types.BasicFunction(defaultFunction, "returns the default value if the first argument is empty"),
	"try":     types.BasicFunction(tryFunction, "returns the fallback if the first expression errors out"),
	"set":     types.BasicFunction(setFunction, "set a value in a variable/document, only really useful with ! modifier (set!)"),
	"delete":  deleteFunction{},
	"empty?":  types.BasicFunction(isEmptyFunction, "returns true when the given value is empty-ish (0, false, null, \"\", ...)"),

	// math
	"+": types.BasicFunction(sumFunction, "returns the sum of all of its arguments"),
	"-": types.BasicFunction(subFunction, "returns arg1 - arg2 - .. - argN"),
	"*": types.BasicFunction(multiplyFunction, "returns the product of all of its arguments"),
	"/": types.BasicFunction(divideFunction, "returns arg1 / arg2 / .. / argN"),

	// math aliases to make bang functions nicer (sum! vs +!)
	"add":  types.BasicFunction(sumFunction, "alias for +"),
	"sub":  types.BasicFunction(subFunction, "alias for -"),
	"mult": types.BasicFunction(multiplyFunction, "alias for *"),
	"div":  types.BasicFunction(divideFunction, "alias for div"),

	// strings
	// "len": lenFunction is defined for lists, but works for strings as well
	// "reverse" also works for strings
	"concat":      types.BasicFunction(concatFunction, "concatenate items in a vector usig a common glue string"),
	"split":       fromStringFunc(splitFunction, 2, "split a string into a vector"),
	"has-prefix?": fromStringFunc(hasPrefixFunction, 2, "returns true if the given string has the prefix"),
	"has-suffix?": fromStringFunc(hasSuffixFunction, 2, "returns true if the given string has the suffix"),
	"trim-prefix": fromStringFunc(trimPrefixFunction, 2, "removes the prefix from the string, if it exists"),
	"trim-suffix": fromStringFunc(trimSuffixFunction, 2, "removes the suffix from the string, if it exists"),
	"to-lower":    fromStringFunc(toLowerFunction, 1, "returns the lowercased version of the given string"),
	"to-upper":    fromStringFunc(toUpperFunction, 1, "returns the uppercased version of the given string"),
	"trim":        fromStringFunc(trimFunction, 1, "returns the given whitespace with leading/trailing whitespace removed"),

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
