// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package functions

import (
	"errors"
	"fmt"
	"testing"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func TestFormCalling(t *testing.T) {
	testcases := []struct {
		name string
		fun  any
		args []ast.Expression
	}{
		{
			name: "no parameters",
			fun: func() (any, error) {
				return true, nil
			},
			args: []ast.Expression{},
		},
		{
			name: "nil values",
			fun: func(v any) (any, error) {
				if v != nil {
					return false, fmt.Errorf("expectednil, got %v", v)
				}

				return true, nil
			},
			args: []ast.Expression{
				ast.Null{},
			},
		},
		{
			name: "basic arguments",
			fun: func(a string, b string) (any, error) {
				if a != "foo" {
					return false, fmt.Errorf("expected %q, got %q", "foo", a)
				}

				if b != "bar" {
					return false, fmt.Errorf("expected %q, got %q", "bar", b)
				}

				return true, nil
			},
			args: []ast.Expression{
				ast.String("foo"),
				ast.String("bar"),
			},
		},
		{
			name: "additional variadic parameter with args",
			fun: func(a string, b string, c ...int64) (any, error) {
				if a != "foo" {
					return false, fmt.Errorf("expected %q, got %q", "foo", a)
				}

				if b != "bar" {
					return false, fmt.Errorf("expected %q, got %q", "bar", b)
				}

				if len(c) != 3 {
					return false, fmt.Errorf("expected 3 ints, got %d", len(c))
				}

				return true, nil
			},
			args: []ast.Expression{
				ast.String("foo"),
				ast.String("bar"),
				ast.Number{Value: 1},
				ast.Number{Value: 2},
				ast.Number{Value: 3},
			},
		},
		{
			name: "only variadic arguments",
			fun: func(c ...int64) (any, error) {
				if len(c) != 3 {
					return false, fmt.Errorf("expected 3 ints, got %d", len(c))
				}

				return true, nil
			},
			args: []ast.Expression{
				ast.Number{Value: 1},
				ast.Number{Value: 2},
				ast.Number{Value: 3},
			},
		},
		{
			name: "vector support",
			fun: func(c []any) (any, error) {
				if len(c) != 3 {
					return false, fmt.Errorf("expected 3 ints, got %d", len(c))
				}

				return true, nil
			},
			args: []ast.Expression{
				ast.VectorNode{
					Expressions: []ast.Expression{
						ast.Number{Value: 1},
						ast.Number{Value: 2},
						ast.Number{Value: 3},
					},
				},
			},
		},
		{
			name: "variadic vector support",
			fun: func(c ...[]any) (any, error) {
				if len(c) != 2 {
					return false, fmt.Errorf("expected 2 vectors, got %d", len(c))
				}

				return true, nil
			},
			args: []ast.Expression{
				ast.VectorNode{
					Expressions: []ast.Expression{
						ast.Number{Value: 1},
						ast.Number{Value: 2},
						ast.Number{Value: 3},
					},
				},
				ast.VectorNode{
					Expressions: []ast.Expression{
						ast.Number{Value: 1},
						ast.Number{Value: 2},
						ast.Number{Value: 3},
					},
				},
			},
		},
		{
			name: "context support",
			fun: func(ctx types.Context) (any, error) {
				if ctx.Coalesce() == nil {
					return false, errors.New("got empty context")
				}

				return true, nil
			},
			args: []ast.Expression{},
		},
		{
			name: "context support with other args",
			fun: func(ctx types.Context, c string) (any, error) {
				if c != "foo" {
					return false, fmt.Errorf("expected %q, got %q", "foo", c)
				}

				return true, nil
			},
			args: []ast.Expression{
				ast.String("foo"),
			},
		},
	}

	coalescer := coalescing.NewHumane()
	ctx := types.NewContext(nil, types.Document{}, nil, nil, coalescer)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := newForm(tc.fun)
			if err != nil {
				t.Fatalf("Failed to create form: %v", err)
			}

			matched, err := f.Match(ctx, convertArgs(tc.args))
			if err != nil {
				t.Fatalf("Failed to run matcher: %v", err)
			}

			if !matched {
				t.Fatalf("Form did not match expression arguments.")
			}

			result, err := f.Call(ctx)
			if err != nil {
				t.Fatalf("Failed to call function: %v", err)
			}

			if result != true {
				t.Fatalf("Expected true, got %v", result)
			}
		})
	}
}
