// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package coalescing

import (
	"fmt"
	"strconv"

	"go.xrstf.de/otto/pkg/lang/ast"
)

type literal interface {
	LiteralValue() any
}

func ToBool(val any) (bool, error) {
	switch v := val.(type) {
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
		lit, ok := val.(literal)
		if !ok {
			return false, fmt.Errorf("cannot coalesce %s into bool", typeName(val))
		}

		return ToBool(lit.LiteralValue())
	}
}

func ToFloat64(val any) (float64, error) {
	switch v := val.(type) {
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
		lit, ok := val.(literal)
		if !ok {
			return 0, fmt.Errorf("cannot coalesce %s into float64", typeName(val))
		}

		return ToFloat64(lit.LiteralValue())
	}
}

func Int64Compatible(val any) bool {
	switch val.(type) {
	case bool:
		return true
	case int64:
		return true
	case nil:
		return true
	default:
		lit, ok := val.(literal)
		if !ok {
			return false
		}

		return Int64Compatible(lit.LiteralValue())
	}
}

func ToInt64(val any) (int64, error) {
	switch v := val.(type) {
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
		lit, ok := val.(literal)
		if !ok {
			return 0, fmt.Errorf("cannot coalesce %s into int64", typeName(val))
		}

		return ToInt64(lit.LiteralValue())
	}
}

func ToString(val any) (string, error) {
	switch v := val.(type) {
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
		lit, ok := val.(literal)
		if !ok {
			return "", fmt.Errorf("cannot coalesce %s into string", typeName(val))
		}

		return ToString(lit.LiteralValue())
	}
}

func typeName(v any) string {
	switch asserted := v.(type) {
	case ast.Node:
		return asserted.NodeName()
	case map[string]interface{}:
		return "object" // lowercase, to distinguish from Objects and ObjectNodes
	case []interface{}:
		return "vector"
	default:
		return fmt.Sprintf("%T", v)
	}
}
