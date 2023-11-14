// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"
	"strings"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval"
	"go.xrstf.de/otto/pkg/lang/eval/coalescing"
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

	glueString, err := coalescing.ToString(glue)
	if err != nil {
		return nil, fmt.Errorf("glue is not stringish: %w", err)
	}

	values, err := evalArgs(ctx, args, 1)
	if err != nil {
		return nil, err
	}

	parts := []string{}
	for i, list := range values {
		vector, ok := list.(ast.Vector)
		if !ok {
			part, err := coalescing.ToString(list)
			if err != nil {
				return nil, fmt.Errorf("argument #%d is neither vector nor stringish: %w", i+1, err)
			}

			parts = append(parts, part)
			continue
		}

		for j, item := range vector.Data {
			part, err := coalescing.ToString(item)
			if err != nil {
				return nil, fmt.Errorf("argument #%d.%d is not stringish: %w", i+1, j, err)
			}
			parts = append(parts, part)
		}
	}

	return ast.String(strings.Join(parts, glueString)), nil
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

	sep, err := coalescing.ToString(values[0])
	if err != nil {
		return nil, fmt.Errorf("separator is not stringish: %w", err)
	}

	source, err := coalescing.ToString(values[1])
	if err != nil {
		return nil, fmt.Errorf("source is not stringish: %w", err)
	}

	parts := strings.Split(source, sep)
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
