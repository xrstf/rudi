// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"
)

type comparisonsTestcase struct {
	expr     string
	expected any
	invalid  bool
}

func (tc *comparisonsTestcase) Test(t *testing.T) {
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

func TestEqFunction(t *testing.T) {
	testcases := []comparisonsTestcase{
		{
			expr:    `(eq)`,
			invalid: true,
		},
		{
			expr:    `(eq true)`,
			invalid: true,
		},
		{
			expr:     `(eq false true)`,
			expected: false,
		},
		{
			expr:     `(eq true true)`,
			expected: true,
		},
		{
			expr:     `(eq 1 1)`,
			expected: true,
		},
		{
			expr:     `(eq 1 2)`,
			expected: false,
		},
		{
			expr:    `(eq 1 "foo")`,
			invalid: true,
		},
		{
			expr:    `(eq 1 true)`,
			invalid: true,
		},
		{
			expr:     `(eq "foo" "bar")`,
			expected: false,
		},
		{
			expr:     `(eq "foo" "Foo")`,
			expected: false,
		},
		{
			expr:     `(eq "foo" " foo")`,
			expected: false,
		},

		// no support for vector comparisons at all

		{
			expr:    `(eq 1 [])`,
			invalid: true,
		},
		{
			expr:    `(eq [] [])`,
			invalid: true,
		},

		// same for objects

		{
			expr:    `(eq 1 {})`,
			invalid: true,
		},
		{
			expr:    `(eq {} {})`,
			invalid: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}
