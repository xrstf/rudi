// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/eval/util/native"
)

var (
	strictCoalescer   = coalescing.NewStrict()
	pedanticCoalescer = coalescing.NewPedantic()
	humaneCoalescer   = coalescing.NewHumane()

	CoreFunctions = types.Functions{
		"default": native.NewFunction(defaultFunction).WithDescription("returns the default value if the first argument is empty"),
		"delete":  deleteFunction{},
		"do":      native.NewFunction(doFunction).WithDescription("eval a sequence of statements where only one expression is valid"),
		"empty?":  native.NewFunction(isEmptyFunction).WithCoalescer(humaneCoalescer).WithDescription("returns true when the given value is empty-ish (0, false, null, \"\", ...)"),
		"error":   native.NewFunction(errorFunction).WithDescription("returns an error"),
		"has?":    native.NewFunction(hasFunction).WithDescription("returns true if the given symbol's path expression points to an existing value"),
		"if":      native.NewFunction(ifElseFunction, ifFunction).WithDescription("evaluate one of two expressions based on a condition"),
		"set":     native.NewFunction(setFunction).WithDescription("set a value in a variable/document, only really useful with ! modifier (set!)"),
		"try":     native.NewFunction(tryWithFallbackFunction, tryFunction).WithDescription("returns the fallback if the first expression errors out"),
	}

	LogicFunctions = types.Functions{
		"and": native.NewFunction(andFunction).WithDescription("returns true if all arguments are true"),
		"or":  native.NewFunction(orFunction).WithDescription("returns true if any of the arguments is true"),
		"not": native.NewFunction(notFunction).WithDescription("negates the given argument"),
	}

	ComparisonFunctions = types.Functions{
		"eq?":        native.NewFunction(eqFunction).WithDescription("equality check: return true if both arguments are the same"),
		"identical?": native.NewFunction(identicalFunction).WithDescription("like `eq?`, but always uses strict coalecsing"),
		"like?":      native.NewFunction(likeFunction).WithDescription("like `eq?`, but always uses humane coalecsing"),

		"lt?":  native.NewFunction(ltCoalescer).WithDescription("returns a < b"),
		"lte?": native.NewFunction(lteCoalescer).WithDescription("returns a <= b"),
		"gt?":  native.NewFunction(gtCoalescer).WithDescription("returns a > b"),
		"gte?": native.NewFunction(gteCoalescer).WithDescription("returns a >= b"),
	}

	// aliases to make bang functions nicer (add! vs +!)
	addRudiFunction      = native.NewFunction(numberAddFunction, integerAddFunction).WithDescription("returns the sum of all of its arguments")
	subRudiFunction      = native.NewFunction(numberSubFunction, integerSubFunction).WithDescription("returns arg1 - arg2 - .. - argN")
	multiplyRudiFunction = native.NewFunction(numberMultFunction, integerMultFunction).WithDescription("returns the product of all of its arguments")
	divideRudiFunction   = native.NewFunction(numberDivFunction, integerDivFunction).WithDescription("returns arg1 / arg2 / .. / argN")

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

	lenRudiFunction      = native.NewFunction(stringLenFunction, vectorLenFunction, objectLenFunction).WithDescription("returns the length of a string, vector or object")
	appendRudiFunction   = native.NewFunction(appendToVectorFunction, appendToStringFunction).WithDescription("appends more strings to a string or arbitrary items into a vector")
	prependRudiFunction  = native.NewFunction(prependToVectorFunction, prependToStringFunction).WithDescription("prepends more strings to a string or arbitrary items into a vector")
	reverseRudiFunction  = native.NewFunction(reverseVectorFunction, reverseStringFunction).WithDescription("reverses a string or the elements of a vector")
	containsRudiFunction = native.NewFunction(stringContainsFunction, vectorContainsFunction).WithDescription("returns true if a string contains a substring or a vector contains the given element")

	StringsFunctions = types.Functions{
		// these ones are shared with ListsFunctions
		"len":       lenRudiFunction,
		"append":    appendRudiFunction,
		"prepend":   prependRudiFunction,
		"reverse":   reverseRudiFunction,
		"contains?": containsRudiFunction,

		"concat":      native.NewFunction(concatFunction).WithDescription("concatenates items in a vector using a common glue string"),
		"split":       native.NewFunction(splitFunction).WithDescription("splits a string into a vector"),
		"has-prefix?": native.NewFunction(hasPrefixFunction).WithDescription("returns true if the given string has the prefix"),
		"has-suffix?": native.NewFunction(hasSuffixFunction).WithDescription("returns true if the given string has the suffix"),
		"trim-prefix": native.NewFunction(trimPrefixFunction).WithDescription("removes the prefix from the string, if it exists"),
		"trim-suffix": native.NewFunction(trimSuffixFunction).WithDescription("removes the suffix from the string, if it exists"),
		"to-lower":    native.NewFunction(toLowerFunction).WithDescription("returns the lowercased version of the given string"),
		"to-upper":    native.NewFunction(toUpperFunction).WithDescription("returns the uppercased version of the given string"),
		"trim":        native.NewFunction(trimFunction).WithDescription("returns the given whitespace with leading/trailing whitespace removed"),
	}

	ListsFunctions = types.Functions{
		// these ones are shared with StringsFunctions
		"len":       lenRudiFunction,
		"append":    appendRudiFunction,
		"prepend":   prependRudiFunction,
		"reverse":   reverseRudiFunction,
		"contains?": containsRudiFunction,

		"range": native.NewFunction(
			rangeVectorFunction,
			rangeObjectFunction,
		).WithDescription("allows to iterate (loop) over a vector or object"),

		"map": native.NewFunction(
			mapVectorExpressionFunction,
			mapObjectExpressionFunction,
			mapVectorAnonymousFunction,
			mapObjectAnonymousFunction,
		).WithDescription("applies an expression to every element in a vector or object"),

		"filter": native.NewFunction(
			filterVectorExpressionFunction,
			filterObjectExpressionFunction,
			filterVectorAnonymousFunction,
			filterObjectAnonymousFunction,
		).WithDescription("returns a copy of a given vector/object with only those elements remaining that satisfy a condition"),
	}

	HashingFunctions = types.Functions{
		"sha1":   native.NewFunction(sha1Function).WithDescription("return the lowercase hex representation of the SHA-1 hash"),
		"sha256": native.NewFunction(sha256Function).WithDescription("return the lowercase hex representation of the SHA-256 hash"),
		"sha512": native.NewFunction(sha512Function).WithDescription("return the lowercase hex representation of the SHA-512 hash"),
	}

	EncodingFunctions = types.Functions{
		"to-base64":   native.NewFunction(toBase64Function).WithDescription("apply base64 encoding to the given string"),
		"from-base64": native.NewFunction(fromBase64Function).WithDescription("decode a base64 encoded string"),
	}

	DateTimeFunctions = types.Functions{
		"now": native.NewFunction(nowFunction).WithDescription("returns the current date & time (UTC), formatted like a Go date"),
	}

	TypeFunctions = types.Functions{
		"type-of": native.NewFunction(typeOfFunction).WithDescription(`returns the type of a given value (e.g. "string" or "number")`),

		// these functions purposefully always uses humane coalescing
		"to-bool":   native.NewFunction(toBoolFunction).WithCoalescer(humaneCoalescer).WithDescription("try to convert the given argument losslessly to a bool"),
		"to-float":  native.NewFunction(toFloatFunction).WithCoalescer(humaneCoalescer).WithDescription("try to convert the given argument losslessly to a float64"),
		"to-int":    native.NewFunction(toIntFunction).WithCoalescer(humaneCoalescer).WithDescription("try to convert the given argument losslessly to an int64"),
		"to-string": native.NewFunction(toStringFunction).WithCoalescer(humaneCoalescer).WithDescription("try to convert the given argument losslessly to a string"),
	}

	CoalescingContextFunctions = types.Functions{
		"strictly":     native.NewFunction(doFunction).WithCoalescer(strictCoalescer).WithDescription("evaluates the child expressions using strict coalescing"),
		"pedantically": native.NewFunction(doFunction).WithCoalescer(pedanticCoalescer).WithDescription("evaluates the child expressions using pedantic coalescing"),
		"humanely":     native.NewFunction(doFunction).WithCoalescer(humaneCoalescer).WithDescription("evaluates the child expressions using humane coalescing"),
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
