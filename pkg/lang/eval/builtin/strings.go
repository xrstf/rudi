// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"
	"strings"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval"
	"go.xrstf.de/otto/pkg/lang/eval/types"
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

// (split SEP:String SOURCE:String)
func splitFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 2 {
		return nil, fmt.Errorf("expected 2 arguments, got %d", size)
	}

	values, err := evalArgs(ctx, args, 0)
	if err != nil {
		return nil, err
	}

	sep, ok := values[0].(ast.String)
	if !ok {
		return nil, fmt.Errorf("separator is not a string, but %T", values[0])
	}

	source, ok := values[1].(ast.String)
	if !ok {
		return nil, fmt.Errorf("source is not a string, but %T", values[1])
	}

	parts := strings.Split(string(source), string(sep))
	result := make([]any, len(parts))
	for i, part := range parts {
		result[i] = ast.String(part)
	}

	return ast.Vector{Data: result}, nil
}

// (trim-suffix SUFFIX:String SOURCE:String)
func trimSuffixFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 2 {
		return nil, fmt.Errorf("expected 2 arguments, got %d", size)
	}

	_, suffix, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	suffixStr, ok := suffix.(ast.String)
	if !ok {
		return nil, fmt.Errorf("argument #0 is not a string, but %T", suffix)
	}

	_, source, err := eval.EvalExpression(ctx, args[1])
	if err != nil {
		return nil, fmt.Errorf("argument #1: %w", err)
	}

	sourceStr, ok := source.(ast.String)
	if !ok {
		return nil, fmt.Errorf("argument #1 is not a string, but %T", source)
	}

	result := strings.TrimSuffix(string(sourceStr), string(suffixStr))

	return ast.String(result), nil
}

// (trim-prefix PREFIX:String SOURCE:String)
func trimPrefixFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 2 {
		return nil, fmt.Errorf("expected 2 arguments, got %d", size)
	}

	_, prefix, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	prefixStr, ok := prefix.(ast.String)
	if !ok {
		return nil, fmt.Errorf("argument #0 is not a string, but %T", prefix)
	}

	_, source, err := eval.EvalExpression(ctx, args[1])
	if err != nil {
		return nil, fmt.Errorf("argument #1: %w", err)
	}

	sourceStr, ok := source.(ast.String)
	if !ok {
		return nil, fmt.Errorf("argument #1 is not a string, but %T", source)
	}

	result := strings.TrimPrefix(string(sourceStr), string(prefixStr))

	return ast.String(result), nil
}

// (to-lower SOURCE:String)
func toLowerFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, prefix, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	str, ok := prefix.(ast.String)
	if !ok {
		return nil, fmt.Errorf("argument is not a string, but %T", prefix)
	}

	result := strings.ToLower(string(str))

	return ast.String(result), nil
}

// (to-upper SOURCE:String)
func toUpperFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, prefix, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	str, ok := prefix.(ast.String)
	if !ok {
		return nil, fmt.Errorf("argument is not a string, but %T", prefix)
	}

	result := strings.ToUpper(string(str))

	return ast.String(result), nil
}
