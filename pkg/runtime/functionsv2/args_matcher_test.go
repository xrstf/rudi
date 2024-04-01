// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package functionsv2

import (
	"context"
	"testing"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/interpreter"
	"go.xrstf.de/rudi/pkg/runtime/types"

	"github.com/google/go-cmp/cmp"
)

type emptyInterface interface{}

type nonemptyInterface interface {
	Test()
}

type emptyStruct struct{}

type nonemptyStruct struct {
	Test string
}

type ExpressionAlias ast.Expression
type ExpressionCopy = ast.Expression
type ExpressionEmbed struct {
	ast.Expression
}

func TestNewArgsMatcherSignatures(t *testing.T) {
	testcases := []struct {
		name    string
		fun     any
		invalid bool
	}{
		// check invalid return values

		{
			name:    "no return values",
			fun:     func() {},
			invalid: true,
		},
		{
			name:    "too few return values",
			fun:     func() int { panic("") },
			invalid: true,
		},
		{
			name:    "too few return values",
			fun:     func() any { panic("") },
			invalid: true,
		},
		{
			name:    "too many return values",
			fun:     func() (any, error, any) { panic("") },
			invalid: true,
		},
		{
			name:    "invalid first type",
			fun:     func() (int, error) { panic("") },
			invalid: true,
		},
		{
			name:    "invalid first type",
			fun:     func() (emptyInterface, error) { panic("") },
			invalid: true,
		},
		{
			name:    "invalid first type",
			fun:     func() (nonemptyInterface, error) { panic("") },
			invalid: true,
		},
		{
			name:    "invalid first type",
			fun:     func() (emptyStruct, error) { panic("") },
			invalid: true,
		},
		{
			name:    "invalid first type",
			fun:     func() (nonemptyStruct, error) { panic("") },
			invalid: true,
		},
		{
			name:    "invalid first type",
			fun:     func() (*emptyStruct, error) { panic("") },
			invalid: true,
		},
		{
			name:    "invalid first type",
			fun:     func() (func(), error) { panic("") },
			invalid: true,
		},
		{
			name:    "invalid second type",
			fun:     func() (any, int) { panic("") },
			invalid: true,
		},
		{
			name:    "invalid second type",
			fun:     func() (any, emptyInterface) { panic("") },
			invalid: true,
		},
		{
			name:    "invalid second type",
			fun:     func() (any, nonemptyInterface) { panic("") },
			invalid: true,
		},
		{
			name:    "invalid second type",
			fun:     func() (any, emptyStruct) { panic("") },
			invalid: true,
		},
		{
			name:    "invalid second type",
			fun:     func() (any, nonemptyStruct) { panic("") },
			invalid: true,
		},
		{
			name:    "invalid second type",
			fun:     func() (*any, emptyStruct) { panic("") },
			invalid: true,
		},

		// check which parameter types are allowed

		{
			name: "no args",
			fun:  func() (any, error) { panic("") },
		},
		{
			name: "accept bool arg",
			fun:  func(bool) (any, error) { panic("") },
		},
		{
			name: "accept int64 arg",
			fun:  func(int64) (any, error) { panic("") },
		},
		{
			name: "accept float64 arg",
			fun:  func(float64) (any, error) { panic("") },
		},
		{
			name: "accept string arg",
			fun:  func(string) (any, error) { panic("") },
		},
		{
			name: "accept any arg",
			fun:  func(any) (any, error) { panic("") },
		},
		{
			name: "accept []any arg",
			fun:  func([]any) (any, error) { panic("") },
		},
		{
			name: "accept map[string]any arg",
			fun:  func(map[string]any) (any, error) { panic("") },
		},
		{
			name: "accept ast.Expression arg",
			fun:  func(ast.Expression) (any, error) { panic("") },
		},
		{
			name: "accept ExpressionAlias arg",
			fun:  func(ExpressionAlias) (any, error) { panic("") },
		},
		{
			name: "accept ExpressionCopy arg",
			fun:  func(ExpressionCopy) (any, error) { panic("") },
		},
		{
			name: "accept types.Context arg",
			fun:  func(types.Context) (any, error) { panic("") },
		},
		{
			name:    "reject custom emptyInterface arg",
			fun:     func(emptyInterface) (any, error) { panic("") },
			invalid: true,
		},
		{
			name:    "reject ExpressionEmbed arg",
			fun:     func(ExpressionEmbed) (any, error) { panic("") },
			invalid: true, // TODO: maybe support this?
		},

		// check which parameter combinations are allowed

		{
			name: "accept multiple plain arguments",
			fun:  func(string, any, int64) (any, error) { panic("") },
		},
		{
			name: "accept multiple plain and complex arguments",
			fun:  func(int64, []any, string, map[string]any) (any, error) { panic("") },
		},
		{
			name:    "reject complex vectors",
			fun:     func([]string, []ast.Expression, []map[string]any) (any, error) { panic("") },
			invalid: true,
		},
		{
			name:    "reject vector of basic maps",
			fun:     func([]map[string]any) (any, error) { panic("") },
			invalid: true,
		},
		{
			name:    "reject complex maps",
			fun:     func(map[string]string) (any, error) { panic("") },
			invalid: true,
		},
		{
			name:    "reject complex maps",
			fun:     func(map[string][]any) (any, error) { panic("") },
			invalid: true,
		},

		// variadic handling

		{
			name: "accept only variadic arg",
			fun:  func(...int64) (any, error) { panic("") },
		},
		{
			name: "accept variadic basic arg",
			fun:  func(string, []any, ...int64) (any, error) { panic("") },
		},
		{
			name: "accept variadic basic arg",
			fun:  func(string, ...any) (any, error) { panic("") },
		},
		{
			name: "accept variadic vector arg",
			fun:  func(string, ...[]any) (any, error) { panic("") },
		},
		{
			name: "accept variadic object arg",
			fun:  func(string, ...map[string]any) (any, error) { panic("") },
		},
		{
			name:    "reject variadic vector of maps arg",
			fun:     func(string, ...[]map[string]any) (any, error) { panic("") },
			invalid: true,
		},
		{
			name:    "reject non-consuming variadics",
			fun:     func(...types.Context) (any, error) { panic("") },
			invalid: true,
		},
		{
			name:    "reject non-consuming variadics",
			fun:     func(string, ...types.Context) (any, error) { panic("") },
			invalid: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := newArgsMatcher(tc.fun)
			if err != nil {
				if !tc.invalid {
					t.Fatalf("Failed to accept signature: %v", err)
				}

				return
			}

			if tc.invalid {
				t.Fatalf("Should not have accepted the invalid function signature, but did.")
			}
		})
	}
}

