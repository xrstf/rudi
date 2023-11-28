//go:build integration

// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package php

import (
	"fmt"
	"strings"
)

// Convert converts a Go value to a rough PHP representation.
func Convert(val any) string {
	switch asserted := val.(type) {
	case nil:
		return `null`
	case bool:
		return fmt.Sprintf("(bool) %v", asserted)
	case int:
		return fmt.Sprintf("(int) %d", asserted)
	case int32:
		return fmt.Sprintf("(int) %d", asserted)
	case int64:
		return fmt.Sprintf("(int) %d", asserted)
	case float32:
		return fmt.Sprintf("(double) %f", asserted)
	case float64:
		return fmt.Sprintf("(double) %f", asserted)
	case string:
		return fmt.Sprintf("%q", asserted)
	case []any:
		items := make([]string, len(asserted))
		for i, item := range asserted {
			phpItem := Convert(item)
			if phpItem == "" {
				return ""
			}

			items[i] = phpItem
		}

		return fmt.Sprintf(`[%s]`, strings.Join(items, ","))

	case map[string]any:
		items := make([]string, len(asserted))
		i := 0
		for key, value := range asserted {
			phpValue := Convert(value)
			if phpValue == "" {
				return ""
			}

			items[i] = fmt.Sprintf(`%q => %s`, key, phpValue)
			i++
		}

		return fmt.Sprintf(`[%s]`, strings.Join(items, ","))
	default:
		return ""
	}
}
