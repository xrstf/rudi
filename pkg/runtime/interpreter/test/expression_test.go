// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/testutil"
)

func TestEvalExpression(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			AST:      ast.Null{},
			Expected: nil,
		},
		{
			AST:      ast.Bool(true),
			Expected: true,
		},
		{
			AST:      ast.String("foo"),
			Expected: "foo",
		},
		{
			AST:      ast.Number{Value: 1},
			Expected: 1,
		},
		{
			AST: ast.ObjectNode{
				Data: []ast.KeyValuePair{
					{
						Key:   ast.Identifier{Name: "foo"},
						Value: ast.String("bar"),
					},
				},
			},
			Expected: map[string]any{"foo": "bar"},
		},
		{
			AST: ast.VectorNode{
				Expressions: []ast.Expression{
					ast.String("foo"),
					ast.Number{Value: 1},
				},
			},
			Expected: []any{"foo", 1},
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = dummyFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}