func TestNewArgsMatcherMinMaxArgs(t *testing.T) {
	testcases := []struct {
		name    string
		fun     any
		minArgs int
		maxArgs int
	}{
		{
			name:    "no parameters",
			fun:     func() (any, error) { panic("") },
			minArgs: 0,
			maxArgs: 0,
		},
		{
			name:    "simple parameters",
			fun:     func(string, int64) (any, error) { panic("") },
			minArgs: 2,
			maxArgs: 2,
		},
		{
			name:    "contexts do not count",
			fun:     func(string, int64, types.Context, string, types.Context) (any, error) { panic("") },
			minArgs: 3,
			maxArgs: 3,
		},
		{
			name:    "simple parameters with variadic",
			fun:     func(string, int64, ...string) (any, error) { panic("") },
			minArgs: 3, // remember, variadics in Rudi require at least 1 arg
			maxArgs: noLimit,
		},
		{
			name:    "pure variadic",
			fun:     func(...string) (any, error) { panic("") },
			minArgs: 1,
			maxArgs: noLimit,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			matcher, err := newArgsMatcher(tc.fun)
			if err != nil {
				t.Fatalf("Failed to accept signature: %v", err)
			}

			if matcher.minArgs != tc.minArgs {
				t.Errorf("expected minArgs=%d, got %d", tc.minArgs, matcher.minArgs)
			}

			if matcher.maxArgs != tc.maxArgs {
				t.Errorf("expected maxArgs=%d, got %d", tc.maxArgs, matcher.maxArgs)
			}
		})
	}
}

