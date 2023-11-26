// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/testutil"
)

func TestToBase64Function(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(to-base64)`,
			Invalid:    true,
		},
		{
			Expression: `(to-base64 "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(to-base64 true)`,
			Invalid:    true,
		},
		{
			Expression: `(to-base64 1)`,
			Invalid:    true,
		},
		{
			Expression: `(to-base64 null)`,
			Invalid:    true,
		},
		{
			Expression: `(to-base64 "")`,
			Expected:   ast.String(""),
		},
		{
			Expression: `(to-base64 " ")`,
			Expected:   ast.String("IA=="),
		},
		{
			Expression: `(to-base64 (concat "" "f" "o" "o"))`,
			Expected:   ast.String("Zm9v"),
		},
		{
			Expression: `(to-base64 "test")`,
			Expected:   ast.String("dGVzdA=="),
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestFromBase64Function(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(from-base64)`,
			Invalid:    true,
		},
		{
			Expression: `(from-base64 "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(from-base64 true)`,
			Invalid:    true,
		},
		{
			Expression: `(from-base64 1)`,
			Invalid:    true,
		},
		{
			Expression: `(from-base64 null)`,
			Invalid:    true,
		},
		{
			Expression: `(from-base64 "definitely-not-base64")`,
			Invalid:    true,
		},
		{
			// should be able to recover
			Expression: `(try (from-base64 "definitely-not-base64") "fallback")`,
			Expected:   ast.String("fallback"),
		},
		{
			Expression: `(from-base64 "")`,
			Expected:   ast.String(""),
		},
		{
			Expression: `(from-base64 "IA==")`,
			Expected:   ast.String(" "),
		},
		{
			Expression: `(from-base64 (concat "" "Z" "m" "9" "v"))`,
			Expected:   ast.String("foo"),
		},
		{
			Expression: `(from-base64 "dGVzdA==")`,
			Expected:   ast.String("test"),
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}
