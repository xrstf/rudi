// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package coalescing

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

type CustomNullCoalescer interface {
	CoalesceToNull(c Coalescer) (bool, error)
}

type CustomBoolCoalescer interface {
	CoalesceToBool(c Coalescer) (bool, error)
}

type CustomInt64Coalescer interface {
	CoalesceToInt64(c Coalescer) (int64, error)
}

type CustomFloat64Coalescer interface {
	CoalesceToFloat64(c Coalescer) (float64, error)
}

type CustomStringCoalescer interface {
	CoalesceToString(c Coalescer) (string, error)
}

type CustomVectorCoalescer interface {
	CoalesceToVector(c Coalescer) ([]any, error)
}

type CustomObjectCoalescer interface {
	CoalesceToObject(c Coalescer) (map[string]any, error)
}

type humane struct{}

func NewHumane() Coalescer {
	return humane{}
}

var _ Coalescer = humane{}

func (h humane) ToNull(val any) (bool, error) {
	switch v := val.(type) {
	case nil:
		return true, nil
	case bool:
		if v {
			return false, fmt.Errorf("cannot coalesce true into null")
		}
		return true, nil
	case int:
		if v != 0 {
			return false, fmt.Errorf("cannot coalesce %v (%T) into null", v, v)
		}
		return true, nil
	case int32:
		if v != 0 {
			return false, fmt.Errorf("cannot coalesce %v (%T) into null", v, v)
		}
		return true, nil
	case int64:
		if v != 0 {
			return false, fmt.Errorf("cannot coalesce %v (%T) into null", v, v)
		}
		return true, nil
	case float32:
		if v != 0 {
			return false, fmt.Errorf("cannot coalesce %v (%T) into null", v, v)
		}
		return true, nil
	case float64:
		if v != 0 {
			return false, fmt.Errorf("cannot coalesce %v (%T) into null", v, v)
		}
		return true, nil
	case string:
		if len(v) != 0 {
			return false, fmt.Errorf("cannot coalesce %q (%T) into null", v, v)
		}
		return true, nil
	case []any:
		if len(v) != 0 {
			return false, errors.New("cannot coalesce non-empty vector into null")
		}
		return true, nil
	case map[string]any:
		if len(v) != 0 {
			return false, errors.New("cannot coalesce non-empty object into null")
		}
		return true, nil
	case CustomNullCoalescer:
		return v.CoalesceToNull(h)
	default:
		return false, fmt.Errorf("cannot coalesce %T into null", v)
	}
}

func (h humane) ToBool(val any) (bool, error) {
	switch v := val.(type) {
	case nil:
		return false, nil
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
	case CustomBoolCoalescer:
		return v.CoalesceToBool(h)
	default:
		return false, fmt.Errorf("cannot coalesce %T into bool", v)
	}
}

func (h humane) ToFloat64(val any) (float64, error) {
	switch v := val.(type) {
	case nil:
		return 0, nil
	case bool:
		if v {
			return 1, nil
		} else {
			return 0, nil
		}
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
		v = strings.TrimSpace(v)
		if v == "" {
			return 0, nil
		}

		parsed, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot coalesce %T into float64", v)
		}
		return parsed, nil
	case CustomFloat64Coalescer:
		return v.CoalesceToFloat64(h)
	default:
		return 0, fmt.Errorf("cannot coalesce %T into float64", v)
	}
}

func (h humane) ToInt64(val any) (int64, error) {
	switch v := val.(type) {
	case nil:
		return 0, nil
	case bool:
		if v {
			return 1, nil
		} else {
			return 0, nil
		}
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
	case string:
		v = strings.TrimSpace(v)
		if v == "" {
			return 0, nil
		}

		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			// allows "2.0" to turn into int64(2)
			parsed, err := strconv.ParseFloat(v, 64)
			if err == nil && parsed == float64(int64(parsed)) {
				return int64(parsed), nil
			}

			return 0, fmt.Errorf("cannot coalesce %T into int64", v)
		}
		return parsed, nil
	case CustomInt64Coalescer:
		return v.CoalesceToInt64(h)
	default:
		return 0, fmt.Errorf("cannot coalesce %T into int64", v)
	}
}

func (h humane) ToNumber(val any) (ast.Number, error) {
	return toNumber(h, val)
}

func (h humane) ToString(val any) (string, error) {
	switch v := val.(type) {
	case nil:
		return "", nil
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
	case string:
		return v, nil
	case CustomStringCoalescer:
		return v.CoalesceToString(h)
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

func (h humane) ToVector(val any) ([]any, error) {
	switch v := val.(type) {
	case nil:
		return []any{}, nil
	case []any:
		return v, nil
	case map[string]any:
		if len(v) == 0 {
			return []any{}, nil
		} else {
			return nil, fmt.Errorf("cannot coalesce %T into vector", v)
		}
	case CustomVectorCoalescer:
		return v.CoalesceToVector(h)
	default:
		return nil, fmt.Errorf("cannot coalesce %T into vector", v)
	}
}

func (h humane) ToObject(val any) (map[string]any, error) {
	switch v := val.(type) {
	case nil:
		return map[string]any{}, nil
	case []any:
		if len(v) == 0 {
			return map[string]any{}, nil
		} else {
			return nil, fmt.Errorf("cannot coalesce %T into object", v)
		}
	case map[string]any:
		return v, nil
	case CustomObjectCoalescer:
		return v.CoalesceToObject(h)
	default:
		return nil, fmt.Errorf("cannot coalesce %T into object", v)
	}
}
