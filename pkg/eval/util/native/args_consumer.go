// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package native

import (
	"reflect"

	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

var (
	dummy          ast.Expression
	expressionType = reflect.TypeOf(&dummy).Elem()
)

type argsConsumer func(ctx types.Context, args []cachedExpression) (asserted []any, remaining []cachedExpression, err error)

func isAny(t reflect.Type) bool {
	return t.Kind() == reflect.Interface && t.Name() == ""
}

func newConsumerFunc(t reflect.Type) argsConsumer {
	switch t.Kind() {
	case reflect.Bool:
		return boolConsumer
	case reflect.Int64:
		return intConsumer
	case reflect.Float64:
		return floatConsumer
	case reflect.String:
		return stringConsumer
	case reflect.Slice:
		// we only support []any
		if isAny(t.Elem()) {
			return vectorConsumer
		}

	case reflect.Map:
		// we only support map[string]any
		// TODO: Check key as well.
		if isAny(t.Elem()) {
			return objectConsumer
		}

	case reflect.Interface:
		// empty interface (any)
		if isAny(t) {
			return anyConsumer
		}

		// allow unevaluated access to the argument expression
		if t.AssignableTo(expressionType) {
			return expressionConsumer
		}
	}

	return nil
}

func boolConsumer(ctx types.Context, args []cachedExpression) (asserted []any, remaining []cachedExpression, err error) {
	if len(args) == 0 {
		return nil, nil, nil
	}

	evaluated, err := args[0].Eval(ctx)
	if err != nil {
		return nil, nil, err
	}

	coalesced, err := ctx.Coalesce().ToBool(evaluated)
	if err != nil {
		return nil, args, nil
	}

	return []any{coalesced}, args[1:], nil
}

func intConsumer(ctx types.Context, args []cachedExpression) (asserted []any, remaining []cachedExpression, err error) {
	if len(args) == 0 {
		return nil, nil, nil
	}

	evaluated, err := args[0].Eval(ctx)
	if err != nil {
		return nil, nil, err
	}

	coalesced, err := ctx.Coalesce().ToInt64(evaluated)
	if err != nil {
		return nil, args, nil
	}

	return []any{coalesced}, args[1:], nil
}

func floatConsumer(ctx types.Context, args []cachedExpression) (asserted []any, remaining []cachedExpression, err error) {
	if len(args) == 0 {
		return nil, nil, nil
	}

	evaluated, err := args[0].Eval(ctx)
	if err != nil {
		return nil, nil, err
	}

	coalesced, err := ctx.Coalesce().ToFloat64(evaluated)
	if err != nil {
		return nil, args, nil
	}

	return []any{coalesced}, args[1:], nil
}

func stringConsumer(ctx types.Context, args []cachedExpression) (asserted []any, remaining []cachedExpression, err error) {
	if len(args) == 0 {
		return nil, nil, nil
	}

	evaluated, err := args[0].Eval(ctx)
	if err != nil {
		return nil, nil, err
	}

	coalesced, err := ctx.Coalesce().ToString(evaluated)
	if err != nil {
		return nil, args, nil
	}

	return []any{coalesced}, args[1:], nil
}

func vectorConsumer(ctx types.Context, args []cachedExpression) (asserted []any, remaining []cachedExpression, err error) {
	if len(args) == 0 {
		return nil, nil, nil
	}

	evaluated, err := args[0].Eval(ctx)
	if err != nil {
		return nil, nil, err
	}

	coalesced, err := ctx.Coalesce().ToVector(evaluated)
	if err != nil {
		return nil, args, nil
	}

	return []any{coalesced}, args[1:], nil
}

func objectConsumer(ctx types.Context, args []cachedExpression) (asserted []any, remaining []cachedExpression, err error) {
	if len(args) == 0 {
		return nil, nil, nil
	}

	evaluated, err := args[0].Eval(ctx)
	if err != nil {
		return nil, nil, err
	}

	coalesced, err := ctx.Coalesce().ToObject(evaluated)
	if err != nil {
		return nil, args, nil
	}

	return []any{coalesced}, args[1:], nil
}

func anyConsumer(ctx types.Context, args []cachedExpression) (asserted []any, remaining []cachedExpression, err error) {
	if len(args) == 0 {
		return nil, nil, nil
	}

	evaluated, err := args[0].Eval(ctx)
	if err != nil {
		return nil, nil, err
	}

	return []any{evaluated}, args[1:], nil
}

func expressionConsumer(ctx types.Context, args []cachedExpression) (asserted []any, remaining []cachedExpression, err error) {
	if len(args) == 0 {
		return nil, nil, nil
	}

	return []any{args[0].expr}, args[1:], nil
}

func toVariadicConsumer(singleConsumer argsConsumer) argsConsumer {
	return func(ctx types.Context, args []cachedExpression) (asserted []any, remaining []cachedExpression, err error) {
		leftover := args
		result := []any{}

		for len(leftover) > 0 {
			var (
				consumed []any
				err      error
			)

			consumed, leftover, err = singleConsumer(ctx, leftover)
			if err != nil {
				return nil, nil, err
			}

			// Variadic consumers must consume all args, so a noMatch is not allowed;
			// it's not an error though, we just have to signal that in total, the
			// matching was not a success.
			if consumed == nil {
				return nil, nil, nil
			}

			result = append(result, consumed[0])
		}

		return result, nil, nil
	}
}
