// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/testutil"
)

func TestLenFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(len)`,
			Invalid:    true,
		},
		{
			Expression: `(len true)`,
			Invalid:    true,
		},
		{
			Expression: `(len 1)`,
			Invalid:    true,
		},
		{
			Expression: `(len null)`,
			Invalid:    true,
		},
		{
			Expression: `(len [] [])`,
			Invalid:    true,
		},
		{
			Expression: `(len "")`,
			Expected:   ast.Number{Value: int64(0)},
		},
		{
			Expression: `(len " foo ")`,
			Expected:   ast.Number{Value: int64(5)},
		},
		{
			Expression: `(len [])`,
			Expected:   ast.Number{Value: int64(0)},
		},
		{
			Expression: `(len [1 2 3])`,
			Expected:   ast.Number{Value: int64(3)},
		},
		{
			Expression: `(len {})`,
			Expected:   ast.Number{Value: int64(0)},
		},
		{
			Expression: `(len {foo "bar" hello "world"})`,
			Expected:   ast.Number{Value: int64(2)},
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestAppendFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(append)`,
			Invalid:    true,
		},
		{
			Expression: `(append [])`,
			Invalid:    true,
		},
		{
			Expression: `(append true 1)`,
			Invalid:    true,
		},
		{
			Expression: `(append 1 1)`,
			Invalid:    true,
		},
		{
			Expression: `(append null 1)`,
			Invalid:    true,
		},
		{
			Expression: `(append {} 1)`,
			Invalid:    true,
		},
		{
			Expression: `(append {} 1)`,
			Invalid:    true,
		},
		{
			Expression: `(append [] 1)`,
			Expected:   ast.Vector{Data: []any{ast.Number{Value: 1}}},
		},
		{
			Expression: `(append [1 2] 3 "foo")`,
			Expected:   ast.Vector{Data: []any{ast.Number{Value: 1}, ast.Number{Value: 2}, ast.Number{Value: 3}, ast.String("foo")}},
		},
		{
			Expression: `(append [] [])`,
			Expected:   ast.Vector{Data: []any{ast.Vector{Data: []any{}}}},
		},
		{
			Expression: `(append [] "foo")`,
			Expected:   ast.Vector{Data: []any{ast.String("foo")}},
		},
		{
			Expression: `(append "foo" [])`,
			Invalid:    true,
		},
		{
			Expression: `(append "foo" "bar" [])`,
			Invalid:    true,
		},
		{
			Expression: `(append "foo" "bar" "test")`,
			Expected:   ast.String("foobartest"),
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestPrependFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(prepend)`,
			Invalid:    true,
		},
		{
			Expression: `(prepend [])`,
			Invalid:    true,
		},
		{
			Expression: `(prepend true 1)`,
			Invalid:    true,
		},
		{
			Expression: `(prepend 1 1)`,
			Invalid:    true,
		},
		{
			Expression: `(prepend null 1)`,
			Invalid:    true,
		},
		{
			Expression: `(prepend {} 1)`,
			Invalid:    true,
		},
		{
			Expression: `(prepend {} 1)`,
			Invalid:    true,
		},
		{
			Expression: `(prepend [] 1)`,
			Expected:   ast.Vector{Data: []any{ast.Number{Value: 1}}},
		},
		{
			Expression: `(prepend [1] 2)`,
			Expected:   ast.Vector{Data: []any{ast.Number{Value: 2}, ast.Number{Value: 1}}},
		},
		{
			Expression: `(prepend [1 2] 3 "foo")`,
			Expected:   ast.Vector{Data: []any{ast.Number{Value: 3}, ast.String("foo"), ast.Number{Value: 1}, ast.Number{Value: 2}}},
		},
		{
			Expression: `(prepend [] [])`,
			Expected:   ast.Vector{Data: []any{ast.Vector{Data: []any{}}}},
		},
		{
			Expression: `(prepend [] "foo")`,
			Expected:   ast.Vector{Data: []any{ast.String("foo")}},
		},
		{
			Expression: `(prepend "foo" [])`,
			Invalid:    true,
		},
		{
			Expression: `(prepend "foo" "bar" [])`,
			Invalid:    true,
		},
		{
			Expression: `(prepend "foo" "bar" "test")`,
			Expected:   ast.String("bartestfoo"),
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestReverseFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(reverse)`,
			Invalid:    true,
		},
		{
			Expression: `(reverse "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(reverse 1)`,
			Invalid:    true,
		},
		{
			Expression: `(reverse true)`,
			Invalid:    true,
		},
		{
			Expression: `(reverse null)`,
			Invalid:    true,
		},
		{
			Expression: `(reverse {})`,
			Invalid:    true,
		},
		{
			Expression: `(reverse "")`,
			Expected:   ast.String(""),
		},
		{
			Expression: `(reverse (concat "" "f" "oo"))`,
			Expected:   ast.String("oof"),
		},
		{
			Expression: `(reverse "abcd")`,
			Expected:   ast.String("dcba"),
		},
		{
			Expression: `(reverse (reverse "abcd"))`,
			Expected:   ast.String("abcd"),
		},
		{
			Expression: `(reverse [])`,
			Expected:   ast.Vector{Data: []any{}},
		},
		{
			Expression: `(reverse [1])`,
			Expected:   ast.Vector{Data: []any{ast.Number{Value: 1}}},
		},
		{
			Expression: `(reverse [1 2 3])`,
			Expected:   ast.Vector{Data: []any{ast.Number{Value: 3}, ast.Number{Value: 2}, ast.Number{Value: 1}}},
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestRangeFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			// missing everything
			Expression: `(range)`,
			Invalid:    true,
		},
		{
			// missing naming vector
			Expression: `(range [1 2 3])`,
			Invalid:    true,
		},
		{
			// missing naming vector
			Expression: `(range [1 2 3] (+ 1 2))`,
			Invalid:    true,
		},
		{
			// naming vector must be 1 or 2 elements long
			Expression: `(range [1 2 3] [] (+ 1 2))`,
			Invalid:    true,
		},
		{
			// naming vector must be 1 or 2 elements long
			Expression: `(range [1 2 3] [a b c] (+ 1 2))`,
			Invalid:    true,
		},
		{
			// do not allow numbers in the naming vector
			Expression: `(range [1 2 3] [1 2] (+ 1 2))`,
			Invalid:    true,
		},
		{
			// do not allow strings in naming vector
			Expression: `(range [1 2 3] ["foo" "bar"] (+ 1 2))`,
			Invalid:    true,
		},
		{
			// cannot range over non-vectors/objects
			Expression: `(range "invalid" [a] (+ 1 2))`,
			Invalid:    true,
		},
		{
			// cannot range over non-vectors/objects
			Expression: `(range 5 [a] (+ 1 2))`,
			Invalid:    true,
		},
		{
			// single simple expression
			Expression: `(range [1 2 3] [a] (+ 1 2))`,
			Expected:   ast.Number{Value: int64(3)},
		},
		{
			// multiple expressions that use a common context
			Expression: `(range [1 2 3] [a] (set! $foo $a) (+ $foo 3))`,
			Expected:   ast.Number{Value: int64(6)},
		},
		{
			// count iterations
			Expression: `(range [1 2 3] [loop-var] (set! $counter (+ (default (try $counter) 0) 1)))`,
			Expected:   ast.Number{Value: int64(3)},
		},
		{
			// value is bound to desired variable
			Expression: `(range [1 2 3] [a] $a)`,
			Expected:   ast.Number{Value: int64(3)},
		},
		{
			// support loop index variable
			Expression: `(range [1 2 3] [idx var] $idx)`,
			Expected:   ast.Number{Value: int64(2)},
		},
		{
			// support loop index variable
			Expression: `(range [1 2 3] [idx var] $var)`,
			Expected:   ast.Number{Value: int64(3)},
		},
		{
			// variables do not leak outside the range
			Expression: `(range [1 2 3] [idx var] $idx) (+ $var 0)`,
			Invalid:    true,
		},
		{
			// variables do not leak outside the range
			Expression: `(range [1 2 3] [idx var] $idx) (+ $idx 0)`,
			Invalid:    true,
		},
		{
			// support ranging over objects
			Expression: `(range {} [key value] $key)`,
			Expected:   ast.Null{},
		},
		{
			Expression: `(range {foo "bar"} [key value] $key)`,
			Expected:   ast.String("foo"),
		},
		{
			Expression: `(range {foo "bar"} [key value] $value)`,
			Expected:   ast.String("bar"),
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestMapFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			// missing everything
			Expression: `(map)`,
			Invalid:    true,
		},
		{
			// missing function identifier
			Expression: `(map [1 2 3])`,
			Invalid:    true,
		},
		{
			// missing naming vector
			Expression: `(map [1 2 3] (+ 1 2))`,
			Invalid:    true,
		},
		{
			// naming vector must be 1 or 2 elements long
			Expression: `(map [1 2 3] [] (+ 1 2))`,
			Invalid:    true,
		},
		{
			// naming vector must be 1 or 2 elements long
			Expression: `(map [1 2 3] [a b c] (+ 1 2))`,
			Invalid:    true,
		},
		{
			// do not allow numbers in the naming vector
			Expression: `(map [1 2 3] [1 2] (+ 1 2))`,
			Invalid:    true,
		},
		{
			// do not allow strings in naming vector
			Expression: `(map [1 2 3] ["foo" "bar"] (+ 1 2))`,
			Invalid:    true,
		},
		{
			// cannot map non-vectors/objects
			Expression: `(map "invalid" [a] (+ 1 2))`,
			Invalid:    true,
		},
		{
			// cannot map non-vectors/objects
			Expression: `(map 5 [a] (+ 1 2))`,
			Invalid:    true,
		},
		{
			// single simple expression
			Expression: `(map ["foo" "bar"] to-upper)`,
			Expected:   ast.Vector{Data: []any{ast.String("FOO"), ast.String("BAR")}},
		},
		{
			Expression: `(map {foo "bar"} to-upper)`,
			Expected:   ast.Object{Data: map[string]any{"foo": ast.String("BAR")}},
		},
		{
			// type safety still applies
			Expression: `(map [1] to-upper)`,
			Invalid:    true,
		},
		{
			// eval expression with variable
			Expression: `(map [1 2 3] [val] (+ $val 3))`,
			Expected:   ast.Vector{Data: []any{ast.Number{Value: 4}, ast.Number{Value: 5}, ast.Number{Value: 6}}},
		},
		{
			// eval with loop index
			Expression: `(map ["foo" "bar"] [idx _] $idx)`,
			Expected:   ast.Vector{Data: []any{ast.Number{Value: 0}, ast.Number{Value: 1}}},
		},
		{
			// last expression controls the result
			Expression: `(map [1 2 3] [val] (+ $val 3) "foo")`,
			Expected:   ast.Vector{Data: []any{ast.String("foo"), ast.String("foo"), ast.String("foo")}},
		},
		{
			// multiple expressions that use a common context
			Expression: `(map [1 2 3] [val] (set! $foo $val) (+ $foo 3))`,
			Expected:   ast.Vector{Data: []any{ast.Number{Value: 4}, ast.Number{Value: 5}, ast.Number{Value: 6}}},
		},
		{
			// context is even shared across elements
			Expression: `(map ["foo" "bar"] [_] (set! $counter (+ (try $counter 0) 1)))`,
			Expected:   ast.Vector{Data: []any{ast.Number{Value: 1}, ast.Number{Value: 2}}},
		},
		{
			// variables do not leak outside the range
			Expression: `(map [1 2 3] [idx var] $idx) (+ $var 0)`,
			Invalid:    true,
		},
		{
			// variables do not leak outside the range
			Expression: `(map [1 2 3] [idx var] $idx) (+ $idx 0)`,
			Invalid:    true,
		},
		// do not modify the source
		{
			Expression: `(set! $foo [1 2 3]) (map $foo [_ __] "bar")`,
			Expected:   ast.Vector{Data: []any{ast.String("bar"), ast.String("bar"), ast.String("bar")}},
		},
		{
			Expression: `(set! $foo [1 2 3]) (map $foo [_ __] "bar") $foo`,
			Expected:   ast.Vector{Data: []any{ast.Number{Value: 1}, ast.Number{Value: 2}, ast.Number{Value: 3}}},
		},
		{
			Expression: `(set! $foo {foo "bar"}) (map $foo [_ __] "new-value") $foo`,
			Expected:   ast.Object{Data: map[string]any{"foo": ast.String("bar")}},
		},
		{
			Expression: `(set! $foo ["foo" "bar"]) (map $foo to-upper)`,
			Expected:   ast.Vector{Data: []any{ast.String("FOO"), ast.String("BAR")}},
		},
		{
			Expression: `(set! $foo ["foo" "bar"]) (map $foo to-upper) $foo`,
			Expected:   ast.Vector{Data: []any{ast.String("foo"), ast.String("bar")}},
		},
		{
			Expression: `(set! $foo {foo "bar"}) (map $foo to-upper) $foo`,
			Expected:   ast.Object{Data: map[string]any{"foo": ast.String("bar")}},
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}
