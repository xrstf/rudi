// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package types

import (
	"testing"

	"go.xrstf.de/rudi/pkg/testutil"
)

func TestToStringFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(to-string)`,
			Invalid:    true,
		},
		{
			Expression: `(to-string "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(to-string "foo")`,
			Expected:   "foo",
		},
		{
			Expression: `(to-string 1)`,
			Expected:   "1",
		},
		{
			Expression: `(to-string 1.5)`,
			Expected:   "1.5",
		},
		{
			Expression: `(to-string 1e3)`,
			Expected:   "1000",
		},
		{
			Expression: `(to-string true)`,
			Expected:   "true",
		},
		{
			Expression: `(to-string false)`,
			Expected:   "false",
		},
		{
			Expression: `(to-string null)`,
			Expected:   "",
		},
		{
			Expression: `(to-string [])`,
			Invalid:    true,
		},
		{
			Expression: `(to-string {})`,
			Invalid:    true,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestToIntFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(to-int)`,
			Invalid:    true,
		},
		{
			Expression: `(to-int "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(to-int 1)`,
			Expected:   int64(1),
		},
		{
			Expression: `(to-int "42")`,
			Expected:   int64(42),
		},
		{
			Expression: `(to-int 1.5)`,
			Invalid:    true,
		},
		{
			Expression: `(to-int "1.5")`,
			Invalid:    true,
		},
		{
			Expression: `(to-int true)`,
			Expected:   int64(1),
		},
		{
			Expression: `(to-int false)`,
			Expected:   int64(0),
		},
		{
			Expression: `(to-int null)`,
			Expected:   int64(0),
		},
		{
			Expression: `(to-int [])`,
			Invalid:    true,
		},
		{
			Expression: `(to-int [0])`,
			Invalid:    true,
		},
		{
			Expression: `(to-int {})`,
			Invalid:    true,
		},
		{
			Expression: `(to-int {"" ""})`,
			Invalid:    true,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestToFloatFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(to-float)`,
			Invalid:    true,
		},
		{
			Expression: `(to-float "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(to-float 1)`,
			Expected:   float64(1),
		},
		{
			Expression: `(to-float 1.5)`,
			Expected:   float64(1.5),
		},
		{
			Expression: `(to-float "3")`,
			Expected:   float64(3),
		},
		{
			Expression: `(to-float "1.5")`,
			Expected:   float64(1.5),
		},
		{
			Expression: `(to-float true)`,
			Expected:   float64(1),
		},
		{
			Expression: `(to-float false)`,
			Expected:   float64(0),
		},
		{
			Expression: `(to-float null)`,
			Expected:   float64(0),
		},
		{
			Expression: `(to-float [])`,
			Invalid:    true,
		},
		{
			Expression: `(to-float [""])`,
			Invalid:    true,
		},
		{
			Expression: `(to-float {})`,
			Invalid:    true,
		},
		{
			Expression: `(to-float {"" ""})`,
			Invalid:    true,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestToBoolFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(to-bool)`,
			Invalid:    true,
		},
		{
			Expression: `(to-bool "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(to-bool 1)`,
			Expected:   true,
		},
		{
			Expression: `(to-bool 0)`,
			Expected:   false,
		},
		{
			Expression: `(to-bool 1.5)`,
			Expected:   true,
		},
		{
			Expression: `(to-bool 0.0)`,
			Expected:   false,
		},
		{
			Expression: `(to-bool "3")`,
			Expected:   true,
		},
		{
			Expression: `(to-bool true)`,
			Expected:   true,
		},
		{
			Expression: `(to-bool false)`,
			Expected:   false,
		},
		{
			Expression: `(to-bool null)`,
			Expected:   false,
		},
		{
			Expression: `(to-bool [])`,
			Expected:   false,
		},
		{
			Expression: `(to-bool [0])`,
			Expected:   true,
		},
		{
			Expression: `(to-bool {})`,
			Expected:   false,
		},
		{
			Expression: `(to-bool {foo "bar"})`,
			Expected:   true,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestTypeOfFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(type-of)`,
			Invalid:    true,
		},
		{
			Expression: `(type-of "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(type-of 1)`,
			Expected:   "number",
		},
		{
			Expression: `(type-of 0)`,
			Expected:   "number",
		},
		{
			Expression: `(type-of 1.5)`,
			Expected:   "number",
		},
		{
			Expression: `(type-of 0.0)`,
			Expected:   "number",
		},
		{
			Expression: `(type-of "3")`,
			Expected:   "string",
		},
		{
			Expression: `(type-of true)`,
			Expected:   "bool",
		},
		{
			Expression: `(type-of false)`,
			Expected:   "bool",
		},
		{
			Expression: `(type-of null)`,
			Expected:   "null",
		},
		{
			Expression: `(type-of [])`,
			Expected:   "vector",
		},
		{
			Expression: `(type-of [0])`,
			Expected:   "vector",
		},
		{
			Expression: `(type-of {})`,
			Expected:   "object",
		},
		{
			Expression: `(type-of {foo "bar"})`,
			Expected:   "object",
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}
