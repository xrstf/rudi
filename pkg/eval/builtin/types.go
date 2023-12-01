// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

// (to-string VAL:any)
func toStringFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, value, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	// this function purposefully always uses humane coalescing
	return coalescing.NewHumane().ToString(value)
}

// (to-int VAL:any)
func toIntFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, value, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	// this function purposefully always uses humane coalescing
	return coalescing.NewHumane().ToInt64(value)
}

// (to-float VAL:any)
func toFloatFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, value, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	// this function purposefully always uses humane coalescing
	return coalescing.NewHumane().ToFloat64(value)
}

// (to-bool VAL:any)
func toBoolFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, value, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	// this function purposefully always uses humane coalescing
	return coalescing.NewHumane().ToBool(value)
}

// (type-of VAL:any)
func typeOfFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, value, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

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
