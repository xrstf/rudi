// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"
)

type mathTestcase struct {
	expr     string
	expected any
	invalid  bool
}

func (tc *mathTestcase) Test(t *testing.T) {
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

func TestSumFunction(t *testing.T) {
	testcases := []mathTestcase{
		{
			expr:    `(+)`,
			invalid: true,
		},
		{
			expr:    `(+ 1)`,
			invalid: true,
		},
		{
			expr:    `(+ 1 "1")`,
			invalid: true,
		},
		{
			expr:    `(+ 1 "foo")`,
			invalid: true,
		},
		{
			expr:    `(+ 1 [])`,
			invalid: true,
		},
		{
			expr:    `(+ 1 {})`,
			invalid: true,
		},
		{
			expr:     `(+ 1 2)`,
			expected: int64(3),
		},
		{
			expr:     `(+ 1 -2 5)`,
			expected: int64(4),
		},
		{
			expr:     `(+ 1 1.5)`,
			expected: float64(2.5),
		},
		{
			expr:     `(+ 1 1.5 (+ 1 2))`,
			expected: float64(5.5),
		},
		{
			expr:     `(+ 0 0.0 -5.6)`,
			expected: float64(-5.6),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestMinusFunction(t *testing.T) {
	testcases := []mathTestcase{
		{
			expr:    `(-)`,
			invalid: true,
		},
		{
			expr:    `(- 1)`,
			invalid: true,
		},
		{
			expr:    `(- 1 "foo")`,
			invalid: true,
		},
		{
			expr:    `(- 1 [])`,
			invalid: true,
		},
		{
			expr:    `(- 1 {})`,
			invalid: true,
		},
		{
			expr:     `(- 1 2)`,
			expected: int64(-1),
		},
		{
			expr:     `(- 1 -2 5)`,
			expected: int64(-2),
		},
		{
			expr:     `(- 1 1.5)`,
			expected: float64(-0.5),
		},
		{
			expr:     `(- 1 1.5 (- 1 2))`,
			expected: float64(0.5),
		},
		{
			expr:     `(- 0 0.0 -5.6)`,
			expected: float64(5.6),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestMultiplyFunction(t *testing.T) {
	testcases := []mathTestcase{
		{
			expr:    `(*)`,
			invalid: true,
		},
		{
			expr:    `(* 1)`,
			invalid: true,
		},
		{
			expr:    `(* 1 "foo")`,
			invalid: true,
		},
		{
			expr:    `(* 1 [])`,
			invalid: true,
		},
		{
			expr:    `(* 1 {})`,
			invalid: true,
		},
		{
			expr:     `(* 1 2)`,
			expected: int64(2),
		},
		{
			expr:     `(* 1 -2 5)`,
			expected: int64(-10),
		},
		{
			expr:     `(* 2 -1.5)`,
			expected: float64(-3.0),
		},
		{
			expr:     `(* 1 1.5 (* 1 2))`,
			expected: float64(3.0),
		},
		{
			expr:     `(* 0 0.0 -5.6)`,
			expected: float64(0),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestDivideFunction(t *testing.T) {
	testcases := []mathTestcase{
		{
			expr:    `(/)`,
			invalid: true,
		},
		{
			expr:    `(/ 1)`,
			invalid: true,
		},
		{
			expr:    `(/ 1 "foo")`,
			invalid: true,
		},
		{
			expr:    `(/ 1 [])`,
			invalid: true,
		},
		{
			expr:    `(/ 1 {})`,
			invalid: true,
		},
		{
			expr:     `(/ 1 2)`,
			expected: float64(0.5),
		},
		{
			expr:     `(/ 1 -2 5)`,
			expected: float64(-0.1),
		},
		{
			expr:     `(/ 2 -1.5)`,
			expected: float64(-1.33333333333333333333),
		},
		{
			expr:     `(/ 1 1.5 (/ 1 2))`,
			expected: float64(1.33333333333333333333),
		},
		{
			expr:    `(/ 0 0.0 -5.6)`,
			invalid: true,
		},
		{
			expr:    `(/ 1 0)`,
			invalid: true,
		},
		{
			expr:    `(/ 1 2 0.0)`,
			invalid: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}
