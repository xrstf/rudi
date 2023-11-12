// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"
	"strings"

	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/coalescing"
	"go.xrstf.de/corel/pkg/lang/eval/types"
)

// (concat GLUE:String ELEMENTS:(Vector/String)+)
func concatFunction(ctx types.Context, args []Argument) (any, error) {
	if size := len(args); size < 2 {
		return nil, fmt.Errorf("expected 2+ arguments, got %d", size)
	}

	_, glue, err := args[0].Eval(ctx)
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

	return ast.String{Value: strings.Join(parts, glueString)}, nil
}

// (split SEP:String SOURCE:String)
func splitFunction(ctx types.Context, args []Argument) (any, error) {
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
		result[i] = ast.String{Value: part}
	}

	return ast.Vector{Data: result}, nil
}
