// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/testutil"
)

func TestEvalIdentifier(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			AST:     ast.Identifier{Name: "foo"},
			Invalid: true,
		},
		{
			AST:     ast.Identifier{Name: "foo", Bang: true},
			Invalid: true,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = dummyFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}
