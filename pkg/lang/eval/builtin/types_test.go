// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"
)

type typesTestcase struct {
	expr     string
	expected any
	invalid  bool
}

func (tc *typesTestcase) Test(t *testing.T) {
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

func TestToStringFunction(t *testing.T) {
	testcases := []typesTestcase{
		{
			expr:    `(to-string)`,
			invalid: true,
		},
		{
			expr:    `(to-string "too" "many")`,
			invalid: true,
		},
		{
			expr:     `(to-string "foo")`,
			expected: "foo",
		},
		{
			expr:     `(to-string 1)`,
			expected: "1",
		},
		{
			expr:     `(to-string (+ 1 3))`,
			expected: "4",
		},
		{
			expr:     `(to-string 1.5)`,
			expected: "1.500000",
		},
		{
			expr:     `(to-string true)`,
			expected: "true",
		},
		{
			expr:     `(to-string null)`,
			expected: "null",
		},
		{
			expr:    `(to-string [])`,
			invalid: true,
		},
		{
			expr:    `(to-string {})`,
			invalid: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestToIntFunction(t *testing.T) {
	testcases := []typesTestcase{
		{
			expr:    `(to-int)`,
			invalid: true,
		},
		{
			expr:    `(to-int "too" "many")`,
			invalid: true,
		},
		{
			expr:     `(to-int 1)`,
			expected: int64(1),
		},
		{
			expr:     `(to-int "42")`,
			expected: int64(42),
		},
		{
			expr:     `(to-int (+ 1 3))`,
			expected: int64(4),
		},
		{
			expr:    `(to-int 1.5)`,
			invalid: true, // should this be allowed?
		},
		{
			expr:    `(to-int "1.5")`,
			invalid: true, // should this be allowed?
		},
		{
			expr:     `(to-int true)`,
			expected: int64(1),
		},
		{
			expr:     `(to-int false)`,
			expected: int64(0),
		},
		{
			expr:     `(to-int null)`,
			expected: int64(0),
		},
		{
			expr:    `(to-int [])`,
			invalid: true,
		},
		{
			expr:    `(to-int {})`,
			invalid: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestToFloatFunction(t *testing.T) {
	testcases := []typesTestcase{
		{
			expr:    `(to-float)`,
			invalid: true,
		},
		{
			expr:    `(to-float "too" "many")`,
			invalid: true,
		},
		{
			expr:     `(to-float 1)`,
			expected: float64(1),
		},
		{
			expr:     `(to-float (+ 1 3))`,
			expected: float64(4),
		},
		{
			expr:     `(to-float 1.5)`,
			expected: float64(1.5),
		},
		{
			expr:     `(to-float "3")`,
			expected: float64(3),
		},
		{
			expr:     `(to-float "1.5")`,
			expected: float64(1.5),
		},
		{
			expr:     `(to-float true)`,
			expected: float64(1),
		},
		{
			expr:     `(to-float false)`,
			expected: float64(0),
		},
		{
			expr:     `(to-float null)`,
			expected: float64(0),
		},
		{
			expr:    `(to-float [])`,
			invalid: true,
		},
		{
			expr:    `(to-float {})`,
			invalid: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestToBoolFunction(t *testing.T) {
	testcases := []typesTestcase{
		{
			expr:    `(to-bool)`,
			invalid: true,
		},
		{
			expr:    `(to-bool "too" "many")`,
			invalid: true,
		},
		{
			expr:     `(to-bool 1)`,
			expected: true,
		},
		{
			expr:     `(to-bool 0)`,
			expected: false,
		},
		{
			expr:     `(to-bool (+ 1 3))`,
			expected: true,
		},
		{
			expr:     `(to-bool 1.5)`,
			expected: true,
		},
		{
			expr:     `(to-bool 0.0)`,
			expected: false,
		},
		{
			expr:     `(to-bool "3")`,
			expected: true,
		},
		{
			expr:     `(to-bool true)`,
			expected: true,
		},
		{
			expr:     `(to-bool false)`,
			expected: false,
		},
		{
			expr:     `(to-bool null)`,
			expected: false,
		},
		{
			expr:     `(to-bool [])`,
			expected: false,
		},
		{
			expr:     `(to-bool [0])`,
			expected: true,
		},
		{
			expr:     `(to-bool {})`,
			expected: false,
		},
		{
			expr:     `(to-bool {foo "bar"})`,
			expected: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}
