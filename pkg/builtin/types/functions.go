// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package types

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/runtime/functions"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

var (
	humaneCoalescer = coalescing.NewHumane()

	Functions = types.Functions{
		"type-of": functions.NewBuilder(typeOfFunction).WithDescription(`returns the type of a given value (e.g. "string" or "number")`).Build(),

		// these functions purposefully always uses humane coalescing
		"to-bool":   functions.NewBuilder(toBoolFunction).WithCoalescer(humaneCoalescer).WithDescription("try to convert the given argument losslessly to a bool").Build(),
		"to-float":  functions.NewBuilder(toFloatFunction).WithCoalescer(humaneCoalescer).WithDescription("try to convert the given argument losslessly to a float64").Build(),
		"to-int":    functions.NewBuilder(toIntFunction).WithCoalescer(humaneCoalescer).WithDescription("try to convert the given argument losslessly to an int64").Build(),
		"to-string": functions.NewBuilder(toStringFunction).WithCoalescer(humaneCoalescer).WithDescription("try to convert the given argument losslessly to a string").Build(),
	}
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
