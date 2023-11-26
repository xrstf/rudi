// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package coalescing

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

type pedantic struct{}

func NewPedantic() Coalescer {
	return pedantic{}
}

var _ Coalescer = pedantic{}

func (pedantic) ToNull(val any) (bool, error) {
	switch v := deliteral(val).(type) {
	case nil:
		return true, nil
	default:
		return false, fmt.Errorf("cannot coalesce %T into null", v)
	}
}

func (pedantic) ToBool(val any) (bool, error) {
	switch v := deliteral(val).(type) {
	case bool:
		return v, nil
	default:
		return false, fmt.Errorf("cannot coalesce %T into bool", v)
	}
}

func (pedantic) ToFloat64(val any) (float64, error) {
	switch v := deliteral(val).(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("cannot coalesce %T into float64", v)
	}
}

func (pedantic) ToInt64(val any) (int64, error) {
	switch v := deliteral(val).(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	default:
		return 0, fmt.Errorf("cannot coalesce %T into int64", v)
	}
}

func (p pedantic) ToNumber(val any) (ast.Number, error) {
	return toNumber(p, val)
}

func (pedantic) ToString(val any) (string, error) {
	switch v := deliteral(val).(type) {
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("cannot coalesce %T into string", v)
	}
}

func (pedantic) ToVector(val any) ([]any, error) {
	switch v := deliteral(val).(type) {
	case []any:
		return v, nil
	default:
		return nil, fmt.Errorf("cannot coalesce %T into vector", v)
	}
}

func (pedantic) ToObject(val any) (map[string]any, error) {
	switch v := deliteral(val).(type) {
	case map[string]any:
		return v, nil
	default:
		return nil, fmt.Errorf("cannot coalesce %T into object", v)
	}
}
