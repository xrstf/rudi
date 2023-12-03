// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"

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
			Expression: `(len [] [])`,
			Invalid:    true,
		},
		{
			// strict coalescing allows null to turn into [] or "", both have len=0
			Expression: `(len null)`,
			Expected:   0,
		},
		{
			Expression: `(len "")`,
			Expected:   0,
		},
		{
			Expression: `(len " foo ")`,
			Expected:   5,
		},
		{
			Expression: `(len [])`,
			Expected:   0,
		},
		{
			Expression: `(len [1 2 3])`,
			Expected:   3,
		},
		{
			Expression: `(len {})`,
			Expected:   0,
		},
		{
			Expression: `(len {foo "bar" hello "world"})`,
			Expected:   2,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = AllFunctions
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
			Expression: `(append {} 1)`,
			Invalid:    true,
		},
		{
			// strict coalescing allows null to turn into [] (null could also be "",
			// which would result in "1", but append/prepend prefer the first arg
			// to be a vector)
			Expression: `(append null 1)`,
			Expected:   []any{int64(1)},
		},
		{
			Expression: `(append [] 1)`,
			Expected:   []any{int64(1)},
		},
		{
			Expression: `(append [1 2] 3 "foo")`,
			Expected:   []any{int64(1), int64(2), int64(3), "foo"},
		},
		{
			Expression: `(append [] [])`,
			Expected:   []any{[]any{}},
		},
		{
			Expression: `(append [] "foo")`,
			Expected:   []any{"foo"},
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
			Expected:   "foobartest",
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = AllFunctions
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
			Expression: `(prepend {} 1)`,
			Invalid:    true,
		},
		{
			// strict coalescing allows null to turn into [] (null could also be "",
			// which would result in "1", but append/prepend prefer the first arg
			// to be a vector)
			Expression: `(prepend null 1)`,
			Expected:   []any{int64(1)},
		},
		{
			Expression: `(prepend [] 1)`,
			Expected:   []any{int64(1)},
		},
		{
			Expression: `(prepend [1] 2)`,
			Expected:   []any{int64(2), int64(1)},
		},
		{
			Expression: `(prepend [1 2] 3 "foo")`,
			Expected:   []any{int64(3), "foo", int64(1), int64(2)},
		},
		{
			Expression: `(prepend [] [])`,
			Expected:   []any{[]any{}},
		},
		{
			Expression: `(prepend [] "foo")`,
			Expected:   []any{"foo"},
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
			Expected:   "bartestfoo",
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = AllFunctions
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
			Expression: `(reverse {})`,
			Invalid:    true,
		},
		{
			// strict coalescing allows null to turn into ""
			Expression: `(reverse null)`,
			Expected:   "",
		},
		{
			Expression: `(reverse "")`,
			Expected:   "",
		},
		{
			Expression: `(reverse (concat "" "f" "oo"))`,
			Expected:   "oof",
		},
		{
			Expression: `(reverse "abcd")`,
			Expected:   "dcba",
		},
		{
			Expression: `(reverse (reverse "abcd"))`,
			Expected:   "abcd",
		},
		{
			Expression: `(reverse [])`,
			Expected:   []any{},
		},
		{
			Expression: `(reverse [1])`,
			Expected:   []any{int64(1)},
		},
		{
			Expression: `(reverse [1 2 3])`,
			Expected:   []any{int64(3), int64(2), int64(1)},
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = AllFunctions
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
			Expected:   int64(3),
		},
		{
			// multiple expressions that use a common context
			Expression: `(range [1 2 3] [a] (set! $foo $a) (+ $foo 3))`,
			Expected:   int64(6),
		},
		{
			// count iterations
			Expression: `(range [1 2 3] [loop-var] (set! $counter (+ (default (try $counter) 0) 1)))`,
			Expected:   int64(3),
		},
		{
			// value is bound to desired variable
			Expression: `(range [1 2 3] [a] $a)`,
			Expected:   int64(3),
		},
		{
			// support loop index variable
			Expression: `(range [1 2 3] [idx var] $idx)`,
			Expected:   2,
		},
		{
			// support loop index variable
			Expression: `(range [1 2 3] [idx var] $var)`,
			Expected:   int64(3),
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
			Expected:   nil,
		},
		{
			Expression: `(range {foo "bar"} [key value] $key)`,
			Expected:   "foo",
		},
		{
			Expression: `(range {foo "bar"} [key value] $value)`,
			Expected:   "bar",
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = AllFunctions
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
			Expected:   []any{"FOO", "BAR"},
		},
		{
			Expression: `(map {foo "bar"} to-upper)`,
			Expected:   map[string]any{"foo": "BAR"},
		},
		{
			// type safety still applies
			Expression: `(map [1] to-upper)`,
			Invalid:    true,
		},
		{
			// eval expression with variable
			Expression: `(map [1 2 3] [val] (+ $val 3))`,
			Expected:   []any{int64(4), int64(5), int64(6)},
		},
		{
			// eval with loop index
			Expression: `(map ["foo" "bar"] [idx _] $idx)`,
			Expected:   []any{0, 1},
		},
		{
			// last expression controls the result
			Expression: `(map [1 2 3] [val] (+ $val 3) "foo")`,
			Expected:   []any{"foo", "foo", "foo"},
		},
		{
			// multiple expressions that use a common context
			Expression: `(map [1 2 3] [val] (set! $foo $val) (+ $foo 3))`,
			Expected:   []any{int64(4), int64(5), int64(6)},
		},
		{
			// context is even shared across elements
			Expression: `(map ["foo" "bar"] [_] (set! $counter (+ (try $counter 0) 1)))`,
			Expected:   []any{int64(1), int64(2)},
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
			Expected:   []any{"bar", "bar", "bar"},
		},
		{
			Expression: `(set! $foo [1 2 3]) (map $foo [_ __] "bar") $foo`,
			Expected:   []any{int64(1), int64(2), int64(3)},
		},
		{
			Expression: `(set! $foo {foo "bar"}) (map $foo [_ __] "new-value") $foo`,
			Expected:   map[string]any{"foo": "bar"},
		},
		{
			Expression: `(set! $foo ["foo" "bar"]) (map $foo to-upper)`,
			Expected:   []any{"FOO", "BAR"},
		},
		{
			Expression: `(set! $foo ["foo" "bar"]) (map $foo to-upper) $foo`,
			Expected:   []any{"foo", "bar"},
		},
		{
			Expression: `(set! $foo {foo "bar"}) (map $foo to-upper) $foo`,
			Expected:   map[string]any{"foo": "bar"},
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = AllFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}
