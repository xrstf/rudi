// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

// Package lists_tests is a standalone package for the tests to prevent the regular
// lists package from having a dependency on the core/math/strings modules. Writing
// tests for higher-order functions like map however makes it kind of necessary to
// have some other helper functions like set! available.
package lists_test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/builtin/core"
	"go.xrstf.de/rudi/pkg/builtin/lists"
	"go.xrstf.de/rudi/pkg/builtin/math"
	"go.xrstf.de/rudi/pkg/builtin/strings"
	"go.xrstf.de/rudi/pkg/testutil"
)

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
			Expression: `(range [1 2 3] [a] (do (set! $foo $a) (+ $foo 3)))`,
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
		testcase.Functions = lists.Functions.DeepCopy().Add(core.Functions).Add(strings.Functions).Add(math.Functions)
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
			Expression: `(map [1 2 3] [val] (do (+ $val 3) "foo"))`,
			Expected:   []any{"foo", "foo", "foo"},
		},
		{
			// multiple expressions that use a common context
			Expression: `(map [1 2 3] [val] (do (set! $foo $val) (+ $foo 3)))`,
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
		testcase.Functions = lists.Functions.DeepCopy().Add(core.Functions).Add(strings.Functions).Add(math.Functions)
		t.Run(testcase.String(), testcase.Run)
	}
}
