// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package functions

import (
	"fmt"
	"reflect"

	"go.xrstf.de/rudi/pkg/runtime/types"
)

type form struct {
	fun any

	matcher *argsMatcher

	// args is filled when calling Match(), so than when calling Call(), we do not have to
	// re-coalesce all the arguments again.
	args []any
}

func newForm(fun any) (form, error) {
	funValue := reflect.ValueOf(fun)
	if funValue.Kind() != reflect.Func {
		return form{}, fmt.Errorf("given value is not a function, but %T", fun)
	}

	matcher, err := newArgsMatcher(fun)
	if err != nil {
		return form{}, err
	}

	return form{
		fun:     fun,
		matcher: matcher,
	}, nil
}

func (f *form) Match(ctx types.Context, args []cachedExpression) (bool, error) {
	consumed, matched, err := f.matcher.Match(ctx, args)
	if err != nil {
		return false, err
	}

	if !matched {
		return false, nil
	}

	f.args = consumed

	return true, nil
}

func (f *form) Call(ctx types.Context) (any, error) {
	reflectArgs := make([]reflect.Value, len(f.args))
	for i, arg := range f.args {
		if arg == nil {
			var e any
			reflectArgs[i] = reflect.ValueOf(&e).Elem()
		} else {
			reflectArgs[i] = reflect.ValueOf(f.args[i])
		}
	}

	results := reflect.ValueOf(f.fun).Call(reflectArgs)

	// Forms can only be constructed with valid signatures,
	// no need to check that 2 values were returned.
	if err := results[1].Interface(); err != nil {
		return nil, err.(error)
	}

	return results[0].Interface(), nil
}
