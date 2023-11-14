// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package coalescing

import (
	"fmt"
	"strconv"
	"strings"

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
		if lower := strings.ToLower(v); lower == "" || lower == "false" || lower == "0" {
			return false, nil
		}

		return true, nil
	case []any:
		return len(v) > 0, nil
	case map[string]any:
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
	case string:
		parsed, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot parse %q as float: %w", v, err)
		}
		return parsed, nil
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
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot parse %q as int: %w", v, err)
		}
		return parsed, nil
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
		return formatFloat(v), nil
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

func formatFloat(f float64) string {
	formatted := fmt.Sprintf("%f", f)
	for strings.HasSuffix(formatted, "0") {
		formatted = strings.TrimSuffix(formatted, "0")
	}

	return strings.TrimSuffix(formatted, ".")
}

func typeName(v any) string {
	switch asserted := v.(type) {
	case ast.Expression:
		return asserted.ExpressionName()
	case map[string]interface{}:
		return "object" // lowercase, to distinguish from Objects and ObjectNodes
	case []interface{}:
		return "vector"
	default:
		return fmt.Sprintf("%T", v)
	}
}

func IsEmpty(val any) (bool, error) {
	switch v := val.(type) {
	case bool:
		return v == false, nil
	case int64:
		return v == 0, nil
	case float64:
		return v == 0, nil
	case nil:
		return true, nil
	case string:
		return len(v) == 0, nil
	case []any:
		return len(v) == 0, nil
	case map[string]any:
		return len(v) == 0, nil
	default:
		lit, ok := val.(literal)
		if !ok {
			return false, fmt.Errorf("cannot termine emptiness of %s", typeName(val))
		}

		return IsEmpty(lit.LiteralValue())
	}
}