func TestArgsMatcher(t *testing.T) {
	testcases := []struct {
		name     string
		fun      any
		args     []ast.Expression
		match    bool
		expected []any
	}{
		{
			name:     "no parameters",
			fun:      func() (any, error) { return nil, nil },
			args:     []ast.Expression{},
			match:    true,
			expected: []any{},
		},
		{
			name: "matching parameters, using coalescing",
			fun:  func(bool, string, int64) (any, error) { return nil, nil },
			args: []ast.Expression{
				ast.String("true"),
				ast.Number{Value: 42},
				ast.Bool(true),
			},
			match:    true,
			expected: []any{true, "42", int64(1)},
		},
		{
			name: "vectors work",
			fun:  func([]any) (any, error) { return nil, nil },
			args: []ast.Expression{
				ast.VectorNode{
					Expressions: []ast.Expression{
						ast.Bool(true),
						ast.String("hello"),
					},
				},
			},
			match:    true,
			expected: []any{[]any{true, "hello"}},
		},
		{
			name: "objects work",
			fun:  func(map[string]any) (any, error) { return nil, nil },
			args: []ast.Expression{
				ast.ObjectNode{
					Data: []ast.KeyValuePair{
						{
							Key:   ast.String("foo"),
							Value: ast.String("bar"),
						},
					},
				},
			},
			match:    true,
			expected: []any{map[string]any{"foo": "bar"}},
		},
		{
			name: "raw data works",
			fun:  func(any) (any, error) { return nil, nil },
			args: []ast.Expression{
				ast.Shim{Value: []string{"foo", "bar"}},
			},
			match:    true,
			expected: []any{[]string{"foo", "bar"}},
		},
		{
			name: "unevaluated expressions work",
			fun:  func(ast.Expression) (any, error) { return nil, nil },
			args: []ast.Expression{
				ast.Shim{Value: 1},
			},
			match:    true,
			expected: []any{ast.Shim{Value: 1}},
		},
		{
			name:  "variadic parameters require at least 1 item",
			fun:   func(...string) (any, error) { return nil, nil },
			args:  []ast.Expression{},
			match: false,
		},
		{
			name:     "basic variadic parameter",
			fun:      func(...string) (any, error) { return nil, nil },
			args:     []ast.Expression{ast.String("foo"), ast.String("bar")},
			match:    true,
			expected: []any{"foo", "bar"},
		},
		{
			name:     "mixed normal and variadic parameters",
			fun:      func(string, ...string) (any, error) { return nil, nil },
			args:     []ast.Expression{ast.String("foo"), ast.String("bar")},
			match:    true,
			expected: []any{"foo", "bar"},
		},
		{
			name:  "variadic cannot be empty",
			fun:   func(string, ...string) (any, error) { return nil, nil },
			args:  []ast.Expression{ast.String("foo")},
			match: false,
		},
		{
			name: "variadic slices",
			fun:  func(string, ...[]any) (any, error) { return nil, nil },
			args: []ast.Expression{
				ast.String("foo"),
				ast.VectorNode{
					Expressions: []ast.Expression{
						ast.String("hello"),
					},
				},
				ast.VectorNode{
					Expressions: []ast.Expression{
						ast.Bool(true),
					},
				},
			},
			match:    true,
			expected: []any{"foo", []any{"hello"}, []any{true}},
		},
		{
			name:  "too many args given",
			fun:   func() (any, error) { return nil, nil },
			args:  []ast.Expression{ast.Bool(true)},
			match: false,
		},
		{
			name:  "too many args given",
			fun:   func(bool) (any, error) { return nil, nil },
			args:  []ast.Expression{ast.Bool(true), ast.Null{}},
			match: false,
		},
		{
			name:  "too few args given",
			fun:   func(bool) (any, error) { return nil, nil },
			args:  []ast.Expression{},
			match: false,
		},
		{
			name:  "too few args given",
			fun:   func(bool, bool) (any, error) { return nil, nil },
			args:  []ast.Expression{ast.Bool(true)},
			match: false,
		},
		{
			name:     "can coalesce to desired type",
			fun:      func(int64) (any, error) { return nil, nil },
			args:     []ast.Expression{ast.String("15")},
			match:    true,
			expected: []any{int64(15)},
		},
		{
			name:  "cannot coalesce to desired type",
			fun:   func(int64) (any, error) { return nil, nil },
			args:  []ast.Expression{ast.String("foo")},
			match: false,
		},
	}

	coalescer := coalescing.NewHumane()
	ctx, err := types.NewContext(interpreter.New(), context.Background(), types.Document{}, nil, nil, coalescer)
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			matcher, err := newArgsMatcher(tc.fun)
			if err != nil {
				t.Fatalf("Failed to create matcher: %v", err)
			}

			result, matched, err := matcher.Match(ctx, convertArgs(tc.args))
			if err != nil {
				t.Fatalf("Failed to run matcher: %v", err)
			}

			if matched != tc.match {
				t.Fatalf("Expected match=%v", tc.match)
			}

			if !cmp.Equal(result, tc.expected) {
				t.Fatalf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}
