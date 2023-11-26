// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/testutil"
)

func TestEvalNumber(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			AST:      ast.Number{Value: 0},
			Expected: ast.Number{Value: 0},
		},
		{
			AST:      ast.Number{Value: 1},
			Expected: ast.Number{Value: 1},
		},
		{
			AST:      ast.Number{Value: 3.14},
			Expected: ast.Number{Value: 3.14},
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = dummyFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}
