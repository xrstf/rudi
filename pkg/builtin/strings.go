// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"
	"strings"

	"go.xrstf.de/rudi/pkg/eval/types"
)

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

	return strings.Join(parts, string(glue)), nil
}

func splitFunction(sep string, source string) (any, error) {
	parts := strings.Split(source, sep)

	// to []any
	result := make([]any, len(parts))
	for i, part := range parts {
		result[i] = part
	}

	return result, nil
}

func splitnFunction(sep string, source string, limit int64) (any, error) {
	parts := strings.SplitN(source, sep, int(limit))

	// to []any
	result := make([]any, len(parts))
	for i, part := range parts {
		result[i] = part
	}

	return result, nil
}

func hasSuffixFunction(source string, suffix string) (any, error) {
	return strings.HasSuffix(source, suffix), nil
}

func hasPrefixFunction(source string, prefix string) (any, error) {
	return strings.HasPrefix(source, prefix), nil
}

func trimSuffixFunction(source string, suffix string) (any, error) {
	return strings.TrimSuffix(source, suffix), nil
}

func trimPrefixFunction(source string, prefix string) (any, error) {
	return strings.TrimPrefix(source, prefix), nil
}

func toLowerFunction(s string) (any, error) {
	return strings.ToLower(s), nil
}

func toUpperFunction(s string) (any, error) {
	return strings.ToUpper(s), nil
}

func trimFunction(s string) (any, error) {
	return strings.TrimSpace(s), nil
}

func replaceAllFunction(s, old, new string) (any, error) {
	return strings.ReplaceAll(s, old, new), nil
}

func replaceLimitFunction(s, old, new string, limit int64) (any, error) {
	return strings.Replace(s, old, new, int(limit)), nil
}
