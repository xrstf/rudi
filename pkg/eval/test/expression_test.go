// SPDX-FileCopyrightText: 2023 Christoph Mewes
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
			Expected: ast.Null{},
		},
		{
			AST:      ast.Bool(true),
			Expected: ast.Bool(true),
		},
		{
			AST:      ast.String("foo"),
			Expected: ast.String("foo"),
		},
		{
			AST:      ast.Number{Value: 1},
			Expected: ast.Number{Value: 1},
		},
		{
			AST:      ast.Object{Data: map[string]any{"foo": "bar"}},
			Expected: ast.Object{Data: map[string]any{"foo": "bar"}},
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
			Expected: ast.Object{Data: map[string]any{"foo": "bar"}},
		},
		{
			AST:      ast.Vector{Data: []any{"foo", 1}},
			Expected: ast.Vector{Data: []any{"foo", 1}},
		},
		{
			AST: ast.VectorNode{
				Expressions: []ast.Expression{
					ast.String("foo"),
					ast.Number{Value: 1},
				},
			},
			Expected: ast.Vector{Data: []any{"foo", 1}},
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = dummyFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}
