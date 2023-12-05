// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/eval/util"
	"go.xrstf.de/rudi/pkg/eval/util/native"
)

var (
	strictCoalescer   = coalescing.NewStrict()
	pedanticCoalescer = coalescing.NewPedantic()
	humaneCoalescer   = coalescing.NewHumane()

	CoreFunctions = types.Functions{
		"default": native.NewFunction("returns the default value if the first argument is empty", defaultFunction),
		"delete":  deleteFunction{},
		"do":      native.NewFunction("eval a sequence of statements where only one expression is valid", doFunction),
		"empty?":  native.NewFunction("returns true when the given value is empty-ish (0, false, null, \"\", ...)", isEmptyFunction).WithCoalescer(humaneCoalescer),
		"error":   native.NewFunction("returns an error", errorFunction),
		"has?":    native.NewFunction("returns true if the given symbol's path expression points to an existing value", hasFunction),
		"if":      native.NewFunction("evaluate one of two expressions based on a condition", ifElseFunction, ifFunction),
		"set":     native.NewFunction("set a value in a variable/document, only really useful with ! modifier (set!)", setFunction),
		"try":     native.NewFunction("returns the fallback if the first expression errors out", tryWithFallbackFunction, tryFunction),
	}

	LogicFunctions = types.Functions{
		"and": native.NewFunction("returns true if all arguments are true", andFunction),
		"or":  native.NewFunction("returns true if any of the arguments is true", orFunction),
		"not": native.NewFunction("negates the given argument", notFunction),
	}

	ComparisonFunctions = types.Functions{
		"eq?":        native.NewFunction("equality check: return true if both arguments are the same", eqFunction),
		"identical?": native.NewFunction("like `eq?`, but always uses strict coalecsing", identicalFunction),
		"like?":      native.NewFunction("like `eq?`, but always uses humane coalecsing", likeFunction),

		"lt?":  native.NewFunction("returns a < b", ltCoalescer),
		"lte?": native.NewFunction("returns a <= b", lteCoalescer),
		"gt?":  native.NewFunction("returns a > b", gtCoalescer),
		"gte?": native.NewFunction("returns a >= b", gteCoalescer),
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

	lenRudiFunction      = native.NewFunction("returns the length of a string, vector or object", stringLenFunction, vectorLenFunction, objectLenFunction)
	appendRudiFunction   = native.NewFunction("appends more strings to a string or arbitrary items into a vector", appendToVectorFunction, appendToStringFunction)
	prependRudiFunction  = native.NewFunction("prepends more strings to a string or arbitrary items into a vector", prependToVectorFunction, prependToStringFunction)
	reverseRudiFunction  = native.NewFunction("reverses a string or the elements of a vector", reverseVectorFunction, reverseStringFunction)
	containsRudiFunction = native.NewFunction("returns true if a string contains a substring or a vector contains the given element", stringContainsFunction, vectorContainsFunction)

	StringsFunctions = types.Functions{
		// these ones are shared with ListsFunctions
		"len":       lenRudiFunction,
		"append":    appendRudiFunction,
		"prepend":   prependRudiFunction,
		"reverse":   reverseRudiFunction,
		"contains?": containsRudiFunction,

		"concat":      native.NewFunction("concatenates items in a vector using a common glue string", concatFunction),
		"split":       native.NewFunction("splits a string into a vector", splitFunction),
		"has-prefix?": native.NewFunction("returns true if the given string has the prefix", hasPrefixFunction),
		"has-suffix?": native.NewFunction("returns true if the given string has the suffix", hasSuffixFunction),
		"trim-prefix": native.NewFunction("removes the prefix from the string, if it exists", trimPrefixFunction),
		"trim-suffix": native.NewFunction("removes the suffix from the string, if it exists", trimSuffixFunction),
		"to-lower":    native.NewFunction("returns the lowercased version of the given string", toLowerFunction),
		"to-upper":    native.NewFunction("returns the uppercased version of the given string", toUpperFunction),
		"trim":        native.NewFunction("returns the given whitespace with leading/trailing whitespace removed", trimFunction),
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
		"sha1":   native.NewFunction("return the lowercase hex representation of the SHA-1 hash", sha1Function),
		"sha256": native.NewFunction("return the lowercase hex representation of the SHA-256 hash", sha256Function),
		"sha512": native.NewFunction("return the lowercase hex representation of the SHA-512 hash", sha512Function),
	}

	EncodingFunctions = types.Functions{
		"to-base64":   native.NewFunction("apply base64 encoding to the given string", toBase64Function),
		"from-base64": native.NewFunction("decode a base64 encoded string", fromBase64Function),
	}

	DateTimeFunctions = types.Functions{
		"now": native.NewFunction("returns the current date & time (UTC), formatted like a Go date", nowFunction),
	}

	TypeFunctions = types.Functions{
		"type-of": native.NewFunction(`returns the type of a given value (e.g. "string" or "number")`, typeOfFunction),

		// these functions purposefully always uses humane coalescing
		"to-bool":   native.NewFunction("try to convert the given argument losslessly to a bool", toBoolFunction).WithCoalescer(humaneCoalescer),
		"to-float":  native.NewFunction("try to convert the given argument losslessly to a float64", toFloatFunction).WithCoalescer(humaneCoalescer),
		"to-int":    native.NewFunction("try to convert the given argument losslessly to an int64", toIntFunction).WithCoalescer(humaneCoalescer),
		"to-string": native.NewFunction("try to convert the given argument losslessly to a string", toStringFunction).WithCoalescer(humaneCoalescer),
	}

	CoalescingContextFunctions = types.Functions{
		"strictly":     native.NewFunction("evaluates the child expressions using strict coalescing", doFunction).WithCoalescer(strictCoalescer),
		"pedantically": native.NewFunction("evaluates the child expressions using pedantic coalescing", doFunction).WithCoalescer(pedanticCoalescer),
		"humanely":     native.NewFunction("evaluates the child expressions using humane coalescing", doFunction).WithCoalescer(humaneCoalescer),
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
