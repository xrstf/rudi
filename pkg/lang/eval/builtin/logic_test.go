// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"
)

type logicTestcase struct {
	expr     string
	expected any
	invalid  bool
}

func (tc *logicTestcase) Test(t *testing.T) {
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

func TestAndFunction(t *testing.T) {
	testcases := []logicTestcase{
		{
			expr:    `(and)`,
			invalid: true,
		},
		{
			expr:     `(and true)`,
			expected: true,
		},
		{
			expr:     `(and false)`,
			expected: false,
		},
		{
			expr:     `(and true false)`,
			expected: false,
		},
		{
			expr:     `(and true true)`,
			expected: true,
		},
		{
			expr:     `(and true 1 "nonempty" 3.1)`,
			expected: true,
		},
		{
			expr:     `(and 0)`,
			expected: false,
		},
		{
			expr:     `(and 0.0)`,
			expected: false,
		},
		{
			expr:     `(and "")`,
			expected: false,
		},
		{
			expr:     `(and null)`,
			expected: false,
		},
		{
			expr:     `(and (or true false) true)`,
			expected: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestOrFunction(t *testing.T) {
	testcases := []logicTestcase{
		{
			expr:    `(or)`,
			invalid: true,
		},
		{
			expr:     `(or true)`,
			expected: true,
		},
		{
			expr:     `(or false)`,
			expected: false,
		},
		{
			expr:     `(or true false)`,
			expected: true,
		},
		{
			expr:     `(or true true)`,
			expected: true,
		},
		{
			expr:     `(or 1)`,
			expected: true,
		},
		{
			expr:     `(or "nonempty")`,
			expected: true,
		},
		{
			expr:     `(or 3.1)`,
			expected: true,
		},
		{
			expr:     `(or 0)`,
			expected: false,
		},
		{
			expr:     `(or 0.0)`,
			expected: false,
		},
		{
			expr:     `(or "")`,
			expected: false,
		},
		{
			expr:     `(or null)`,
			expected: false,
		},
		{
			expr:     `(or (or true false) true)`,
			expected: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestNotFunction(t *testing.T) {
	testcases := []logicTestcase{
		{
			expr:    `(not)`,
			invalid: true,
		},
		{
			expr:     `(not false)`,
			expected: true,
		},
		{
			expr:     `(not true)`,
			expected: false,
		},
		{
			expr:     `(not (not (not (not true))))`,
			expected: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}
