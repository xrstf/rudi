// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
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
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(and false)`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(and null)`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(and true false)`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(and true true)`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(and (eq? 1 1) true)`,
			Expected:   ast.Bool(true),
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
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(or false)`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(or null)`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(or true false)`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(or true true)`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(or (eq? 1 1) true)`,
			Expected:   ast.Bool(true),
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
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(not true)`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(not null)`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(not (not (not (not true))))`,
			Expected:   ast.Bool(true),
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}
