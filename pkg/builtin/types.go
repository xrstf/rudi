// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/eval/types"
)

// (to-string VAL:any)
func toStringFunction(ctx types.Context, args []any) (any, error) {
	// this function purposefully always uses humane coalescing
	return coalescing.NewHumane().ToString(args[0])
}

// (to-int VAL:any)
func toIntFunction(ctx types.Context, args []any) (any, error) {
	// this function purposefully always uses humane coalescing
	return coalescing.NewHumane().ToInt64(args[0])
}

// (to-float VAL:any)
func toFloatFunction(ctx types.Context, args []any) (any, error) {
	// this function purposefully always uses humane coalescing
	return coalescing.NewHumane().ToFloat64(args[0])
}

// (to-bool VAL:any)
func toBoolFunction(ctx types.Context, args []any) (any, error) {
	// this function purposefully always uses humane coalescing
	return coalescing.NewHumane().ToBool(args[0])
}

// (type-of VAL:any)
func typeOfFunction(ctx types.Context, args []any) (any, error) {
	var typeName string

	switch val := args[0].(type) {
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
		typeName = fmt.Sprintf("%T", val)
	}

	return typeName, nil
}
