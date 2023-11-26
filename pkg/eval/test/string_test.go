// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/testutil"
)

func TestEvalString(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			AST:      ast.String(""),
			Expected: ast.String(""),
		},
		{
			AST:      ast.String("foo"),
			Expected: ast.String("foo"),
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = dummyFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}
