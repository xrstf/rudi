// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package functions

import (
	"context"
	"testing"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/interpreter"
	"go.xrstf.de/rudi/pkg/runtime/types"

	"github.com/google/go-cmp/cmp"
)

type argsConsumerTestcase struct {
	name      string
	args      []ast.Expression
	expected  []any
	remaining int
	invalid   bool
}

func (tc *argsConsumerTestcase) Test(t *testing.T, ctx types.Context, consumer argsConsumer) {
	consumed, remaining, err := consumer(ctx, convertArgs(tc.args))
	if err != nil {
		if !tc.invalid {
			t.Fatalf("Failed to consume: %v", err)
		}

		return
	}

	if tc.invalid {
		t.Fatalf("Should not have consumed the input, but did: returned %v, remaining %d", consumed, len(remaining))
	}

	if !cmp.Equal(consumed, tc.expected) {
		t.Fatalf("Expected %v, but got %v", tc.expected, consumed)
	}

	if len(remaining) != tc.remaining {
		t.Fatalf("Expected %d remaining args, but %d remain.", tc.remaining, len(remaining))
	}
}

func TestBoolConsumer(t *testing.T) {
	testcases := []argsConsumerTestcase{
		{
			name: "simple bool",
			args: []ast.Expression{
				ast.Bool(true),
			},
			expected:  []any{true},
			remaining: 0,
		},
		{
			name: "keep rest",
			args: []ast.Expression{
				ast.Bool(true),
				ast.String("foo"),
				ast.Null{},
			},
			expected:  []any{true},
			remaining: 2,
		},
		{
			name: "apply coalescing",
			args: []ast.Expression{
				ast.String(""),
			},
			expected:  []any{false},
			remaining: 0,
		},
		{
			name: "coalescing can fail, then input should not be consumed",
			args: []ast.Expression{
				ast.Shim{Value: func() {}},
			},
			expected:  nil,
			remaining: 1,
		},
	}

	coalescer := coalescing.NewHumane()
	ctx, err := types.NewContext(interpreter.New(), context.Background(), types.Document{}, nil, nil, coalescer)
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.Test(t, ctx, boolConsumer)
		})
	}
}

func TestStringConsumer(t *testing.T) {
	testcases := []argsConsumerTestcase{
		{
			name: "simple string",
			args: []ast.Expression{
				ast.String("foo"),
			},
			expected:  []any{"foo"},
			remaining: 0,
		},
		{
			name: "keep rest",
			args: []ast.Expression{
				ast.String("foo"),
				ast.Bool(true),
				ast.Null{},
			},
			expected:  []any{"foo"},
			remaining: 2,
		},
		{
			name: "apply coalescing",
			args: []ast.Expression{
				ast.Bool(true),
			},
			expected:  []any{"true"},
			remaining: 0,
		},
		{
			name: "coalescing can fail, then input should not be consumed",
			args: []ast.Expression{
				ast.Shim{Value: func() {}},
			},
			expected:  nil,
			remaining: 1,
		},
	}

	coalescer := coalescing.NewHumane()
	ctx, err := types.NewContext(interpreter.New(), context.Background(), types.Document{}, nil, nil, coalescer)
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.Test(t, ctx, stringConsumer)
		})
	}
}

func TestAnyConsumer(t *testing.T) {
	testcases := []argsConsumerTestcase{
		{
			name: "simple expression",
			args: []ast.Expression{
				ast.String("foo"),
			},
			expected:  []any{"foo"},
			remaining: 0,
		},
		{
			name: "keep rest",
			args: []ast.Expression{
				ast.String("foo"),
				ast.Bool(true),
				ast.Null{},
			},
			expected:  []any{"foo"},
			remaining: 2,
		},
		{
			name: "allows nulls",
			args: []ast.Expression{
				ast.Null{},
			},
			expected:  []any{nil},
			remaining: 0,
		},
		{
			name: "no coalescing",
			args: []ast.Expression{
				// no coalescer would ever handle anything but []any;
				// []string was chosen here because it is comparable by cmp
				ast.Shim{Value: []string{"foo", "bar"}},
			},
			expected:  []any{[]string{"foo", "bar"}},
			remaining: 0,
		},
	}

	coalescer := coalescing.NewHumane()
	ctx, err := types.NewContext(interpreter.New(), context.Background(), types.Document{}, nil, nil, coalescer)
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.Test(t, ctx, anyConsumer)
		})
	}
}

func TestExpressionConsumer(t *testing.T) {
	testcases := []argsConsumerTestcase{
		{
			name: "simple expression",
			args: []ast.Expression{
				ast.String("foo"),
			},
			expected:  []any{ast.String("foo")},
			remaining: 0,
		},
		{
			name: "keep rest",
			args: []ast.Expression{
				ast.String("foo"),
				ast.Bool(true),
				ast.Null{},
			},
			expected:  []any{ast.String("foo")},
			remaining: 2,
		},
	}

	coalescer := coalescing.NewHumane()
	ctx, err := types.NewContext(interpreter.New(), context.Background(), types.Document{}, nil, nil, coalescer)
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.Test(t, ctx, expressionConsumer)
		})
	}
}

func TestVariadicConsumer(t *testing.T) {
	testcases := []argsConsumerTestcase{
		{
			name:      "variadic request at least 1 item",
			args:      []ast.Expression{},
			expected:  nil,
			remaining: 0,
		},
		{
			name: "consume all strings",
			args: []ast.Expression{
				ast.String("foo"),
				ast.Bool(true),
				ast.Null{},
			},
			expected:  []any{"foo", "true", ""},
			remaining: 0,
		},
	}

	coalescer := coalescing.NewHumane()
	ctx, err := types.NewContext(interpreter.New(), context.Background(), types.Document{}, nil, nil, coalescer)
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.Test(t, ctx, toVariadicConsumer(stringConsumer))
		})
	}
}
