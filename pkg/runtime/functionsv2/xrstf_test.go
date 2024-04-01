// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package functionsv2

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func theFuncToCall(i int) (any, error) {
	fmt.Printf("i: %#v\n", i)

	return nil, nil
}

func nilPtrTo[T any]() *T {
	var v *T
	return v
}

func ptrTo[T any](v T) *T {
	return &v
}

func getArg[T any](v T) any {
	return v
}

var ErrInvalidCall = errors.New("invalid call")

func callFunc(callee any, args []any, autoZeroValues bool) (any, error) {
	rCallee := reflect.ValueOf(callee)
	funType := rCallee.Type()

	// variadic := funType.IsVariadic()
	totalParams := funType.NumIn()

	calleeArgs := make([]reflect.Value, len(args))
	for i := 0; i < totalParams; i++ {
		if i > len(args)-1 {
			return nil, fmt.Errorf("%w: function has %d parameters, but only %d arguments were provided", ErrInvalidCall, totalParams, len(args))
		}

		parameterType := funType.In(i)

		var arg reflect.Value

		if args[i] == nil {
			// fmt.Println("arg is untyped nil.")
			// turn an untyped nil into a typed one based on the parameter type
			//
			// f(*T) called with nil => create *T pointer
			// f(T) called with nil  => create *T pointer
			//
			// Turning this pointer potentially into the zero value is done later.
			ptrType := parameterType
			if parameterType.Kind() != reflect.Ptr {
				ptrType = reflect.PointerTo(parameterType)
			}

			arg = reflect.Zero(ptrType)
		} else {
			arg = reflect.ValueOf(args[i])
		}

		argType := arg.Type()

		// check whether the provided arg directly matches the required parameter type
		if !argType.AssignableTo(parameterType) {
			fixed := false

			// check if the func has a *T param, but T was given
			if !fixed && parameterType.Kind() == reflect.Ptr {
				paramPtrType := parameterType.Elem()

				if argType.AssignableTo(paramPtrType) {
					// create pointer to arg value
					newArg := reflect.New(arg.Type())
					newArg.Elem().Set(arg)

					arg = newArg
					fixed = true
				}
			}

			// if not, check the opposite: a T parameter and *T argument
			if !fixed && argType.Kind() == reflect.Ptr {
				if arg.IsNil() {
					if !autoZeroValues {
						return nil, fmt.Errorf("%w: cannot dereference nil pointer", ErrInvalidCall)
					}

					arg = reflect.Zero(parameterType)
					fixed = true
				} else {
					// deref pointer and check if the pointer pointed to something compatible
					if newArg := arg.Elem(); newArg.Type().AssignableTo(parameterType) {
						arg = newArg
						fixed = true
					}
				}
			}

			if !fixed {
				return nil, fmt.Errorf("%w: argument type %v is not assignable to %v", ErrInvalidCall, argType.String(), parameterType.String())
			}
		}

		calleeArgs[i] = arg
	}

	results := rCallee.Call(calleeArgs)

	// Forms can only be constructed with valid signatures,
	// no need to check that 2 values were returned.
	if err := results[1].Interface(); err != nil {
		return nil, err.(error)
	}

	return results[0].Interface(), nil
}

type xrstfTestcase struct {
	name        string
	f           any
	zeroing     bool
	args        []any
	expected    any
	expectedErr bool
	invalid     bool
}

func (tc *xrstfTestcase) Test(t *testing.T) {
	if tc.expectedErr && tc.invalid {
		panic("Invalid testcase: Cannot assert both invalid signature and error from calling the callee at the same time.")
	}

	result, err := callFunc(tc.f, tc.args, tc.zeroing)
	if err != nil {
		if errors.Is(err, ErrInvalidCall) {
			if !tc.invalid {
				t.Fatalf("Arguments should have matched, but did not: %v", err)
			}
		} else {
			if !tc.expectedErr {
				t.Fatalf("Did not expect an error from the callee, but got: %v", err)
			}
		}

		t.Logf("Test returned error (as expected): %v", err)
		return
	}

	if tc.invalid {
		t.Fatalf("Args should not have matched, but callee was still called and returned: %v (%T)", result, result)
	}

	if tc.expectedErr {
		t.Fatalf("Callee should have returned an error, but returned: %v (%T)", result, result)
	}

	if !cmp.Equal(result, tc.expected) {
		t.Fatalf("Expected %v (%T), but got %v (%T)", tc.expected, tc.expected, result, result)
	}
}

