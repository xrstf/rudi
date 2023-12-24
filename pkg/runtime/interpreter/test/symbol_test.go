// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"fmt"
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/pathexpr"
	"go.xrstf.de/rudi/pkg/runtime/types"
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

type customType struct {
	Data string
}

type customObjGetter struct {
	value any
}

var _ pathexpr.ObjectReader = customObjGetter{}

func (g customObjGetter) GetObjectKey(name string) (any, error) {
	if name == "value" {
		return g.value, nil
	}

	return nil, fmt.Errorf("cannot get property %q", name)
}

type customVecGetter struct {
	magic int
	value any
}

var _ pathexpr.VectorReader = customVecGetter{}

func (g customVecGetter) GetVectorItem(index int) (any, error) {
	if index == g.magic {
		return g.value, nil
	}

	return nil, fmt.Errorf("cannot get index %d", index)
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
		// $var with custom data type
		{
			AST: makeSymbol("var", nil),
			Variables: types.Variables{
				"var": customType{Data: "foo"},
			},
			Expected: customType{Data: "foo"},
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
		// $custom.value
		{
			AST: makeSymbol("custom", &ast.PathExpression{Steps: []ast.Expression{ast.Identifier{Name: "value"}}}),
			Variables: types.Variables{
				"custom": customObjGetter{value: "foo"},
			},
			Expected: "foo",
		},
		// $custom[7]
		{
			AST: makeSymbol("custom", &ast.PathExpression{Steps: []ast.Expression{ast.Number{Value: 7}}}),
			Variables: types.Variables{
				"custom": customVecGetter{magic: 7, value: "foo"},
			},
			Expected: "foo",
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = dummyFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}
