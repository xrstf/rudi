// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package coalescing

import (
	"fmt"
	"strconv"
)

type literal interface {
	LiteralValue() any
}

func ToBool(val any) (bool, error) {
	lit, ok := val.(literal)
	if !ok {
		return false, fmt.Errorf("cannot coalesce %T into bool", val)
	}

	litVal := lit.LiteralValue()

	switch v := litVal.(type) {
	case bool:
		return v, nil
	case int64:
		return v != 0, nil
	case float64:
		return v != 0, nil
	case string:
		return len(v) > 0, nil
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("cannot coalesce %T into bool", litVal)
	}
}

func ToFloat64(val any) (float64, error) {
	lit, ok := val.(literal)
	if !ok {
		return 0, fmt.Errorf("cannot coalesce %T into float64", val)
	}

	litVal := lit.LiteralValue()

	switch v := litVal.(type) {
	case bool:
		if v {
			return 1, nil
		} else {
			return 0, nil
		}
	case int64:
		return float64(v), nil
	case float64:
		return v, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot coalesce %T into float64", litVal)
	}
}

func Int64Compatible(val any) bool {
	lit, ok := val.(literal)
	if !ok {
		return false
	}

	switch lit.LiteralValue().(type) {
	case bool:
		return true
	case int64:
		return true
	case nil:
		return true
	default:
		return false
	}
}

func ToInt64(val any) (int64, error) {
	lit, ok := val.(literal)
	if !ok {
		return 0, fmt.Errorf("cannot coalesce %T into int64", val)
	}

	litVal := lit.LiteralValue()

	switch v := litVal.(type) {
	case bool:
		if v {
			return 1, nil
		} else {
			return 0, nil
		}
	case int64:
		return v, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot coalesce %T into int64", litVal)
	}
}

func ToString(val any) (string, error) {
	lit, ok := val.(literal)
	if !ok {
		return "", fmt.Errorf("cannot coalesce %T into string", val)
	}

	litVal := lit.LiteralValue()

	switch v := litVal.(type) {
	case bool:
		return strconv.FormatBool(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float64:
		return fmt.Sprintf("%f", v), nil
	case nil:
		return "null", nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("cannot coalesce %T into string", litVal)
	}
}
