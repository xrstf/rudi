// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package coalescing

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

type Coalescer interface {
	ToNull(val any) (bool, error)
	ToBool(val any) (bool, error)
	ToInt64(val any) (int64, error)
	ToFloat64(val any) (float64, error)
	ToNumber(val any) (ast.Number, error)
	ToString(val any) (string, error)
	ToVector(val any) ([]any, error)
	ToObject(val any) (map[string]any, error)
}

func deliteral(val any) any {
	lit, ok := val.(ast.Literal)
	if ok {
		return lit.LiteralValue()
	}

	return val
}

func toNumber(c Coalescer, val any) (ast.Number, error) {
	i, err := c.ToInt64(val)
	if err == nil {
		return ast.Number{Value: i}, nil
	}

	f, err := c.ToFloat64(val)
	if err == nil {
		return ast.Number{Value: f}, nil
	}

	return ast.Number{}, fmt.Errorf("cannot convert %v losslessly to number", val)
}
