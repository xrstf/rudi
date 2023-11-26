// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"
	"strings"

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

	value, err = types.WrapNative(value)
	if err != nil {
		return nil, err
	}

	expr, ok := value.(ast.Literal)
	if !ok {
		return nil, fmt.Errorf("expected expression, but got %T", value)
	}

	name := strings.ToLower(expr.ExpressionName())

	return name, nil
}
