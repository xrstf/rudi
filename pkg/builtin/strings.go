// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"
	"strings"

	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/eval/util"
)

// (concat GLUE:String ELEMENTS:(Vector/String)+)
func concatFunction(ctx types.Context, args []any) (any, error) {
	glue, err := ctx.Coalesce().ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	parts := []string{}
	for i, value := range args[1:] {
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

	return strings.Join(parts, string(glue)), nil
}

type stringifiedLiteralFunction func(args []string) (any, error)

func stringifyArgs(fun stringifiedLiteralFunction) util.LiteralFunction {
	return func(ctx types.Context, args []any) (any, error) {
		stringValues := make([]string, len(args))
		for i, value := range args {
			strValue, err := ctx.Coalesce().ToString(value)
			if err != nil {
				return nil, fmt.Errorf("argument #%d is not a string: %w", i, err)
			}

			stringValues[i] = string(strValue)
		}

		return fun(stringValues)
	}
}

// (split SEP:String SOURCE:String)
func splitFunction(args []string) (any, error) {
	parts := strings.Split(args[1], args[0])

	// to []any
	result := make([]any, len(parts))
	for i, part := range parts {
		result[i] = part
	}

	return result, nil
}

// (has-suffix SOURCE:String SUFFIX:String)
func hasSuffixFunction(args []string) (any, error) {
	result := strings.HasSuffix(args[0], args[1])

	return result, nil
}

// (has-prefix SOURCE:String PREFIX:String)
func hasPrefixFunction(args []string) (any, error) {
	result := strings.HasPrefix(args[0], args[1])

	return result, nil
}

// (trim-suffix SOURCE:String SUFFIX:String)
func trimSuffixFunction(args []string) (any, error) {
	result := strings.TrimSuffix(args[0], args[1])

	return result, nil
}

// (trim-prefix SOURCE:String PREFIX:String)
func trimPrefixFunction(args []string) (any, error) {
	result := strings.TrimPrefix(args[0], args[1])

	return result, nil
}

// (to-lower SOURCE:String)
func toLowerFunction(args []string) (any, error) {
	result := strings.ToLower(args[0])

	return result, nil
}

// (to-upper SOURCE:String)
func toUpperFunction(args []string) (any, error) {
	result := strings.ToUpper(args[0])

	return result, nil
}

// (trim SOURCE:String)
func trimFunction(args []string) (any, error) {
	result := strings.TrimSpace(args[0])

	return result, nil
}
