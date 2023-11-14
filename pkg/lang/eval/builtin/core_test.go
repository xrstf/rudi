// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"
)

type coreTestcase struct {
	expr     string
	expected any
	invalid  bool
}

func (tc *coreTestcase) Test(t *testing.T) {
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

	if result != tc.expected {
		t.Fatalf("Expected %v (%T), but got %v (%T)", tc.expected, tc.expected, result, result)
	}
}

func TestIfFunction(t *testing.T) {
	testcases := []coreTestcase{
		{
			expr:    `(if)`,
			invalid: true,
		},
		{
			expr:    `(if true)`,
			invalid: true,
		},
		{
			expr:    `(if true "yes" "no" "extra")`,
			invalid: true,
		},
		{
			expr:     `(if true 3)`,
			expected: int64(3),
		},
		{
			expr:     `(if (eq 1 1) 3)`,
			expected: int64(3),
		},
		{
			expr:     `(if (eq 1 2) 3)`,
			expected: nil,
		},
		{
			expr:     `(if (eq 1 2) "yes" "else")`,
			expected: "else",
		},
		{
			expr:     `(if false "yes" (+ 1 4))`,
			expected: int64(5),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestDoFunction(t *testing.T) {
	testcases := []coreTestcase{
		{
			expr:    `(do)`,
			invalid: true,
		},
		{
			expr:     `(do 3)`,
			expected: int64(3),
		},

		// test that the runtime context is inherited from one step to another
		{
			expr:     `(do (set $var "foo") $var)`,
			expected: "foo",
		},
		{
			expr:     `(do (set $var "foo") $var (set $var "new") $var)`,
			expected: "new",
		},

		// test that the runtime context doesn't leak
		{
			expr:     `(set $var "outer") (do (set $var "inner")) (concat $var [1 2])`,
			expected: "1outer2",
		},
		{
			expr:    `(do (set $var "inner")) (concat $var [1 2])`,
			invalid: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestDefaultFunction(t *testing.T) {
	testcases := []coreTestcase{
		{
			expr:    `(default)`,
			invalid: true,
		},
		{
			expr:    `(default true)`,
			invalid: true,
		},
		{
			expr:     `(default null 3)`,
			expected: int64(3),
		},

		// coalescing should be applied

		{
			expr:     `(default false 3)`,
			expected: int64(3),
		},
		{
			expr:     `(default [] 3)`,
			expected: int64(3),
		},

		// errors are not swallowed

		{
			expr:    `(default (eq 3 "foo") 3)`,
			invalid: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestTryFunction(t *testing.T) {
	testcases := []coreTestcase{
		{
			expr:    `(try)`,
			invalid: true,
		},
		{
			expr:     `(try (+ 1 2))`,
			expected: int64(3),
		},

		// coalescing should be not applied

		{
			expr:     `(try false)`,
			expected: false,
		},
		{
			expr:     `(try null)`,
			expected: nil,
		},
		{
			expr:     `(try null "fallback")`,
			expected: nil,
		},

		// swallow errors

		{
			expr:     `(try (eq 3 "foo"))`,
			expected: nil,
		},
		{
			expr:     `(try (eq 3 "foo") "fallback")`,
			expected: "fallback",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}
