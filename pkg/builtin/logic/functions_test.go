// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package logic

import (
	"testing"

	"go.xrstf.de/rudi/pkg/testutil"
)

func TestAndFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(and)`,
			Invalid:    true,
		},
		{
			Expression: `(and 1)`,
			Invalid:    true,
		},
		{
			Expression: `(and 1.1)`,
			Invalid:    true,
		},
		{
			Expression: `(and "")`,
			Invalid:    true,
		},
		{
			Expression: `(and "nonempty")`,
			Invalid:    true,
		},
		{
			Expression: `(and {})`,
			Invalid:    true,
		},
		{
			Expression: `(and {foo "bar"})`,
			Invalid:    true,
		},
		{
			Expression: `(and [])`,
			Invalid:    true,
		},
		{
			Expression: `(and ["bar"])`,
			Invalid:    true,
		},
		{
			Expression: `(and true)`,
			Expected:   true,
		},
		{
			Expression: `(and false)`,
			Expected:   false,
		},
		{
			Expression: `(and null)`,
			Expected:   false,
		},
		{
			Expression: `(and true false)`,
			Expected:   false,
		},
		{
			Expression: `(and true true)`,
			Expected:   true,
		},
		{
			Expression: `(and (not false) true)`,
			Expected:   true,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestOrFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(or)`,
			Invalid:    true,
		},
		{
			Expression: `(or 1)`,
			Invalid:    true,
		},
		{
			Expression: `(or 1.1)`,
			Invalid:    true,
		},
		{
			Expression: `(or "")`,
			Invalid:    true,
		},
		{
			Expression: `(or "nonempty")`,
			Invalid:    true,
		},
		{
			Expression: `(or {})`,
			Invalid:    true,
		},
		{
			Expression: `(or {foo "bar"})`,
			Invalid:    true,
		},
		{
			Expression: `(or [])`,
			Invalid:    true,
		},
		{
			Expression: `(or ["bar"])`,
			Invalid:    true,
		},
		{
			Expression: `(or true)`,
			Expected:   true,
		},
		{
			Expression: `(or false)`,
			Expected:   false,
		},
		{
			Expression: `(or null)`,
			Expected:   false,
		},
		{
			Expression: `(or true false)`,
			Expected:   true,
		},
		{
			Expression: `(or true true)`,
			Expected:   true,
		},
		{
			Expression: `(or (not false) true)`,
			Expected:   true,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestNotFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(not)`,
			Invalid:    true,
		},
		{
			Expression: `(not true true)`,
			Invalid:    true,
		},
		{
			Expression: `(not 1)`,
			Invalid:    true,
		},
		{
			Expression: `(not 1.1)`,
			Invalid:    true,
		},
		{
			Expression: `(not "")`,
			Invalid:    true,
		},
		{
			Expression: `(not "nonempty")`,
			Invalid:    true,
		},
		{
			Expression: `(not {})`,
			Invalid:    true,
		},
		{
			Expression: `(not {foo "bar"})`,
			Invalid:    true,
		},
		{
			Expression: `(not [])`,
			Invalid:    true,
		},
		{
			Expression: `(not ["bar"])`,
			Invalid:    true,
		},
		{
			Expression: `(not false)`,
			Expected:   true,
		},
		{
			Expression: `(not true)`,
			Expected:   false,
		},
		{
			Expression: `(not null)`,
			Expected:   true,
		},
		{
			Expression: `(not (not (not (not true))))`,
			Expected:   true,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}
