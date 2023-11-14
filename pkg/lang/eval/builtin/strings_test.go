// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type stringsTestcase struct {
	expr     string
	expected any
	invalid  bool
}

func (tc *stringsTestcase) Test(t *testing.T) {
	t.Helper()

	result, err := runExpression(t, tc.expr, nil)
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

func TestConcatFunction(t *testing.T) {
	testcases := []stringsTestcase{
		{
			expr:    `(concat)`,
			invalid: true,
		},
		{
			expr:    `(concat "foo")`,
			invalid: true,
		},
		{
			expr:    `(concat [] "foo")`,
			invalid: true,
		},
		{
			expr:    `(concat {} "foo")`,
			invalid: true,
		},
		{
			expr:    `(concat "g" {})`,
			invalid: true,
		},
		{
			expr:    `(concat "g" [{}])`,
			invalid: true,
		},
		{
			expr:    `(concat "g" [["foo"]])`,
			invalid: true,
		},
		{
			expr:    `(concat "-" "foo" 1)`,
			invalid: true,
		},
		{
			expr:    `(concat true "foo" "bar")`,
			invalid: true,
		},
		{
			expr:     `(concat "g" "foo")`,
			expected: "foo",
		},
		{
			expr:     `(concat "-" "foo" "bar" "test")`,
			expected: "foo-bar-test",
		},
		{
			expr:     `(concat "" "foo" "bar")`,
			expected: "foobar",
		},
		{
			expr:     `(concat "" ["foo" "bar"])`,
			expected: "foobar",
		},
		{
			expr:     `(concat "-" ["foo" "bar"] "test" ["suffix"])`,
			expected: "foo-bar-test-suffix",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestSplitFunction(t *testing.T) {
	testcases := []stringsTestcase{
		{
			expr:    `(split)`,
			invalid: true,
		},
		{
			expr:    `(split "foo")`,
			invalid: true,
		},
		{
			expr:    `(split [] "foo")`,
			invalid: true,
		},
		{
			expr:    `(split {} "foo")`,
			invalid: true,
		},
		{
			expr:    `(split "g" {})`,
			invalid: true,
		},
		{
			expr:    `(split "g" [{}])`,
			invalid: true,
		},
		{
			expr:     `(split "" "")`,
			expected: []any{},
		},
		{
			expr:     `(split "g" "")`,
			expected: []any{""},
		},
		{
			expr:     `(split "g" "foo")`,
			expected: []any{"foo"},
		},
		{
			expr:     `(split "-" "foo-bar-test-")`,
			expected: []any{"foo", "bar", "test", ""},
		},
		{
			expr:     `(split "" "foobar")`,
			expected: []any{"f", "o", "o", "b", "a", "r"},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}
