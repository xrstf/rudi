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
	"concat":      types.BasicFunction(concatFunction, "concatenate items in a vector using a common glue string"),
	"split":       fromStringFunc(splitFunction, 2, "split a string into a vector"),
	"has-prefix?": fromStringFunc(hasPrefixFunction, 2, "returns true if the given string has the prefix"),
	"has-suffix?": fromStringFunc(hasSuffixFunction, 2, "returns true if the given string has the suffix"),
	"trim-prefix": fromStringFunc(trimPrefixFunction, 2, "removes the prefix from the string, if it exists"),
	"trim-suffix": fromStringFunc(trimSuffixFunction, 2, "removes the suffix from the string, if it exists"),
	"to-lower":    fromStringFunc(toLowerFunction, 1, "returns the lowercased version of the given string"),
	"to-upper":    fromStringFunc(toUpperFunction, 1, "returns the uppercased version of the given string"),
	"trim":        fromStringFunc(trimFunction, 1, "returns the given whitespace with leading/trailing whitespace removed"),

	// lists
	"len":       types.BasicFunction(lenFunction, "returns the length of a string, vector or object"),
	"append":    types.BasicFunction(appendFunction, "appends more strings to a string or arbitrary items into a vector"),
	"prepend":   types.BasicFunction(prependFunction, "prepends more strings to a string or arbitrary items into a vector"),
	"reverse":   types.BasicFunction(reverseFunction, "reverses a string or the elements of a vector"),
	"range":     types.BasicFunction(rangeFunction, "allows to iterate (loop) over a vector or object"),
	"map":       types.BasicFunction(mapFunction, "applies an expression to every element in a vector or object"),
	"filter":    types.BasicFunction(filterFunction, "returns a copy of a given vector/object with only those elements remaining that satisfy a condition"),
	"contains?": types.BasicFunction(containsFunction, "returns true if a string contains a substring or a vector contains the given element"),

	// logic
	"and": types.BasicFunction(andFunction, "returns true if all arguments are true"),
	"or":  types.BasicFunction(orFunction, "returns true if any of the arguments is true"),
	"not": types.BasicFunction(notFunction, "negates the given argument"),

	// comparisons
	"eq?":   types.BasicFunction(eqFunction, "equality check: return true if both arguments are the same"),
	"like?": types.BasicFunction(likeFunction, `like eq?, but does lossless type conversions so 1 == "1"`),

	"lt?": makeNumberComparatorFunc(
		func(a, b int64) (ast.Bool, error) { return ast.Bool(a < b), nil },
		func(a, b float64) (ast.Bool, error) { return ast.Bool(a < b), nil },
		"returns a < b",
	),
	"lte?": makeNumberComparatorFunc(
		func(a, b int64) (ast.Bool, error) { return ast.Bool(a <= b), nil },
		func(a, b float64) (ast.Bool, error) { return ast.Bool(a <= b), nil },
		"return a <= b",
	),
	"gt?": makeNumberComparatorFunc(
		func(a, b int64) (ast.Bool, error) { return ast.Bool(a > b), nil },
		func(a, b float64) (ast.Bool, error) { return ast.Bool(a > b), nil },
		"returns a > b",
	),
	"gte?": makeNumberComparatorFunc(
		func(a, b int64) (ast.Bool, error) { return ast.Bool(a >= b), nil },
		func(a, b float64) (ast.Bool, error) { return ast.Bool(a >= b), nil },
		"returns a >= b",
	),

	// types
	"type-of":   types.BasicFunction(typeOfFunction, `returns the type of a given value (e.g. "string" or "number")`),
	"to-string": types.BasicFunction(toStringFunction, "try to convert the given argument losslessly to a string"),
	"to-int":    types.BasicFunction(toIntFunction, "try to convert the given argument losslessly to an int64"),
	"to-float":  types.BasicFunction(toFloatFunction, "try to convert the given argument losslessly to a float64"),
	"to-bool":   types.BasicFunction(toBoolFunction, "try to convert the given argument losslessly to a bool"),

	// hashes
	"sha1":   types.BasicFunction(sha1Function, "return the lowercase hex representation of the SHA-1 hash"),
	"sha256": types.BasicFunction(sha256Function, "return the lowercase hex representation of the SHA-256 hash"),
	"sha512": types.BasicFunction(sha512Function, "return the lowercase hex representation of the SHA-512 hash"),

	// encoding
	"to-base64":   types.BasicFunction(toBase64Function, "apply base64 encoding to the given string"),
	"from-base64": types.BasicFunction(fromBase64Function, "decode a base64 encoded string"),

	// dates & time
	"now": types.BasicFunction(nowFunction, "returns the current date & time (UTC), formatted like a Go date"),
}
