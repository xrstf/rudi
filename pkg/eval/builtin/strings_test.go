// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/testutil"
)

func TestConcatFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(concat)`,
			Invalid:    true,
		},
		{
			Expression: `(concat "foo")`,
			Invalid:    true,
		},
		{
			Expression: `(concat [] "foo")`,
			Invalid:    true,
		},
		{
			Expression: `(concat {} "foo")`,
			Invalid:    true,
		},
		{
			Expression: `(concat "g" {})`,
			Invalid:    true,
		},
		{
			Expression: `(concat "g" [{}])`,
			Invalid:    true,
		},
		{
			Expression: `(concat "g" [["foo"]])`,
			Invalid:    true,
		},
		{
			Expression: `(concat "-" "foo" 1)`,
			Invalid:    true,
		},
		{
			Expression: `(concat true "foo" "bar")`,
			Invalid:    true,
		},
		{
			Expression: `(concat "g" "foo")`,
			Expected:   ast.String("foo"),
		},
		{
			Expression: `(concat "-" "foo" "bar" "test")`,
			Expected:   ast.String("foo-bar-test"),
		},
		{
			Expression: `(concat "" "foo" "bar")`,
			Expected:   ast.String("foobar"),
		},
		{
			Expression: `(concat "" ["foo" "bar"])`,
			Expected:   ast.String("foobar"),
		},
		{
			Expression: `(concat "-" ["foo" "bar"] "test" ["suffix"])`,
			Expected:   ast.String("foo-bar-test-suffix"),
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestSplitFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(split)`,
			Invalid:    true,
		},
		{
			Expression: `(split "foo")`,
			Invalid:    true,
		},
		{
			Expression: `(split [] "foo")`,
			Invalid:    true,
		},
		{
			Expression: `(split {} "foo")`,
			Invalid:    true,
		},
		{
			Expression: `(split "g" {})`,
			Invalid:    true,
		},
		{
			Expression: `(split "g" [{}])`,
			Invalid:    true,
		},
		{
			Expression: `(split "" "")`,
			Expected: ast.Vector{
				Data: []any{},
			},
		},
		{
			Expression: `(split "g" "")`,
			Expected: ast.Vector{
				Data: []any{ast.String("")},
			},
		},
		{
			Expression: `(split "g" "foo")`,
			Expected: ast.Vector{
				Data: []any{ast.String("foo")},
			},
		},
		{
			Expression: `(split "-" "foo-bar-test-")`,
			Expected: ast.Vector{
				Data: []any{ast.String("foo"), ast.String("bar"), ast.String("test"), ast.String("")},
			},
		},
		{
			Expression: `(split "" "foobar")`,
			Expected: ast.Vector{
				Data: []any{ast.String("f"), ast.String("o"), ast.String("o"), ast.String("b"), ast.String("a"), ast.String("r")},
			},
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestToUpperFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(to-upper)`,
			Invalid:    true,
		},
		{
			Expression: `(to-upper "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(to-upper true)`,
			Invalid:    true,
		},
		{
			Expression: `(to-upper [])`,
			Invalid:    true,
		},
		{
			Expression: `(to-upper {})`,
			Invalid:    true,
		},
		{
			Expression: `(to-upper "")`,
			Expected:   ast.String(""),
		},
		{
			Expression: `(to-upper " TeSt ")`,
			Expected:   ast.String(" TEST "),
		},
		{
			Expression: `(to-upper " test ")`,
			Expected:   ast.String(" TEST "),
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestToLowerFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(to-lower)`,
			Invalid:    true,
		},
		{
			Expression: `(to-lower "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(to-lower true)`,
			Invalid:    true,
		},
		{
			Expression: `(to-lower [])`,
			Invalid:    true,
		},
		{
			Expression: `(to-lower {})`,
			Invalid:    true,
		},
		{
			Expression: `(to-lower "")`,
			Expected:   ast.String(""),
		},
		{
			Expression: `(to-lower " TeSt ")`,
			Expected:   ast.String(" test "),
		},
		{
			Expression: `(to-lower " TEST ")`,
			Expected:   ast.String(" test "),
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}
