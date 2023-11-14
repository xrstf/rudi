// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"
	"strings"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval"
	"go.xrstf.de/otto/pkg/lang/eval/coalescing"
	"go.xrstf.de/otto/pkg/lang/eval/types"
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

	coalesced, err := coalescing.ToString(value)
	if err != nil {
		return nil, err
	}

	return ast.String(coalesced), nil
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

	coalesced, err := coalescing.ToInt64(value)
	if err != nil {
		return nil, err
	}

	return ast.Number{Value: coalesced}, nil
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

	coalesced, err := coalescing.ToFloat64(value)
	if err != nil {
		return nil, err
	}

	return ast.Number{Value: coalesced}, nil
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

	coalesced, err := coalescing.ToBool(value)
	if err != nil {
		return nil, err
	}

	return ast.Bool(coalesced), nil
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

	// EvalExpression() will return a number or bool, but those do count as expressions
	expr, ok := value.(ast.Expression)
	if !ok {
		return nil, fmt.Errorf("expected expression, but got %T", value)
	}

	name := strings.ToLower(expr.ExpressionName())

	return ast.String(name), nil
}
