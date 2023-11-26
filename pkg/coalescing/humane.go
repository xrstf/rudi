// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package coalescing

import (
	"fmt"
	"strconv"
	"strings"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

type humane struct{}

func NewHumane() Coalescer {
	return humane{}
}

var _ Coalescer = humane{}

func (humane) ToNull(val any) (bool, error) {
	switch v := deliteral(val).(type) {
	case nil:
		return true, nil
	case bool:
		return !v, nil
	default:
		return false, fmt.Errorf("cannot coalesce %T into null", v)
	}
}

func (humane) ToBool(val any) (bool, error) {
	switch v := deliteral(val).(type) {
	case bool:
		return v, nil
	case int:
		return v != 0, nil
	case int32:
		return v != 0, nil
	case int64:
		return v != 0, nil
	case float32:
		return v != 0, nil
	case float64:
		return v != 0, nil
	case string:
		if v == "" || v == "0" {
			return false, nil
		}

		return !strings.EqualFold(v, "false"), nil
	case []any:
		return len(v) > 0, nil
	case map[string]any:
		return len(v) > 0, nil
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("cannot coalesce %T into bool", v)
	}
}

func (humane) ToFloat64(val any) (float64, error) {
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
	case string:
		parsed, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert %q losslessly to float64", v)
		}
		return parsed, nil
	case bool:
		if v {
			return 1, nil
		} else {
			return 0, nil
		}
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot coalesce %T into float64", v)
	}
}

func (humane) ToInt64(val any) (int64, error) {
	switch v := deliteral(val).(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert %q losslessly to int64", v)
		}
		return parsed, nil
	case bool:
		if v {
			return 1, nil
		} else {
			return 0, nil
		}
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot coalesce %T into int64", v)
	}
}

func (h humane) ToNumber(val any) (ast.Number, error) {
	return toNumber(h, val)
}

func (humane) ToString(val any) (string, error) {
	switch v := deliteral(val).(type) {
	case string:
		return v, nil
	case bool:
		return strconv.FormatBool(v), nil
	case int:
		return strconv.FormatInt(int64(v), 10), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float64:
		return formatFloat(v), nil
	case nil:
		return "", nil
	default:
		return "", fmt.Errorf("cannot coalesce %T into string", v)
	}
}

func formatFloat(f float64) string {
	formatted := fmt.Sprintf("%f", f)
	for strings.HasSuffix(formatted, "0") {
		formatted = strings.TrimSuffix(formatted, "0")
	}

	return strings.TrimSuffix(formatted, ".")
}

func (humane) ToVector(val any) ([]any, error) {
	switch v := deliteral(val).(type) {
	case nil:
		return []any{}, nil
	case []any:
		return v, nil
	default:
		return nil, fmt.Errorf("cannot coalesce %T into vector", v)
	}
}

func (humane) ToObject(val any) (map[string]any, error) {
	switch v := deliteral(val).(type) {
	case nil:
		return map[string]any{}, nil
	case map[string]any:
		return v, nil
	default:
		return nil, fmt.Errorf("cannot coalesce %T into object", v)
	}
}