func TestCallFunc(t *testing.T) {
	testcases := []xrstfTestcase{
		{
			name:     "f()<--()",
			f:        func() (any, error) { return nil, nil },
			args:     []any{},
			expected: nil,
		},
		{
			name:     "f(int)<--(int)",
			f:        func(i int) (any, error) { return i, nil },
			args:     []any{getArg(1)},
			expected: 1,
		},
		{
			name:    "f(int)<--(string) [invalid]",
			f:       func(i int) (any, error) { return i, nil },
			args:    []any{getArg("foo")},
			invalid: true,
		},
		{
			name:     "f(int)<--(*int) (auto-deref)",
			f:        func(i int) (any, error) { return i, nil },
			args:     []any{getArg(ptrTo(1))},
			expected: 1,
		},
		{
			name:    "f(int)<--(*int<nil>) (auto-deref) (no auto-zeroing) [invalid]",
			f:       func(i int) (any, error) { return i, nil },
			args:    []any{getArg(nilPtrTo[int]())},
			invalid: true,
		},
		{
			name:     "f(int)<--(*int<nil>) (auto-deref) (auto-zeroing)",
			f:        func(i int) (any, error) { return i, nil },
			zeroing:  true,
			args:     []any{getArg(nilPtrTo[int]())},
			expected: 0,
		},
		{
			name:    "f(int)<--(*string<nil>) (auto-deref) (auto-zeroing)",
			f:       func(i int) (any, error) { return i, nil },
			zeroing: true,
			args:    []any{getArg(nilPtrTo[string]())},
			// This is different from Go, of course, where you cannot call f() with
			// a *string variable, because Go is statically typed. But in Rudi the
			// type checks happen at runtime, so it's fine to turn a nil *string
			// pointer into a *int and ultimately int value.
			expected: 0,
		},
		{
			name:    "f(int)<--(nil) (auto-deref) (no auto-zeroing) [invalid]",
			f:       func(i int) (any, error) { return i, nil },
			args:    []any{nil},
			invalid: true,
		},
		{
			name:     "f(int)<--(nil) (auto-deref) (auto-zeroing) [invalid]",
			f:        func(i int) (any, error) { return i, nil },
			zeroing:  true,
			args:     []any{nil},
			expected: 0,
		},
		{
			name:     "f(*int)<--(*int)",
			f:        func(i *int) (any, error) { return i, nil },
			args:     []any{getArg(ptrTo(1))},
			expected: ptrTo(1),
		},
		{
			name:    "f(*int)<--(*string) [invalid]",
			f:       func(i *int) (any, error) { return i, nil },
			args:    []any{getArg(ptrTo("foo"))},
			invalid: true,
		},
		{
			name:     "f(*int)<--(int) (auto-pointer)",
			f:        func(i *int) (any, error) { return i, nil },
			args:     []any{getArg(1)},
			expected: ptrTo(1),
		},
		{
			name:    "f(*int)<--(string) (auto-pointer) [invalid]",
			f:       func(i *int) (any, error) { return i, nil },
			args:    []any{getArg("foo")},
			invalid: true,
		},
		{
			name:     "f(*int)<--(plain int) (auto-pointer)",
			f:        func(i *int) (any, error) { return i, nil },
			args:     []any{1},
			expected: ptrTo(1),
		},
		{
			name:     "f(*int)<--(*int<nil>)",
			f:        func(i *int) (any, error) { return i, nil },
			args:     []any{nilPtrTo[int]()},
			expected: nilPtrTo[int](),
		},
		{
			name:     "f(*int)<--(nil)",
			f:        func(i *int) (any, error) { return i, nil },
			args:     []any{nil},
			expected: nilPtrTo[int](),
		},
		{
			name:     "f(**int)<--(**int)",
			f:        func(i **int) (any, error) { return i, nil },
			args:     []any{getArg(ptrTo(ptrTo(1)))},
			expected: ptrTo(ptrTo(1)),
		},
		{
			name:     "f(**int)<--(*int) (auto-pointer)",
			f:        func(i **int) (any, error) { return i, nil },
			args:     []any{getArg(ptrTo(1))},
			expected: ptrTo(ptrTo(1)),
		},
		{
			name:     "f(**int)<--(*int<nil>) (auto-pointer)",
			f:        func(i **int) (any, error) { return i, nil },
			args:     []any{getArg(nilPtrTo[int]())},
			expected: ptrTo(nilPtrTo[int]()),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, testcase.Test)
	}
}
