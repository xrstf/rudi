// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"
	"strings"

	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

// (concat GLUE:String ELEMENTS:(Vector/String)+)
func concatFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size < 2 {
		return nil, fmt.Errorf("expected 2+ arguments, got %d", size)
	}

	_, glue, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	glueString, err := ctx.Coalesce().ToString(glue)
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	values, err := evalArgs(ctx, args, 1)
	if err != nil {
		return nil, err
	}

	parts := []string{}
	for i, value := range values {
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

	return strings.Join(parts, string(glueString)), nil
}

type genericStringFunc func(ctx types.Context, args []string) (any, error)

func fromStringFunc(f genericStringFunc, expectedArgs int, desc string) types.Function {
	return types.BasicFunction(func(ctx types.Context, args []ast.Expression) (any, error) {
		if size := len(args); size != expectedArgs {
			return nil, fmt.Errorf("expected %d argument(s), got %d", expectedArgs, size)
		}

		values, err := evalArgs(ctx, args, 0)
		if err != nil {
			return nil, err
		}

		stringValues := make([]string, len(values))
		for i, value := range values {
			strValue, err := ctx.Coalesce().ToString(value)
			if err != nil {
				return nil, fmt.Errorf("argument #%d: %w", i, err)
			}

			stringValues[i] = string(strValue)
		}

		return f(ctx, stringValues)
	}, desc)
}

// (split SEP:String SOURCE:String)
func splitFunction(ctx types.Context, args []string) (any, error) {
	parts := strings.Split(args[1], args[0])

	// to []any
	result := make([]any, len(parts))
	for i, part := range parts {
		result[i] = part
	}

	return result, nil
}

// (has-suffix SOURCE:String SUFFIX:String)
func hasSuffixFunction(ctx types.Context, args []string) (any, error) {
	result := strings.HasSuffix(args[0], args[1])

	return result, nil
}

// (has-prefix SOURCE:String PREFIX:String)
func hasPrefixFunction(ctx types.Context, args []string) (any, error) {
	result := strings.HasPrefix(args[0], args[1])

	return result, nil
}

// (trim-suffix SOURCE:String SUFFIX:String)
func trimSuffixFunction(ctx types.Context, args []string) (any, error) {
	result := strings.TrimSuffix(args[0], args[1])

	return result, nil
}

// (trim-prefix SOURCE:String PREFIX:String)
func trimPrefixFunction(ctx types.Context, args []string) (any, error) {
	result := strings.TrimPrefix(args[0], args[1])

	return result, nil
}

// (to-lower SOURCE:String)
func toLowerFunction(ctx types.Context, args []string) (any, error) {
	result := strings.ToLower(args[0])

	return result, nil
}

// (to-upper SOURCE:String)
func toUpperFunction(ctx types.Context, args []string) (any, error) {
	result := strings.ToUpper(args[0])

	return result, nil
}

// (trim SOURCE:String)
func trimFunction(ctx types.Context, args []string) (any, error) {
	result := strings.TrimSpace(args[0])

	return result, nil
}
