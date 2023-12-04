// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/eval/util"
)

var (
	CoreFunctions = types.Functions{
		"default": util.NewRawFunction(defaultFunction, "returns the default value if the first argument is empty").MinArgs(2).MaxArgs(2),
		"delete":  deleteFunction{},
		"do":      util.NewRawFunction(doFunction, "eval a sequence of statements where only one expression is valid").MinArgs(1),
		"empty?":  util.NewLiteralFunction(isEmptyFunction, "returns true when the given value is empty-ish (0, false, null, \"\", ...)").MinArgs(1).MaxArgs(1),
		"error":   util.NewLiteralFunction(errorFunction, "returns an error").MinArgs(1),
		"has?":    util.NewRawFunction(hasFunction, "returns true if the given symbol's path expression points to an existing value").MinArgs(1).MaxArgs(1),
		"if":      util.NewRawFunction(ifFunction, "evaluate one of two expressions based on a condition").MinArgs(2).MaxArgs(3),
		"set":     util.NewRawFunction(setFunction, "set a value in a variable/document, only really useful with ! modifier (set!)").MinArgs(2).MaxArgs(2),
		"try":     util.NewRawFunction(tryFunction, "returns the fallback if the first expression errors out").MinArgs(1).MaxArgs(2),
	}

	LogicFunctions = types.Functions{
		"and": util.NewRawFunction(andFunction, "returns true if all arguments are true").MinArgs(1),
		"or":  util.NewRawFunction(orFunction, "returns true if any of the arguments is true").MinArgs(1),
		"not": util.NewLiteralFunction(notFunction, "negates the given argument").MinArgs(1).MaxArgs(1),
	}

	ComparisonFunctions = types.Functions{
		"eq?":        util.NewLiteralFunction(eqFunction, "equality check: return true if both arguments are the same").MinArgs(2).MaxArgs(2),
		"identical?": util.NewLiteralFunction(identicalFunction, "like `eq?`, but always uses strict coalecsing").MinArgs(2).MaxArgs(2),
		"like?":      util.NewLiteralFunction(likeFunction, "like `eq?`, but always uses humane coalecsing").MinArgs(2).MaxArgs(2),

		"lt?":  util.NewLiteralFunction(ltCoalescer, "returns a < b").MinArgs(2).MaxArgs(2),
		"lte?": util.NewLiteralFunction(lteCoalescer, "returns a <= b").MinArgs(2).MaxArgs(2),
		"gt?":  util.NewLiteralFunction(gtCoalescer, "returns a > b").MinArgs(2).MaxArgs(2),
		"gte?": util.NewLiteralFunction(gteCoalescer, "returns a >= b").MinArgs(2).MaxArgs(2),
	}

	// aliases to make bang functions nicer (add! vs +!)
	addRudiFunction      = util.NewLiteralFunction(numberifyArgs(addFunction), "returns the sum of all of its arguments").MinArgs(2)
	subRudiFunction      = util.NewLiteralFunction(numberifyArgs(subFunction), "returns arg1 - arg2 - .. - argN").MinArgs(2)
	multiplyRudiFunction = util.NewLiteralFunction(numberifyArgs(multiplyFunction), "returns the product of all of its arguments").MinArgs(2)
	divideRudiFunction   = util.NewLiteralFunction(numberifyArgs(divideFunction), "returns arg1 / arg2 / .. / argN").MinArgs(2)

	MathFunctions = types.Functions{
		"+": addRudiFunction,
		"-": subRudiFunction,
		"*": multiplyRudiFunction,
		"/": divideRudiFunction,

		// aliases to make bang functions nicer (add! vs +!)
		"add":  addRudiFunction,
		"sub":  subRudiFunction,
		"mult": multiplyRudiFunction,
		"div":  divideRudiFunction,
	}

	lenRudiFunction      = util.NewLiteralFunction(lenFunction, "returns the length of a string, vector or object").MinArgs(1).MaxArgs(1)
	appendRudiFunction   = util.NewLiteralFunction(appendFunction, "appends more strings to a string or arbitrary items into a vector").MinArgs(2)
	prependRudiFunction  = util.NewLiteralFunction(prependFunction, "prepends more strings to a string or arbitrary items into a vector").MinArgs(2)
	reverseRudiFunction  = util.NewLiteralFunction(reverseFunction, "reverses a string or the elements of a vector").MinArgs(1).MaxArgs(1)
	containsRudiFunction = util.NewLiteralFunction(containsFunction, "returns true if a string contains a substring or a vector contains the given element").MinArgs(2).MaxArgs(2)

	StringsFunctions = types.Functions{
		// these ones are shared with ListsFunctions
		"len":       lenRudiFunction,
		"append":    appendRudiFunction,
		"prepend":   prependRudiFunction,
		"reverse":   reverseRudiFunction,
		"contains?": containsRudiFunction,

		"concat":      util.NewLiteralFunction(concatFunction, "concatenates items in a vector using a common glue string").MinArgs(2),
		"split":       util.NewLiteralFunction(stringifyArgs(splitFunction), "splits a string into a vector").MinArgs(2).MaxArgs(2),
		"has-prefix?": util.NewLiteralFunction(stringifyArgs(hasPrefixFunction), "returns true if the given string has the prefix").MinArgs(2).MaxArgs(2),
		"has-suffix?": util.NewLiteralFunction(stringifyArgs(hasSuffixFunction), "returns true if the given string has the suffix").MinArgs(2).MaxArgs(2),
		"trim-prefix": util.NewLiteralFunction(stringifyArgs(trimPrefixFunction), "removes the prefix from the string, if it exists").MinArgs(2).MaxArgs(2),
		"trim-suffix": util.NewLiteralFunction(stringifyArgs(trimSuffixFunction), "removes the suffix from the string, if it exists").MinArgs(2).MaxArgs(2),
		"to-lower":    util.NewLiteralFunction(stringifyArgs(toLowerFunction), "returns the lowercased version of the given string").MinArgs(1).MaxArgs(1),
		"to-upper":    util.NewLiteralFunction(stringifyArgs(toUpperFunction), "returns the uppercased version of the given string").MinArgs(1).MaxArgs(1),
		"trim":        util.NewLiteralFunction(stringifyArgs(trimFunction), "returns the given whitespace with leading/trailing whitespace removed").MinArgs(1).MaxArgs(1),
	}

	ListsFunctions = types.Functions{
		// these ones are shared with StringsFunctions
		"len":       lenRudiFunction,
		"append":    appendRudiFunction,
		"prepend":   prependRudiFunction,
		"reverse":   reverseRudiFunction,
		"contains?": containsRudiFunction,

		"range":  util.NewRawFunction(rangeFunction, "allows to iterate (loop) over a vector or object").MinArgs(3),
		"map":    util.NewRawFunction(mapFunction, "applies an expression to every element in a vector or object").MinArgs(2),
		"filter": util.NewRawFunction(filterFunction, "returns a copy of a given vector/object with only those elements remaining that satisfy a condition").MinArgs(2),
	}

	HashingFunctions = types.Functions{
		"sha1":   util.NewLiteralFunction(sha1Function, "return the lowercase hex representation of the SHA-1 hash").MinArgs(1).MaxArgs(1),
		"sha256": util.NewLiteralFunction(sha256Function, "return the lowercase hex representation of the SHA-256 hash").MinArgs(1).MaxArgs(1),
		"sha512": util.NewLiteralFunction(sha512Function, "return the lowercase hex representation of the SHA-512 hash").MinArgs(1).MaxArgs(1),
	}

	EncodingFunctions = types.Functions{
		"to-base64":   util.NewLiteralFunction(toBase64Function, "apply base64 encoding to the given string").MinArgs(1).MaxArgs(1),
		"from-base64": util.NewLiteralFunction(fromBase64Function, "decode a base64 encoded string").MinArgs(1).MaxArgs(1),
	}

	DateTimeFunctions = types.Functions{
		"now": util.NewLiteralFunction(nowFunction, "returns the current date & time (UTC), formatted like a Go date").MinArgs(1).MaxArgs(1),
	}

	TypeFunctions = types.Functions{
		"type-of":   util.NewLiteralFunction(typeOfFunction, `returns the type of a given value (e.g. "string" or "number")`).MinArgs(1).MaxArgs(1),
		"to-bool":   util.NewLiteralFunction(toBoolFunction, "try to convert the given argument losslessly to a bool").MinArgs(1).MaxArgs(1),
		"to-float":  util.NewLiteralFunction(toFloatFunction, "try to convert the given argument losslessly to a float64").MinArgs(1).MaxArgs(1),
		"to-int":    util.NewLiteralFunction(toIntFunction, "try to convert the given argument losslessly to an int64").MinArgs(1).MaxArgs(1),
		"to-string": util.NewLiteralFunction(toStringFunction, "try to convert the given argument losslessly to a string").MinArgs(1).MaxArgs(1),
	}

	CoalescingContextFunctions = types.Functions{
		"strictly":     util.NewRawFunction(strictlyFunction, "evaluates the child expressions using strict coalescing").MinArgs(1).MaxArgs(1),
		"pedantically": util.NewRawFunction(pedanticallyFunction, "evaluates the child expressions using pedantic coalescing").MinArgs(1).MaxArgs(1),
		"humanely":     util.NewRawFunction(humanelyFunction, "evaluates the child expressions using humane coalescing").MinArgs(1).MaxArgs(1),
	}

	AllFunctions = types.Functions{}.
			Add(CoreFunctions).
			Add(LogicFunctions).
			Add(ComparisonFunctions).
			Add(MathFunctions).
			Add(StringsFunctions).
			Add(ListsFunctions).
			Add(HashingFunctions).
			Add(EncodingFunctions).
			Add(DateTimeFunctions).
			Add(TypeFunctions).
			Add(CoalescingContextFunctions)
)
