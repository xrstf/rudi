// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"
)

// The actual conversions happen in the pattern matching, since the functions
// are already configured to use the human coalescer for that process.

func toBoolFunction(b bool) (any, error) {
	return b, nil
}

func toIntFunction(i int64) (any, error) {
	return i, nil
}

func toFloatFunction(f float64) (any, error) {
	return f, nil
}

func toStringFunction(s string) (any, error) {
	return s, nil
}

func typeOfFunction(value any) (any, error) {
	var typeName string

	switch value.(type) {
	case nil:
		typeName = "null"
	case bool:
		typeName = "bool"
	case int64:
		typeName = "number"
	case float64:
		typeName = "number"
	case string:
		typeName = "string"
	case []any:
		typeName = "vector"
	case map[string]any:
		typeName = "object"
	default:
		// should never happen
		typeName = fmt.Sprintf("%T", value)
	}

	return typeName, nil
}
