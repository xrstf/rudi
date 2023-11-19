// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type listsTestcase struct {
	expr     string
	expected any
	invalid  bool
}

func (tc *listsTestcase) Test(t *testing.T) {
	t.Helper()

	result, err := runExpression(t, tc.expr, nil, nil)
	if err != nil {
		if !tc.invalid {
			t.Fatalf("Failed to run %s: %v", tc.expr, err)
		}

		return
	}

	if tc.invalid {
		t.Fatalf("Should not have been able to run %s, but got: %v", tc.expr, result)
	}

	if !cmp.Equal(result, tc.expected) {
		t.Fatalf("Did not receive expected output:\n%s", cmp.Diff(tc.expected, result))
	}
}

func TestLenFunction(t *testing.T) {
	testcases := []listsTestcase{
		{
			expr:    `(len)`,
			invalid: true,
		},
		{
			expr:    `(len true)`,
			invalid: true,
		},
		{
			expr:    `(len 1)`,
			invalid: true,
		},
		{
			expr:    `(len null)`,
			invalid: true,
		},
		{
			expr:    `(len [] [])`,
			invalid: true,
		},
		{
			expr:     `(len "")`,
			expected: int64(0),
		},
		{
			expr:     `(len " foo ")`,
			expected: int64(5),
		},
		{
			expr:     `(len [])`,
			expected: int64(0),
		},
		{
			expr:     `(len [1 2 3])`,
			expected: int64(3),
		},
		{
			expr:     `(len {})`,
			expected: int64(0),
		},
		{
			expr:     `(len {foo "bar" hello "world"})`,
			expected: int64(2),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestAppendFunction(t *testing.T) {
	testcases := []listsTestcase{
		{
			expr:    `(append)`,
			invalid: true,
		},
		{
			expr:    `(append [])`,
			invalid: true,
		},
		{
			expr:    `(append true 1)`,
			invalid: true,
		},
		{
			expr:    `(append 1 1)`,
			invalid: true,
		},
		{
			expr:    `(append null 1)`,
			invalid: true,
		},
		{
			expr:    `(append {} 1)`,
			invalid: true,
		},
		{
			expr:    `(append {} 1)`,
			invalid: true,
		},
		{
			expr:     `(append [] 1)`,
			expected: []any{int64(1)},
		},
		{
			expr:     `(append [1 2] 3 "foo")`,
			expected: []any{int64(1), int64(2), int64(3), "foo"},
		},
		{
			expr:     `(append [] [])`,
			expected: []any{[]any{}},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestPrependFunction(t *testing.T) {
	testcases := []listsTestcase{
		{
			expr:    `(prepend)`,
			invalid: true,
		},
		{
			expr:    `(prepend [])`,
			invalid: true,
		},
		{
			expr:    `(prepend true 1)`,
			invalid: true,
		},
		{
			expr:    `(prepend 1 1)`,
			invalid: true,
		},
		{
			expr:    `(prepend null 1)`,
			invalid: true,
		},
		{
			expr:    `(prepend {} 1)`,
			invalid: true,
		},
		{
			expr:    `(prepend {} 1)`,
			invalid: true,
		},
		{
			expr:     `(prepend [] 1)`,
			expected: []any{int64(1)},
		},
		{
			expr:     `(prepend [1] 2)`,
			expected: []any{int64(2), int64(1)},
		},
		{
			expr:     `(prepend [1 2] 3 "foo")`,
			expected: []any{int64(3), "foo", int64(1), int64(2)},
		},
		{
			expr:     `(prepend [] [])`,
			expected: []any{[]any{}},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestReverseFunction(t *testing.T) {
	testcases := []listsTestcase{
		{
			expr:    `(reverse)`,
			invalid: true,
		},
		{
			expr:    `(reverse "too" "many")`,
			invalid: true,
		},
		{
			expr:    `(reverse 1)`,
			invalid: true,
		},
		{
			expr:    `(reverse true)`,
			invalid: true,
		},
		{
			expr:    `(reverse null)`,
			invalid: true,
		},
		{
			expr:    `(reverse {})`,
			invalid: true,
		},
		{
			expr:     `(reverse "")`,
			expected: "",
		},
		{
			expr:     `(reverse (concat "" "f" "oo"))`,
			expected: "oof",
		},
		{
			expr:     `(reverse "abcd")`,
			expected: "dcba",
		},
		{
			expr:     `(reverse (reverse "abcd"))`,
			expected: "abcd",
		},
		{
			expr:     `(reverse [])`,
			expected: []any{},
		},
		{
			expr:     `(reverse [1])`,
			expected: []any{int64(1)},
		},
		{
			expr:     `(reverse [1 2 3])`,
			expected: []any{int64(3), int64(2), int64(1)},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestRangeFunction(t *testing.T) {
	testcases := []listsTestcase{
		{
			// missing everything
			expr:    `(range)`,
			invalid: true,
		},
		{
			// missing naming vector
			expr:    `(range [1 2 3])`,
			invalid: true,
		},
		{
			// missing naming vector
			expr:    `(range [1 2 3] (+ 1 2))`,
			invalid: true,
		},
		{
			// naming vector must be 1 or 2 elements long
			expr:    `(range [1 2 3] [] (+ 1 2))`,
			invalid: true,
		},
		{
			// naming vector must be 1 or 2 elements long
			expr:    `(range [1 2 3] [a b c] (+ 1 2))`,
			invalid: true,
		},
		{
			// do not allow numbers in the naming vector
			expr:    `(range [1 2 3] [1 2] (+ 1 2))`,
			invalid: true,
		},
		{
			// do not allow strings in naming vector
			expr:    `(range [1 2 3] ["foo" "bar"] (+ 1 2))`,
			invalid: true,
		},
		{
			// cannot range over non-vectors/objects
			expr:    `(range "invalid" [a] (+ 1 2))`,
			invalid: true,
		},
		{
			// cannot range over non-vectors/objects
			expr:    `(range 5 [a] (+ 1 2))`,
			invalid: true,
		},
		{
			// single simple expression
			expr:     `(range [1 2 3] [a] (+ 1 2))`,
			expected: int64(3),
		},
		{
			// multiple expressions that use a common context
			expr:     `(range [1 2 3] [a] (set $foo $a) (+ $foo 3))`,
			expected: int64(6),
		},
		{
			// count iterations
			expr:     `(range [1 2 3] [loop-var] (set $counter (+ (default (try $counter) 0) 1)))`,
			expected: int64(3),
		},
		{
			// value is bound to desired variable
			expr:     `(range [1 2 3] [a] $a)`,
			expected: int64(3),
		},
		{
			// support loop index variable
			expr:     `(range [1 2 3] [idx var] $idx)`,
			expected: int64(2),
		},
		{
			// support loop index variable
			expr:     `(range [1 2 3] [idx var] $var)`,
			expected: int64(3),
		},
		{
			// variables do not leak outside the range
			expr:    `(range [1 2 3] [idx var] $idx) (+ $var 0)`,
			invalid: true,
		},
		{
			// variables do not leak outside the range
			expr:    `(range [1 2 3] [idx var] $idx) (+ $idx 0)`,
			invalid: true,
		},
		{
			// support ranging over objects
			expr:     `(range {} [key value] $key)`,
			expected: nil,
		},
		{
			expr:     `(range {foo "bar"} [key value] $key)`,
			expected: "foo",
		},
		{
			expr:     `(range {foo "bar"} [key value] $value)`,
			expected: "bar",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestMapFunction(t *testing.T) {
	testcases := []listsTestcase{
		{
			// missing everything
			expr:    `(map)`,
			invalid: true,
		},
		{
			// missing function identifier
			expr:    `(map [1 2 3])`,
			invalid: true,
		},
		{
			// missing naming vector
			expr:    `(map [1 2 3] (+ 1 2))`,
			invalid: true,
		},
		{
			// naming vector must be 1 or 2 elements long
			expr:    `(map [1 2 3] [] (+ 1 2))`,
			invalid: true,
		},
		{
			// naming vector must be 1 or 2 elements long
			expr:    `(map [1 2 3] [a b c] (+ 1 2))`,
			invalid: true,
		},
		{
			// do not allow numbers in the naming vector
			expr:    `(map [1 2 3] [1 2] (+ 1 2))`,
			invalid: true,
		},
		{
			// do not allow strings in naming vector
			expr:    `(map [1 2 3] ["foo" "bar"] (+ 1 2))`,
			invalid: true,
		},
		{
			// cannot map non-vectors/objects
			expr:    `(map "invalid" [a] (+ 1 2))`,
			invalid: true,
		},
		{
			// cannot map non-vectors/objects
			expr:    `(map 5 [a] (+ 1 2))`,
			invalid: true,
		},
		{
			// single simple expression
			expr:     `(map ["foo" "bar"] to-upper)`,
			expected: []any{"FOO", "BAR"},
		},
		{
			expr:     `(map {foo "bar"} to-upper)`,
			expected: map[string]any{"foo": "BAR"},
		},
		{
			// type safety still applies
			expr:    `(map [1] to-upper)`,
			invalid: true,
		},
		{
			// eval expression with variable
			expr:     `(map [1 2 3] [val] (+ $val 3))`,
			expected: []any{int64(4), int64(5), int64(6)},
		},
		{
			// eval with loop index
			expr:     `(map ["foo" "bar"] [idx _] $idx)`,
			expected: []any{int64(0), int64(1)},
		},
		{
			// last expression controls the result
			expr:     `(map [1 2 3] [val] (+ $val 3) "foo")`,
			expected: []any{"foo", "foo", "foo"},
		},
		{
			// multiple expressions that use a common context
			expr:     `(map [1 2 3] [val] (set $foo $val) (+ $foo 3))`,
			expected: []any{int64(4), int64(5), int64(6)},
		},
		{
			// context is even shared across elements
			expr:     `(map ["foo" "bar"] [_] (set $counter (+ (try $counter 0) 1)))`,
			expected: []any{int64(1), int64(2)},
		},
		{
			// variables do not leak outside the range
			expr:    `(map [1 2 3] [idx var] $idx) (+ $var 0)`,
			invalid: true,
		},
		{
			// variables do not leak outside the range
			expr:    `(map [1 2 3] [idx var] $idx) (+ $idx 0)`,
			invalid: true,
		},
		// do not modify the source
		{
			expr:     `(set $foo [1 2 3]) (map $foo [_ __] "bar")`,
			expected: []any{"bar", "bar", "bar"},
		},
		{
			expr:     `(set $foo [1 2 3]) (map $foo [_ __] "bar") $foo`,
			expected: []any{int64(1), int64(2), int64(3)},
		},
		{
			expr:     `(set $foo {foo "bar"}) (map $foo [_ __] "new-value") $foo`,
			expected: map[string]any{"foo": "bar"},
		},
		{
			expr:     `(set $foo ["foo" "bar"]) (map $foo to-upper)`,
			expected: []any{"FOO", "BAR"},
		},
		{
			expr:     `(set $foo ["foo" "bar"]) (map $foo to-upper) $foo`,
			expected: []any{"foo", "bar"},
		},
		{
			expr:     `(set $foo {foo "bar"}) (map $foo to-upper) $foo`,
			expected: map[string]any{"foo": "bar"},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}
