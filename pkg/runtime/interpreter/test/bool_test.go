// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/testutil"
)

func TestEvalBool(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			AST:      ast.Bool(true),
			Expected: true,
		},
		{
			AST:      ast.Bool(false),
			Expected: false,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = dummyFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}
