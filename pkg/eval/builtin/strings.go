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

	glueString, ok := glue.(ast.String)
	if !ok {
		return nil, fmt.Errorf("glue is not string, but %T", glue)
	}

	values, err := evalArgs(ctx, args, 1)
	if err != nil {
		return nil, err
	}

	parts := []string{}
	for i, value := range values {
		vector, ok := value.(ast.Vector)
		if !ok {
			part, ok := value.(ast.String)
			if !ok {
				return nil, fmt.Errorf("argument #%d is neither vector nor string, but %T", i+1, value)
			}

			parts = append(parts, string(part))
			continue
		}

		for j, item := range vector.Data {
			part, ok := item.(ast.String)
			if !ok {
				return nil, fmt.Errorf("argument #%d.%d is not a string, but %T", i+1, j, item)
			}

			parts = append(parts, string(part))
		}
	}

	return ast.String(strings.Join(parts, string(glueString))), nil
}

type genericStringFunc func(ctx types.Context, args []string) (any, error)

func fromStringFunc(f genericStringFunc, expectedArgs int) types.Function {
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
			strValue, ok := value.(ast.String)
			if !ok {
				return nil, fmt.Errorf("argument #%d is not a string, but %T", i, value)
			}

			stringValues[i] = string(strValue)
		}

		return f(ctx, stringValues)
	}, "")
}

// (split SEP:String SOURCE:String)
func splitFunction(ctx types.Context, args []string) (any, error) {
	parts := strings.Split(args[1], args[0])
	result := make([]any, len(parts))
	for i, part := range parts {
		result[i] = ast.String(part)
	}

	return ast.Vector{Data: result}, nil
}

// (has-suffix SOURCE:String SUFFIX:String)
func hasSuffixFunction(ctx types.Context, args []string) (any, error) {
	result := strings.HasSuffix(args[0], args[1])

	return ast.Bool(result), nil
}

// (has-prefix SOURCE:String PREFIX:String)
func hasPrefixFunction(ctx types.Context, args []string) (any, error) {
	result := strings.HasPrefix(args[0], args[1])

	return ast.Bool(result), nil
}

// (trim-suffix SOURCE:String SUFFIX:String)
func trimSuffixFunction(ctx types.Context, args []string) (any, error) {
	result := strings.TrimSuffix(args[0], args[1])

	return ast.String(result), nil
}

// (trim-prefix SOURCE:String PREFIX:String)
func trimPrefixFunction(ctx types.Context, args []string) (any, error) {
	result := strings.TrimPrefix(args[0], args[1])

	return ast.String(result), nil
}

// (to-lower SOURCE:String)
func toLowerFunction(ctx types.Context, args []string) (any, error) {
	result := strings.ToLower(args[0])

	return ast.String(result), nil
}

// (to-upper SOURCE:String)
func toUpperFunction(ctx types.Context, args []string) (any, error) {
	result := strings.ToUpper(args[0])

	return ast.String(result), nil
}

// (trim SOURCE:String)
func trimFunction(ctx types.Context, args []string) (any, error) {
	result := strings.TrimSpace(args[0])

	return ast.String(result), nil
}
