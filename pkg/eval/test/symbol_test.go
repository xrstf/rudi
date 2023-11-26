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
				"var": "foo",
			},
			Expected: "foo",
		},
		// $var.foo
		{
			AST: makeSymbol("var", &ast.PathExpression{Steps: []ast.Expression{ast.Identifier{Name: "foo"}}}),
			Variables: types.Variables{
				"var": map[string]any{
					"foo": "foobar",
				},
			},
			Expected: "foobar",
		},
		// $aVector.foo
		{
			AST: makeSymbol("aVector", &ast.PathExpression{Steps: []ast.Expression{ast.Identifier{Name: "foo"}}}),
			Variables: types.Variables{
				"aVector": []any{"first"},
			},
			Invalid: true,
		},
		// $var[1]
		{
			AST: makeSymbol("var", &ast.PathExpression{Steps: []ast.Expression{ast.Number{Value: 1}}}),
			Variables: types.Variables{
				"var": []any{
					"first",
					"second",
				},
			},
			Expected: "second",
		},
		// $aString[1]
		{
			AST: makeSymbol("aString", &ast.PathExpression{Steps: []ast.Expression{ast.Number{Value: 1}}}),
			Variables: types.Variables{
				"aString": "bar",
			},
			Invalid: true,
		},
		// .
		{
			AST:      makeSymbol("", &ast.PathExpression{}),
			Expected: nil,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = dummyFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}
