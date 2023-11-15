// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package coalescing

import (
	"fmt"
	"strconv"
	"strings"

	"go.xrstf.de/otto/pkg/lang/ast"
)

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
		lit, ok := val.(ast.Literal)
		if !ok {
			return false, fmt.Errorf("cannot coalesce %T into bool", val)
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
		lit, ok := val.(ast.Literal)
		if !ok {
			return 0, fmt.Errorf("cannot coalesce %T into float64", val)
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
		lit, ok := val.(ast.Literal)
		if !ok {
			return 0, fmt.Errorf("cannot coalesce %T into int64", val)
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
		lit, ok := val.(ast.Literal)
		if !ok {
			return "", fmt.Errorf("cannot coalesce %T into string", val)
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
		lit, ok := val.(ast.Literal)
		if !ok {
			return false, fmt.Errorf("cannot termine emptiness oT %s", val)
		}

		return IsEmpty(lit.LiteralValue())
	}
}
