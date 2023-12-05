// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/eval/functions"
	"go.xrstf.de/rudi/pkg/eval/types"
)

var (
	strictCoalescer   = coalescing.NewStrict()
	pedanticCoalescer = coalescing.NewPedantic()
	humaneCoalescer   = coalescing.NewHumane()

	CoreFunctions = types.Functions{
		"default": functions.NewBuilder(defaultFunction).WithDescription("returns the default value if the first argument is empty").Build(),
		"delete":  functions.NewBuilder(deleteFunction).WithBangHandler(deleteBangHandler).WithDescription("removes a key from an object or an item from a vector").Build(),
		"do":      functions.NewBuilder(doFunction).WithDescription("eval a sequence of statements where only one expression is valid").Build(),
		"empty?":  functions.NewBuilder(isEmptyFunction).WithCoalescer(humaneCoalescer).WithDescription("returns true when the given value is empty-ish (0, false, null, \"\", ...)").Build(),
		"error":   functions.NewBuilder(errorFunction, fmtErrorFunction).WithDescription("returns an error").Build(),
		"has?":    functions.NewBuilder(hasFunction).WithDescription("returns true if the given symbol's path expression points to an existing value").Build(),
		"if":      functions.NewBuilder(ifElseFunction, ifFunction).WithDescription("evaluate one of two expressions based on a condition").Build(),
		"set":     functions.NewBuilder(setFunction).WithDescription("set a value in a variable/document, only really useful with ! modifier (set!)").Build(),
		"try":     functions.NewBuilder(tryWithFallbackFunction, tryFunction).WithDescription("returns the fallback if the first expression errors out").Build(),
	}

	LogicFunctions = types.Functions{
		"and": functions.NewBuilder(andFunction).WithDescription("returns true if all arguments are true").Build(),
		"or":  functions.NewBuilder(orFunction).WithDescription("returns true if any of the arguments is true").Build(),
		"not": functions.NewBuilder(notFunction).WithDescription("negates the given argument").Build(),
	}

	ComparisonFunctions = types.Functions{
		"eq?":        functions.NewBuilder(eqFunction).WithDescription("equality check: return true if both arguments are the same").Build(),
		"identical?": functions.NewBuilder(identicalFunction).WithDescription("like `eq?`, but always uses strict coalecsing").Build(),
		"like?":      functions.NewBuilder(likeFunction).WithDescription("like `eq?`, but always uses humane coalecsing").Build(),

		"lt?":  functions.NewBuilder(ltFunction).WithDescription("returns a < b").Build(),
		"lte?": functions.NewBuilder(lteFunction).WithDescription("returns a <= b").Build(),
		"gt?":  functions.NewBuilder(gtFunction).WithDescription("returns a > b").Build(),
		"gte?": functions.NewBuilder(gteFunction).WithDescription("returns a >= b").Build(),
	}

	addRudiFunction      = functions.NewBuilder(integerAddFunction, numberAddFunction).WithDescription("returns the sum of all of its arguments").Build()
	subRudiFunction      = functions.NewBuilder(integerSubFunction, numberSubFunction).WithDescription("returns arg1 - arg2 - .. - argN").Build()
	multiplyRudiFunction = functions.NewBuilder(integerMultFunction, numberMultFunction).WithDescription("returns the product of all of its arguments").Build()
	divideRudiFunction   = functions.NewBuilder(numberDivFunction).WithDescription("returns arg1 / arg2 / .. / argN (always a floating point division, regardless of arguments)").Build()

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

	lenRudiFunction      = functions.NewBuilder(stringLenFunction, vectorLenFunction, objectLenFunction).WithDescription("returns the length of a string, vector or object").Build()
	appendRudiFunction   = functions.NewBuilder(appendToVectorFunction, appendToStringFunction).WithDescription("appends more strings to a string or arbitrary items into a vector").Build()
	prependRudiFunction  = functions.NewBuilder(prependToVectorFunction, prependToStringFunction).WithDescription("prepends more strings to a string or arbitrary items into a vector").Build()
	reverseRudiFunction  = functions.NewBuilder(reverseStringFunction, reverseVectorFunction).WithDescription("reverses a string or the elements of a vector").Build()
	containsRudiFunction = functions.NewBuilder(stringContainsFunction, vectorContainsFunction).WithDescription("returns true if a string contains a substring or a vector contains the given element").Build()

	StringsFunctions = types.Functions{
		// these ones are shared with ListsFunctions
		"len":       lenRudiFunction,
		"append":    appendRudiFunction,
		"prepend":   prependRudiFunction,
		"reverse":   reverseRudiFunction,
		"contains?": containsRudiFunction,

		"concat":      functions.NewBuilder(concatFunction).WithDescription("concatenates items in a vector using a common glue string").Build(),
		"split":       functions.NewBuilder(splitFunction, splitnFunction).WithDescription("splits a string into a vector").Build(),
		"has-prefix?": functions.NewBuilder(hasPrefixFunction).WithDescription("returns true if the given string has the prefix").Build(),
		"has-suffix?": functions.NewBuilder(hasSuffixFunction).WithDescription("returns true if the given string has the suffix").Build(),
		"trim-prefix": functions.NewBuilder(trimPrefixFunction).WithDescription("removes the prefix from the string, if it exists").Build(),
		"trim-suffix": functions.NewBuilder(trimSuffixFunction).WithDescription("removes the suffix from the string, if it exists").Build(),
		"to-lower":    functions.NewBuilder(toLowerFunction).WithDescription("returns the lowercased version of the given string").Build(),
		"to-upper":    functions.NewBuilder(toUpperFunction).WithDescription("returns the uppercased version of the given string").Build(),
		"trim":        functions.NewBuilder(trimFunction).WithDescription("returns the given whitespace with leading/trailing whitespace removed").Build(),
		"replace":     functions.NewBuilder(replaceAllFunction, replaceLimitFunction).WithDescription("returns a copy of a string with the a substring replaced by another").Build(),
	}

	ListsFunctions = types.Functions{
		// these ones are shared with StringsFunctions
		"len":       lenRudiFunction,
		"append":    appendRudiFunction,
		"prepend":   prependRudiFunction,
		"reverse":   reverseRudiFunction,
		"contains?": containsRudiFunction,

		"range": functions.
			NewBuilder(
				rangeVectorFunction,
				rangeObjectFunction,
			).
			WithDescription("allows to iterate (loop) over a vector or object").
			Build(),

		"map": functions.
			NewBuilder(
				mapVectorExpressionFunction,
				mapObjectExpressionFunction,
				mapVectorAnonymousFunction,
				mapObjectAnonymousFunction,
			).
			WithDescription("applies an expression to every element in a vector or object").
			Build(),

		"filter": functions.
			NewBuilder(
				filterVectorExpressionFunction,
				filterObjectExpressionFunction,
				filterVectorAnonymousFunction,
				filterObjectAnonymousFunction,
			).
			WithDescription("returns a copy of a given vector/object with only those elements remaining that satisfy a condition").
			Build(),
	}

	HashingFunctions = types.Functions{
		"sha1":   functions.NewBuilder(sha1Function).WithDescription("return the lowercase hex representation of the SHA-1 hash").Build(),
		"sha256": functions.NewBuilder(sha256Function).WithDescription("return the lowercase hex representation of the SHA-256 hash").Build(),
		"sha512": functions.NewBuilder(sha512Function).WithDescription("return the lowercase hex representation of the SHA-512 hash").Build(),
	}

	EncodingFunctions = types.Functions{
		"to-base64":   functions.NewBuilder(toBase64Function).WithDescription("apply base64 encoding to the given string").Build(),
		"from-base64": functions.NewBuilder(fromBase64Function).WithDescription("decode a base64 encoded string").Build(),
		"to-json":     functions.NewBuilder(toJSONFunction).WithDescription("encode the given value using JSON").Build(),
		"from-json":   functions.NewBuilder(fromJSONFunction).WithDescription("decode a JSON string").Build(),
	}

	DateTimeFunctions = types.Functions{
		"now": functions.NewBuilder(nowFunction).WithDescription("returns the current date & time (UTC), formatted like a Go date").Build(),
	}

	TypeFunctions = types.Functions{
		"type-of": functions.NewBuilder(typeOfFunction).WithDescription(`returns the type of a given value (e.g. "string" or "number")`).Build(),

		// these functions purposefully always uses humane coalescing
		"to-bool":   functions.NewBuilder(toBoolFunction).WithCoalescer(humaneCoalescer).WithDescription("try to convert the given argument losslessly to a bool").Build(),
		"to-float":  functions.NewBuilder(toFloatFunction).WithCoalescer(humaneCoalescer).WithDescription("try to convert the given argument losslessly to a float64").Build(),
		"to-int":    functions.NewBuilder(toIntFunction).WithCoalescer(humaneCoalescer).WithDescription("try to convert the given argument losslessly to an int64").Build(),
		"to-string": functions.NewBuilder(toStringFunction).WithCoalescer(humaneCoalescer).WithDescription("try to convert the given argument losslessly to a string").Build(),
	}

	CoalescingContextFunctions = types.Functions{
		"strictly":     functions.NewBuilder(doFunction).WithCoalescer(strictCoalescer).WithDescription("evaluates the child expressions using strict coalescing").Build(),
		"pedantically": functions.NewBuilder(doFunction).WithCoalescer(pedanticCoalescer).WithDescription("evaluates the child expressions using pedantic coalescing").Build(),
		"humanely":     functions.NewBuilder(doFunction).WithCoalescer(humaneCoalescer).WithDescription("evaluates the child expressions using humane coalescing").Build(),
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
