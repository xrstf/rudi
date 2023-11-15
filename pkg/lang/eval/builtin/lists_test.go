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
