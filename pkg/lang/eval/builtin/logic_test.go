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
			expr:    `(and 1)`,
			invalid: true,
		},
		{
			expr:    `(and 1.1)`,
			invalid: true,
		},
		{
			expr:    `(and null)`,
			invalid: true,
		},
		{
			expr:    `(and "")`,
			invalid: true,
		},
		{
			expr:    `(and "nonempty")`,
			invalid: true,
		},
		{
			expr:    `(and {})`,
			invalid: true,
		},
		{
			expr:    `(and {foo "bar"})`,
			invalid: true,
		},
		{
			expr:    `(and [])`,
			invalid: true,
		},
		{
			expr:    `(and ["bar"])`,
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
			expr:     `(and (eq 1 1) true)`,
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
			expr:    `(or 1)`,
			invalid: true,
		},
		{
			expr:    `(or 1.1)`,
			invalid: true,
		},
		{
			expr:    `(or null)`,
			invalid: true,
		},
		{
			expr:    `(or "")`,
			invalid: true,
		},
		{
			expr:    `(or "nonempty")`,
			invalid: true,
		},
		{
			expr:    `(or {})`,
			invalid: true,
		},
		{
			expr:    `(or {foo "bar"})`,
			invalid: true,
		},
		{
			expr:    `(or [])`,
			invalid: true,
		},
		{
			expr:    `(or ["bar"])`,
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
			expr:     `(or (eq 1 1) true)`,
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
			expr:    `(not true true)`,
			invalid: true,
		},
		{
			expr:    `(not 1)`,
			invalid: true,
		},
		{
			expr:    `(not 1.1)`,
			invalid: true,
		},
		{
			expr:    `(not null)`,
			invalid: true,
		},
		{
			expr:    `(not "")`,
			invalid: true,
		},
		{
			expr:    `(not "nonempty")`,
			invalid: true,
		},
		{
			expr:    `(not {})`,
			invalid: true,
		},
		{
			expr:    `(not {foo "bar"})`,
			invalid: true,
		},
		{
			expr:    `(not [])`,
			invalid: true,
		},
		{
			expr:    `(not ["bar"])`,
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
