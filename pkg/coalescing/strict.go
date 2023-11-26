// SPDX-FileCopyrightText: 2023 Christoph Mewes
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
	switch v := deliteral(val).(type) {
	case nil:
		return true, nil
	default:
		return false, fmt.Errorf("cannot coalesce %T into null", v)
	}
}

func (strict) ToBool(val any) (bool, error) {
	switch v := deliteral(val).(type) {
	case bool:
		return v, nil
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("cannot coalesce %T into bool", v)
	}
}

func (strict) ToFloat64(val any) (float64, error) {
	switch v := deliteral(val).(type) {
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
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot coalesce %T into float64", v)
	}
}

func (strict) ToInt64(val any) (int64, error) {
	switch v := deliteral(val).(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case ast.Number:
		intVal, ok := v.ToInteger()
		if !ok {
			return 0, fmt.Errorf("cannot convert %f losslessly to int64", val)
		}
		return intVal, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot coalesce %T into int64", v)
	}
}

func (s strict) ToNumber(val any) (ast.Number, error) {
	return toNumber(s, val)
}

func (strict) ToString(val any) (string, error) {
	switch v := deliteral(val).(type) {
	case string:
		return v, nil
	case nil:
		return "", nil
	default:
		return "", fmt.Errorf("cannot coalesce %T into string", v)
	}
}

func (strict) ToVector(val any) ([]any, error) {
	switch v := deliteral(val).(type) {
	case []any:
		return v, nil
	case nil:
		return []any{}, nil
	default:
		return nil, fmt.Errorf("cannot coalesce %T into vector", v)
	}
}

func (strict) ToObject(val any) (map[string]any, error) {
	switch v := deliteral(val).(type) {
	case map[string]any:
		return v, nil
	case nil:
		return map[string]any{}, nil
	default:
		return nil, fmt.Errorf("cannot coalesce %T into object", v)
	}
}
