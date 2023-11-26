// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/testutil"
)

func TestEvalVectorNode(t *testing.T) {
	testcases := []testutil.Testcase{
		// []
		{
			AST:      ast.VectorNode{},
			Expected: ast.Vector{},
		},
		// [identifier]
		{
			AST: ast.VectorNode{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "identifier"},
				},
			},
			Invalid: true,
		},
		// [true "foo" (eval "evaled")]
		{
			AST: ast.VectorNode{
				Expressions: []ast.Expression{
					ast.Bool(true),
					ast.String("foo"),
					ast.Tuple{
						Expressions: []ast.Expression{
							ast.Identifier{Name: "eval"},
							ast.String("evaled"),
						},
					},
				},
			},
			Expected: ast.Vector{
				Data: []any{
					true,
					ast.String("foo"),
					ast.String("evaled"),
				},
			},
		},
		// [true "foo"][1]
		{
			AST: ast.VectorNode{
				Expressions: []ast.Expression{
					ast.Bool(true),
					ast.String("foo"),
				},
				PathExpression: &ast.PathExpression{
					Steps: []ast.Expression{
						ast.Number{Value: 1},
					},
				},
			},
			Expected: ast.String("foo"),
		},
		// ["foo"][1]
		{
			AST: ast.VectorNode{
				Expressions: []ast.Expression{
					ast.String("foo"),
				},
				PathExpression: &ast.PathExpression{
					Steps: []ast.Expression{
						ast.Number{Value: 1},
					},
				},
			},
			Invalid: true,
		},
		// ["foo"].ident
		{
			AST: ast.VectorNode{
				Expressions: []ast.Expression{
					ast.String("foo"),
				},
				PathExpression: &ast.PathExpression{
					Steps: []ast.Expression{
						ast.Identifier{Name: "foo"},
					},
				},
			},
			Invalid: true,
		},
		// ["foo"][1.2]
		{
			AST: ast.VectorNode{
				Expressions: []ast.Expression{
					ast.String("foo"),
				},
				PathExpression: &ast.PathExpression{
					Steps: []ast.Expression{
						ast.Number{Value: 1.2},
					},
				},
			},
			Invalid: true,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = dummyFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}
