// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/testutil"
)

func TestEvalStatement(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			AST:      ast.Statement{Expression: ast.Null{}},
			Expected: nil,
		},
		{
			AST:      ast.Statement{Expression: ast.Bool(true)},
			Expected: true,
		},
		{
			AST:      ast.Statement{Expression: ast.String("foo")},
			Expected: "foo",
		},
		{
			AST:      ast.Statement{Expression: ast.Number{Value: 1}},
			Expected: 1,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = dummyFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}
