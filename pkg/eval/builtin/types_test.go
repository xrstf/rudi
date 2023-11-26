// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
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
			Expected:   ast.String("foo"),
		},
		{
			Expression: `(to-string 1)`,
			Expected:   ast.String("1"),
		},
		{
			Expression: `(to-string (+ 1 3))`,
			Expected:   ast.String("4"),
		},
		{
			Expression: `(to-string 1.5)`,
			Expected:   ast.String("1.5"),
		},
		{
			Expression: `(to-string 1e3)`,
			Expected:   ast.String("1000"),
		},
		{
			Expression: `(to-string true)`,
			Expected:   ast.String("true"),
		},
		{
			Expression: `(to-string null)`,
			Expected:   ast.String("null"),
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
			Expected:   ast.Number{Value: int64(1)},
		},
		{
			Expression: `(to-int "42")`,
			Expected:   ast.Number{Value: int64(42)},
		},
		{
			Expression: `(to-int (+ 1 3))`,
			Expected:   ast.Number{Value: int64(4)},
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
			Expected:   ast.Number{Value: int64(1)},
		},
		{
			Expression: `(to-int false)`,
			Expected:   ast.Number{Value: int64(0)},
		},
		{
			Expression: `(to-int null)`,
			Expected:   ast.Number{Value: int64(0)},
		},
		{
			Expression: `(to-int [])`,
			Invalid:    true,
		},
		{
			Expression: `(to-int {})`,
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
			Expected:   ast.Number{Value: float64(1)},
		},
		{
			Expression: `(to-float (+ 1 3))`,
			Expected:   ast.Number{Value: float64(4)},
		},
		{
			Expression: `(to-float 1.5)`,
			Expected:   ast.Number{Value: float64(1.5)},
		},
		{
			Expression: `(to-float "3")`,
			Expected:   ast.Number{Value: float64(3)},
		},
		{
			Expression: `(to-float "1.5")`,
			Expected:   ast.Number{Value: float64(1.5)},
		},
		{
			Expression: `(to-float true)`,
			Expected:   ast.Number{Value: float64(1)},
		},
		{
			Expression: `(to-float false)`,
			Expected:   ast.Number{Value: float64(0)},
		},
		{
			Expression: `(to-float null)`,
			Expected:   ast.Number{Value: float64(0)},
		},
		{
			Expression: `(to-float [])`,
			Invalid:    true,
		},
		{
			Expression: `(to-float {})`,
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
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(to-bool 0)`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(to-bool (+ 1 3))`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(to-bool 1.5)`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(to-bool 0.0)`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(to-bool "3")`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(to-bool true)`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(to-bool false)`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(to-bool null)`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(to-bool [])`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(to-bool [0])`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(to-bool {})`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(to-bool {foo "bar"})`,
			Expected:   ast.Bool(true),
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
			Expected:   ast.String("number"),
		},
		{
			Expression: `(type-of 0)`,
			Expected:   ast.String("number"),
		},
		{
			Expression: `(type-of (+ 1 3))`,
			Expected:   ast.String("number"),
		},
		{
			Expression: `(type-of 1.5)`,
			Expected:   ast.String("number"),
		},
		{
			Expression: `(type-of 0.0)`,
			Expected:   ast.String("number"),
		},
		{
			Expression: `(type-of "3")`,
			Expected:   ast.String("string"),
		},
		{
			Expression: `(type-of true)`,
			Expected:   ast.String("bool"),
		},
		{
			Expression: `(type-of false)`,
			Expected:   ast.String("bool"),
		},
		{
			Expression: `(type-of null)`,
			Expected:   ast.String("null"),
		},
		{
			Expression: `(type-of [])`,
			Expected:   ast.String("vector"),
		},
		{
			Expression: `(type-of (append [] "test"))`,
			Expected:   ast.String("vector"),
		},
		{
			Expression: `(type-of [0])`,
			Expected:   ast.String("vector"),
		},
		{
			Expression: `(type-of {})`,
			Expected:   ast.String("object"),
		},
		{
			Expression: `(type-of {foo "bar"})`,
			Expected:   ast.String("object"),
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}
