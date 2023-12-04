// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package native

import (
	"errors"
	"fmt"
	"reflect"

	"go.xrstf.de/rudi/pkg/eval/types"
)

type argsMatcher struct {
	consumers []argsConsumer
}

// newArgsMatcher creates a new argsMatcher for a given function. The function must follow the
// allowed signature (i.e. return (any, error) etc.) and only use allowed parameter types.
func newArgsMatcher(fun any) (*argsMatcher, error) {
	funType := reflect.TypeOf(fun)

	// check return signature, must match (any, error)
	if err := checkReturnValueSignature(funType); err != nil {
		return nil, err
	}

	// create a list of consumers that will match the function's signature
	consumers, err := createConsumers(funType)
	if err != nil {
		return nil, err
	}

	return &argsMatcher{
		consumers: consumers,
	}, nil
}

// checkReturnValueSignature ensures that the function returns (any, error).
func checkReturnValueSignature(funType reflect.Type) error {
	if retvals := funType.NumOut(); retvals != 2 {
		return fmt.Errorf("function must return (any, error), but has %d return values", retvals)
	}
	if o := funType.Out(0); !isAny(o) {
		return errors.New("function must return (any, error)")
	}
	if o := funType.Out(1); o.Kind() != reflect.Interface || o.Name() != "error" {
		return errors.New("function must return (any, error)")
	}

	return nil
}

// createConsumers converts each function parameter into a consumer, returning the stack of
// consumers.
func createConsumers(funType reflect.Type) ([]argsConsumer, error) {
	consumers := []argsConsumer{}
	variadic := funType.IsVariadic()
	totalParams := funType.NumIn()

	// For each of the function's parameters, create a consumer function that
	// attempts to read and coalesce an argument into the desired parameter type.
	for i := 0; i < totalParams; i++ {
		parameterType := funType.In(i)
		variadicArg := variadic && i == totalParams-1

		// A variadic function like func(int, ...int) has its last parameter
		// report []int as the type; we "unwrap" this here to get a plain int
		// consumer first.
		if variadicArg {
			parameterType = parameterType.Elem()
		}

		consumer := newConsumerFunc(parameterType)
		if consumer == nil {
			return nil, fmt.Errorf("cannot handle %v parameters", parameterType)
		}

		// Wrap the single consumer into a variadic consumer
		// that just keeps consuming until all args are gone.
		if variadicArg {
			consumer = toVariadicConsumer(consumer)
		}

		consumers = append(consumers, consumer)
	}

	return consumers, nil
}

func (c *argsMatcher) Match(ctx types.Context, args []cachedExpression) ([]any, bool, error) {
	result := []any{}
	remaining := args

	// Run each consumer func in succession, making each consume as many args as it wants.
	for _, consumer := range c.consumers {
		var (
			consumed []any
			err      error
		)

		consumed, remaining, err = consumer(ctx, remaining)
		if err != nil {
			return nil, false, err
		}

		// the consumer didn't match
		if consumed == nil {
			return nil, false, nil
		}

		result = append(result, consumed...)
	}

	// not all arguments consumed => no match
	if len(remaining) > 0 {
		return nil, false, nil
	}

	return result, true, nil
}
