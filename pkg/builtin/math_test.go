// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"

	"go.xrstf.de/rudi/pkg/testutil"
)

func TestSumFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(+)`,
			Invalid:    true,
		},
		{
			Expression: `(+ 1)`,
			Invalid:    true,
		},
		{
			Expression: `(+ 1 "1")`,
			Invalid:    true,
		},
		{
			Expression: `(+ 1 "foo")`,
			Invalid:    true,
		},
		{
			Expression: `(+ 1 [])`,
			Invalid:    true,
		},
		{
			Expression: `(+ 1 {})`,
			Invalid:    true,
		},
		{
			Expression: `(+ 1 2)`,
			Expected:   int64(3),
		},
		{
			Expression: `(+ 1 -2 5)`,
			Expected:   int64(4),
		},
		{
			Expression: `(+ 1 1.5)`,
			Expected:   float64(2.5),
		},
		{
			Expression: `(+ 1 1.5 (+ 1 2))`,
			Expected:   float64(5.5),
		},
		{
			Expression: `(+ 0 0.0 -5.6)`,
			Expected:   float64(-5.6),
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = AllFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestMinusFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(-)`,
			Invalid:    true,
		},
		{
			Expression: `(- 1)`,
			Invalid:    true,
		},
		{
			Expression: `(- 1 "foo")`,
			Invalid:    true,
		},
		{
			Expression: `(- 1 [])`,
			Invalid:    true,
		},
		{
			Expression: `(- 1 {})`,
			Invalid:    true,
		},
		{
			Expression: `(- 1 2)`,
			Expected:   int64(-1),
		},
		{
			Expression: `(- 1 -2 5)`,
			Expected:   int64(-2),
		},
		{
			Expression: `(- 1 1.5)`,
			Expected:   float64(-0.5),
		},
		{
			Expression: `(- 1 1.5 (- 1 2))`,
			Expected:   float64(0.5),
		},
		{
			Expression: `(- 0 0.0 -5.6)`,
			Expected:   float64(5.6),
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = AllFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestMultiplyFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(*)`,
			Invalid:    true,
		},
		{
			Expression: `(* 1)`,
			Invalid:    true,
		},
		{
			Expression: `(* 1 "foo")`,
			Invalid:    true,
		},
		{
			Expression: `(* 1 [])`,
			Invalid:    true,
		},
		{
			Expression: `(* 1 {})`,
			Invalid:    true,
		},
		{
			Expression: `(* 1 2)`,
			Expected:   int64(2),
		},
		{
			Expression: `(* 1 -2 5)`,
			Expected:   int64(-10),
		},
		{
			Expression: `(* 2 -1.5)`,
			Expected:   float64(-3.0),
		},
		{
			Expression: `(* 1 1.5 (* 1 2))`,
			Expected:   float64(3.0),
		},
		{
			Expression: `(* 0 0.0 -5.6)`,
			Expected:   float64(0),
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = AllFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestDivideFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(/)`,
			Invalid:    true,
		},
		{
			Expression: `(/ 1)`,
			Invalid:    true,
		},
		{
			Expression: `(/ 1 "foo")`,
			Invalid:    true,
		},
		{
			Expression: `(/ 1 [])`,
			Invalid:    true,
		},
		{
			Expression: `(/ 1 {})`,
			Invalid:    true,
		},
		{
			Expression: `(/ 1 2)`,
			Expected:   float64(0.5),
		},
		{
			Expression: `(/ 1 -2 5)`,
			Expected:   float64(-0.1),
		},
		{
			Expression: `(/ 2 -1.5)`,
			Expected:   float64(-1.33333333333333333333),
		},
		{
			Expression: `(/ 1 1.5 (/ 1 2))`,
			Expected:   float64(1.33333333333333333333),
		},
		{
			Expression: `(/ 0 0.0 -5.6)`,
			Invalid:    true,
		},
		{
			Expression: `(/ 1 0)`,
			Invalid:    true,
		},
		{
			Expression: `(/ 1 2 0.0)`,
			Invalid:    true,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = AllFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}
