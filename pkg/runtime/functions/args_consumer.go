// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package functions

import (
	"reflect"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

var (
	dummyExpression ast.Expression

	expressionType = reflect.TypeOf(&dummyExpression).Elem()
	contextType    = reflect.TypeOf(types.Context{})
	numberType     = reflect.TypeOf(ast.Number{})
)

type argsConsumer func(ctx types.Context, args []cachedExpression) (asserted []any, remaining []cachedExpression, err error)

func isAny(t reflect.Type) bool {
	return t.Kind() == reflect.Interface && t.Name() == ""
}

func newConsumerFunc(t reflect.Type) (consumer argsConsumer, argsConsumed int) {
	switch t.Kind() {
	case reflect.Bool:
		return boolConsumer, 1
	case reflect.Int64:
		return intConsumer, 1
	case reflect.Float64:
		return floatConsumer, 1
	case reflect.String:
		return stringConsumer, 1
	case reflect.Slice:
		// we only support []any
		if isAny(t.Elem()) {
			return vectorConsumer, 1
		}

	case reflect.Map:
		// we only support map[string]any
		// TODO: Check key as well.
		if isAny(t.Elem()) {
			return objectConsumer, 1
		}

	case reflect.Interface:
		// empty interface (any)
		if isAny(t) {
			return anyConsumer, 1
		}

		// allow unevaluated access to the argument expression
		if t.AssignableTo(expressionType) {
			return expressionConsumer, 1
		}

	case reflect.Struct:
		// allow to inject the context when required
		if t.AssignableTo(contextType) {
			// this consumer does not consume an argument expression
			return contextConsumer, 0
		}

		if t.AssignableTo(numberType) {
			return numberConsumer, 1
		}
	}

	return nil, 0
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

func numberConsumer(ctx types.Context, args []cachedExpression) (asserted []any, remaining []cachedExpression, err error) {
	if len(args) == 0 {
		return nil, nil, nil
	}

	evaluated, err := args[0].Eval(ctx)
	if err != nil {
		return nil, nil, err
	}

	coalesced, err := ctx.Coalesce().ToNumber(evaluated)
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

func contextConsumer(ctx types.Context, args []cachedExpression) (asserted []any, remaining []cachedExpression, err error) {
	return []any{ctx}, args, nil
}

// toVariadicConsumer wraps a singular consumer to consume all remaining args. In contrast to
// Go, variadic arguments must have at least 1 item (i.e. calling func(foo string, a ...int) with
// ("abc") only is invalid).
func toVariadicConsumer(singleConsumer argsConsumer) argsConsumer {
	return func(ctx types.Context, args []cachedExpression) (asserted []any, remaining []cachedExpression, err error) {
		// at least one argument is required, otherwise variadics do not match
		if len(args) == 0 {
			return nil, nil, nil
		}

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
