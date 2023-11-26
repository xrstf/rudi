// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/testutil"
)

func makeSymbol(name string, path *ast.PathExpression) ast.Symbol {
	sym := ast.Symbol{
		PathExpression: path,
	}

	if len(name) > 0 {
		variable := ast.Variable(name)
		sym.Variable = &variable
	}

	return sym
}

func TestEvalSymbol(t *testing.T) {
	testcases := []testutil.Testcase{
		// <utterly invalid Symbol>
		{
			AST:     ast.Symbol{},
			Invalid: true,
		},
		// $undefined
		{
			AST:     makeSymbol("undefined", nil),
			Invalid: true,
		},
		// $var
		{
			AST: makeSymbol("var", nil),
			Variables: types.Variables{
				"var": ast.String("foo"),
			},
			Expected: ast.String("foo"),
		},
		// $native
		{
			AST: makeSymbol("native", nil),
			Variables: types.Variables{
				"native": "foo",
			},
			Expected: ast.String("foo"),
		},
		// $var.foo
		{
			AST: makeSymbol("var", &ast.PathExpression{Steps: []ast.Expression{ast.Identifier{Name: "foo"}}}),
			Variables: types.Variables{
				"var": map[string]any{
					"foo": ast.String("foobar"),
				},
			},
			Expected: ast.String("foobar"),
		},
		// $aVector.foo
		{
			AST: makeSymbol("aVector", &ast.PathExpression{Steps: []ast.Expression{ast.Identifier{Name: "foo"}}}),
			Variables: types.Variables{
				"var": ast.Vector{
					Data: []any{ast.String("first")},
				},
			},
			Invalid: true,
		},
		// $var[1]
		{
			AST: makeSymbol("var", &ast.PathExpression{Steps: []ast.Expression{ast.Number{Value: 1}}}),
			Variables: types.Variables{
				"var": ast.Vector{
					Data: []any{
						ast.String("first"),
						ast.String("second"),
					},
				},
			},
			Expected: ast.String("second"),
		},
		// $aString[1]
		{
			AST: makeSymbol("aString", &ast.PathExpression{Steps: []ast.Expression{ast.Number{Value: 1}}}),
			Variables: types.Variables{
				"var": ast.String("bar"),
			},
			Invalid: true,
		},
		// .
		{
			AST:      makeSymbol("", &ast.PathExpression{}),
			Expected: ast.Null{},
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = dummyFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}
