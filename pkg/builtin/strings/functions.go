// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package strings

import (
	"fmt"
	stdstrings "strings"

	"go.xrstf.de/rudi/pkg/equality"
	"go.xrstf.de/rudi/pkg/runtime/functions"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

var (
	Functions = types.Functions{
		// these ones also work with lists
		"len":       functions.NewBuilder(stringLenFunction, vectorLenFunction, objectLenFunction).WithDescription("returns the length of a string, vector or object").Build(),
		"append":    functions.NewBuilder(appendToVectorFunction, appendToStringFunction).WithDescription("appends more strings to a string or arbitrary items into a vector").Build(),
		"prepend":   functions.NewBuilder(prependToVectorFunction, prependToStringFunction).WithDescription("prepends more strings to a string or arbitrary items into a vector").Build(),
		"reverse":   functions.NewBuilder(reverseStringFunction, reverseVectorFunction).WithDescription("reverses a string or the elements of a vector").Build(),
		"contains?": functions.NewBuilder(stringContainsFunction, vectorContainsFunction).WithDescription("returns true if a string contains a substring or a vector contains the given element").Build(),

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
)

func vectorLenFunction(vec []any) (any, error) {
	return len(vec), nil
}

func objectLenFunction(obj map[string]any) (any, error) {
	return len(obj), nil
}

func stringLenFunction(s string) (any, error) {
	return len(s), nil
}

// (append VEC ITEMS+)
func appendToVectorFunction(base []any, args ...any) (any, error) {
	result := []any{}
	result = append(result, base...)
	result = append(result, args...)

	return result, nil
}

// (append STR ITEMS+)
func appendToStringFunction(base string, args ...string) (any, error) {
	return base + stdstrings.Join(args, ""), nil
}

// (prepend VEC ITEMS+)
func prependToVectorFunction(base []any, args ...any) (any, error) {
	return append(args, base...), nil
}

// (prepend STR ITEMS+)
func prependToStringFunction(base string, args ...string) (any, error) {
	return stdstrings.Join(args, "") + base, nil
}

// (reverse STR)
func reverseStringFunction(s string) (any, error) {
	// thank you https://stackoverflow.com/a/10030772
	result := []rune(s)
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result), nil
}

// (reverse VEC)
func reverseVectorFunction(vec []any) (any, error) {
	// clone original data
	result := append([]any{}, vec...)

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}

// (contains? STR STR)
func stringContainsFunction(haystack string, needle string) (any, error) {
	return stdstrings.Contains(haystack, needle), nil
}

// (contains? VEC ITEM)
func vectorContainsFunction(ctx types.Context, haystack []any, needle any) (any, error) {
	for _, val := range haystack {
		equal, err := equality.Equal(ctx.Coalesce(), val, needle)
		if err != nil {
			return false, err
		}
		if equal {
			return true, nil
		}
	}

	return false, nil
}

// (concat GLUE:String ELEMENTS:(Vector/String)+)
func concatFunction(ctx types.Context, glue string, args ...any) (any, error) {
	parts := []string{}
	for i, value := range args {
		vector, err := ctx.Coalesce().ToVector(value)
		if err != nil {
			part, err := ctx.Coalesce().ToString(value)
			if err != nil {
				return nil, fmt.Errorf("argument #%d is neither vector nor string, but %T", i+1, value)
			}

			parts = append(parts, string(part))
			continue
		}

		for j, item := range vector {
			part, err := ctx.Coalesce().ToString(item)
			if err != nil {
				return nil, fmt.Errorf("argument #%d.%d: %w", i+1, j, err)
			}

			parts = append(parts, string(part))
		}
	}

	return stdstrings.Join(parts, string(glue)), nil
}

func splitFunction(sep string, source string) (any, error) {
	parts := stdstrings.Split(source, sep)

	// to []any
	result := make([]any, len(parts))
	for i, part := range parts {
		result[i] = part
	}

	return result, nil
}

func splitnFunction(sep string, source string, limit int64) (any, error) {
	parts := stdstrings.SplitN(source, sep, int(limit))

	// to []any
	result := make([]any, len(parts))
	for i, part := range parts {
		result[i] = part
	}

	return result, nil
}

func hasSuffixFunction(source string, suffix string) (any, error) {
	return stdstrings.HasSuffix(source, suffix), nil
}

func hasPrefixFunction(source string, prefix string) (any, error) {
	return stdstrings.HasPrefix(source, prefix), nil
}

func trimSuffixFunction(source string, suffix string) (any, error) {
	return stdstrings.TrimSuffix(source, suffix), nil
}

func trimPrefixFunction(source string, prefix string) (any, error) {
	return stdstrings.TrimPrefix(source, prefix), nil
}

func toLowerFunction(s string) (any, error) {
	return stdstrings.ToLower(s), nil
}

func toUpperFunction(s string) (any, error) {
	return stdstrings.ToUpper(s), nil
}

func trimFunction(s string) (any, error) {
	return stdstrings.TrimSpace(s), nil
}

func replaceAllFunction(s, old, new string) (any, error) {
	return stdstrings.ReplaceAll(s, old, new), nil
}

func replaceLimitFunction(s, old, new string, limit int64) (any, error) {
	return stdstrings.Replace(s, old, new, int(limit)), nil
}
