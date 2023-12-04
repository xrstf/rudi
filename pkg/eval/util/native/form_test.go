// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package native

import (
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
			name: "additional variadic parameter without args",
			fun: func(a string, b string, c ...int64) (any, error) {
				if a != "foo" {
					return false, fmt.Errorf("expected %q, got %q", "foo", a)
				}

				if b != "bar" {
					return false, fmt.Errorf("expected %q, got %q", "bar", b)
				}

				if len(c) != 0 {
					return false, fmt.Errorf("expected 0 ints, got %d", len(c))
				}

				return true, nil
			},
			args: []ast.Expression{
				ast.String("foo"),
				ast.String("bar"),
			},
		},
		{
			name: "empty variadic function",
			fun: func(c ...int64) (any, error) {
				if len(c) != 0 {
					return false, fmt.Errorf("expected 0 ints, got %d", len(c))
				}

				return true, nil
			},
			args: []ast.Expression{},
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
	}

	coalescer := coalescing.NewHumane()
	ctx := types.NewContext(types.Document{}, nil, nil, coalescer)

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

// func mySingleStringFunc(a string, b int64) (any, error) {
// 	fmt.Printf("a: %q\n", a)
// 	fmt.Printf("b: %d\n", b)
// 	return "i was called!", nil
// }

// func TestFoo(t *testing.T) {
// 	native, err := newNativeFunction(mySingleStringFunc)
// 	if err != nil {
// 		t.Fatalf("Failed to create native func: %v", err)
// 	}

// 	args := []ast.Expression{
// 		ast.Number{Value: 3.14},
// 		ast.Tuple{
// 			Expressions: []ast.Expression{
// 				ast.Identifier{Name: "+"},
// 				ast.Number{Value: 1},
// 				ast.Number{Value: 2},
// 			},
// 		},
// 	}

// 	ctx := types.NewContext(types.Document{}, nil, builtin.AllFunctions, coalescing.NewHumane())
// 	cachedArgs := convertArgs(args)

// 	matches, err := native.Match(ctx, cachedArgs)
// 	if err != nil {
// 		t.Fatalf("Match() failed: %v", err)
// 	}

// 	if !matches {
// 		fmt.Println("does not match")
// 		return
// 	}

// 	result, err := native.Call(ctx)
// 	if err != nil {
// 		t.Fatalf("Call() failed: %v", err)
// 	}

// 	fmt.Println(result)
// }
