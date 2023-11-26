// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"

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
			Expected:   "foo",
		},
		{
			Expression: `(concat "-" "foo" "bar" "test")`,
			Expected:   "foo-bar-test",
		},
		{
			Expression: `(concat "" "foo" "bar")`,
			Expected:   "foobar",
		},
		{
			Expression: `(concat "" ["foo" "bar"])`,
			Expected:   "foobar",
		},
		{
			Expression: `(concat "-" ["foo" "bar"] "test" ["suffix"])`,
			Expected:   "foo-bar-test-suffix",
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
			Expected:   []any{},
		},
		{
			Expression: `(split "g" "")`,
			Expected:   []any{""},
		},
		{
			Expression: `(split "g" "foo")`,
			Expected:   []any{"foo"},
		},
		{
			Expression: `(split "-" "foo-bar-test-")`,
			Expected:   []any{"foo", "bar", "test", ""},
		},
		{
			Expression: `(split "" "foobar")`,
			Expected:   []any{"f", "o", "o", "b", "a", "r"},
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
			Expected:   "",
		},
		{
			Expression: `(to-upper " TeSt ")`,
			Expected:   " TEST ",
		},
		{
			Expression: `(to-upper " test ")`,
			Expected:   " TEST ",
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
			Expected:   "",
		},
		{
			Expression: `(to-lower " TeSt ")`,
			Expected:   " test ",
		},
		{
			Expression: `(to-lower " TEST ")`,
			Expected:   " test ",
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}
