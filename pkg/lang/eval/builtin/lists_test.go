// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"
)

type listsTestcase struct {
	expr     string
	expected any
	invalid  bool
}

func (tc *listsTestcase) Test(t *testing.T) {
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
			expected: 0,
		},
		{
			expr:     `(len " foo ")`,
			expected: 5,
		},
		{
			expr:     `(len [])`,
			expected: 0,
		},
		{
			expr:     `(len [1 2 3])`,
			expected: 3,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}
