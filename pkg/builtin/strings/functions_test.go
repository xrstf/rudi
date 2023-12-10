// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package strings

import (
	"testing"

	"go.xrstf.de/rudi/pkg/testutil"
)

func TestLenFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(len)`,
			Invalid:    true,
		},
		{
			Expression: `(len true)`,
			Invalid:    true,
		},
		{
			Expression: `(len 1)`,
			Invalid:    true,
		},
		{
			Expression: `(len [] [])`,
			Invalid:    true,
		},
		{
			// strict coalescing allows null to turn into [] or "", both have len=0
			Expression: `(len null)`,
			Expected:   0,
		},
		{
			Expression: `(len "")`,
			Expected:   0,
		},
		{
			Expression: `(len " foo ")`,
			Expected:   5,
		},
		{
			Expression: `(len [])`,
			Expected:   0,
		},
		{
			Expression: `(len [1 2 3])`,
			Expected:   3,
		},
		{
			Expression: `(len {})`,
			Expected:   0,
		},
		{
			Expression: `(len {foo "bar" hello "world"})`,
			Expected:   2,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestAppendFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(append)`,
			Invalid:    true,
		},
		{
			Expression: `(append [])`,
			Invalid:    true,
		},
		{
			Expression: `(append true 1)`,
			Invalid:    true,
		},
		{
			Expression: `(append 1 1)`,
			Invalid:    true,
		},
		{
			Expression: `(append {} 1)`,
			Invalid:    true,
		},
		{
			// strict coalescing allows null to turn into [] (null could also be "",
			// which would result in "1", but append/prepend prefer the first arg
			// to be a vector)
			Expression: `(append null 1)`,
			Expected:   []any{int64(1)},
		},
		{
			Expression: `(append [] 1)`,
			Expected:   []any{int64(1)},
		},
		{
			Expression: `(append [1 2] 3 "foo")`,
			Expected:   []any{int64(1), int64(2), int64(3), "foo"},
		},
		{
			Expression: `(append [] [])`,
			Expected:   []any{[]any{}},
		},
		{
			Expression: `(append [] "foo")`,
			Expected:   []any{"foo"},
		},
		{
			Expression: `(append "foo" [])`,
			Invalid:    true,
		},
		{
			Expression: `(append "foo" "bar" [])`,
			Invalid:    true,
		},
		{
			Expression: `(append "foo" "bar" "test")`,
			Expected:   "foobartest",
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestPrependFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(prepend)`,
			Invalid:    true,
		},
		{
			Expression: `(prepend [])`,
			Invalid:    true,
		},
		{
			Expression: `(prepend true 1)`,
			Invalid:    true,
		},
		{
			Expression: `(prepend 1 1)`,
			Invalid:    true,
		},
		{
			Expression: `(prepend {} 1)`,
			Invalid:    true,
		},
		{
			// strict coalescing allows null to turn into [] (null could also be "",
			// which would result in "1", but append/prepend prefer the first arg
			// to be a vector)
			Expression: `(prepend null 1)`,
			Expected:   []any{int64(1)},
		},
		{
			Expression: `(prepend [] 1)`,
			Expected:   []any{int64(1)},
		},
		{
			Expression: `(prepend [1] 2)`,
			Expected:   []any{int64(2), int64(1)},
		},
		{
			Expression: `(prepend [1 2] 3 "foo")`,
			Expected:   []any{int64(3), "foo", int64(1), int64(2)},
		},
		{
			Expression: `(prepend [] [])`,
			Expected:   []any{[]any{}},
		},
		{
			Expression: `(prepend [] "foo")`,
			Expected:   []any{"foo"},
		},
		{
			Expression: `(prepend "foo" [])`,
			Invalid:    true,
		},
		{
			Expression: `(prepend "foo" "bar" [])`,
			Invalid:    true,
		},
		{
			Expression: `(prepend "foo" "bar" "test")`,
			Expected:   "bartestfoo",
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestReverseFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(reverse)`,
			Invalid:    true,
		},
		{
			Expression: `(reverse "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(reverse 1)`,
			Invalid:    true,
		},
		{
			Expression: `(reverse true)`,
			Invalid:    true,
		},
		{
			Expression: `(reverse {})`,
			Invalid:    true,
		},
		{
			// strict coalescing allows null to turn into ""
			Expression: `(reverse null)`,
			Expected:   "",
		},
		{
			Expression: `(reverse "")`,
			Expected:   "",
		},
		{
			Expression: `(reverse "abcd")`,
			Expected:   "dcba",
		},
		{
			Expression: `(reverse (reverse "abcd"))`,
			Expected:   "abcd",
		},
		{
			Expression: `(reverse [])`,
			Expected:   []any{},
		},
		{
			Expression: `(reverse [1])`,
			Expected:   []any{int64(1)},
		},
		{
			Expression: `(reverse [1 2 3])`,
			Expected:   []any{int64(3), int64(2), int64(1)},
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

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
