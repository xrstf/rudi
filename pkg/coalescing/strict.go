// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package coalescing

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

type strict struct{}

func NewStrict() Coalescer {
	return strict{}
}

var _ Coalescer = strict{}

func (strict) ToNull(val any) (bool, error) {
	switch v := val.(type) {
	case nil:
		return true, nil
	default:
		return false, fmt.Errorf("cannot coalesce %T into null", v)
	}
}

func (strict) ToBool(val any) (bool, error) {
	switch v := val.(type) {
	case nil:
		return false, nil
	case bool:
		return v, nil
	default:
		return false, fmt.Errorf("cannot coalesce %T into bool", v)
	}
}

func (strict) ToFloat64(val any) (float64, error) {
	switch v := val.(type) {
	case nil:
		return 0, nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("cannot coalesce %T into float64", v)
	}
}

func (strict) ToInt64(val any) (int64, error) {
	switch v := val.(type) {
	case nil:
		return 0, nil
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float32:
		if v == float32(int32(v)) {
			return int64(v), nil
		}
		return 0, fmt.Errorf("cannot convert %s losslessly to int64", formatFloat(float64(v)))
	case float64:
		if v == float64(int64(v)) {
			return int64(v), nil
		}
		return 0, fmt.Errorf("cannot convert %s losslessly to int64", formatFloat(v))
	default:
		return 0, fmt.Errorf("cannot coalesce %T into int64", v)
	}
}

func (s strict) ToNumber(val any) (ast.Number, error) {
	return toNumber(s, val)
}

func (strict) ToString(val any) (string, error) {
	switch v := val.(type) {
	case nil:
		return "", nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("cannot coalesce %T into string", v)
	}
}

func (strict) ToVector(val any) ([]any, error) {
	switch v := val.(type) {
	case nil:
		return []any{}, nil
	case []any:
		return v, nil
	default:
		return nil, fmt.Errorf("cannot coalesce %T into vector", v)
	}
}

func (strict) ToObject(val any) (map[string]any, error) {
	switch v := val.(type) {
	case nil:
		return map[string]any{}, nil
	case map[string]any:
		return v, nil
	default:
		return nil, fmt.Errorf("cannot coalesce %T into object", v)
	}
}
