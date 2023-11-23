// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/equality"
	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
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
	testcases := []struct {
		input     ast.Symbol
		expected  ast.Literal
		variables types.Variables
		document  any
		invalid   bool
	}{
		// <utterly invalid Symbol>
		{
			input:   ast.Symbol{},
			invalid: true,
		},
		// $undefined
		{
			input:   makeSymbol("undefined", nil),
			invalid: true,
		},
		// $var
		{
			input: makeSymbol("var", nil),
			variables: types.Variables{
				"var": ast.String("foo"),
			},
			expected: ast.String("foo"),
		},
		// $native
		{
			input: makeSymbol("native", nil),
			variables: types.Variables{
				"native": "foo",
			},
			expected: ast.String("foo"),
		},
		// $var.foo
		{
			input: makeSymbol("var", &ast.PathExpression{Steps: []ast.Expression{ast.Identifier{Name: "foo"}}}),
			variables: types.Variables{
				"var": map[string]any{
					"foo": ast.String("foobar"),
				},
			},
			expected: ast.String("foobar"),
		},
		// $aVector.foo
		{
			input: makeSymbol("aVector", &ast.PathExpression{Steps: []ast.Expression{ast.Identifier{Name: "foo"}}}),
			variables: types.Variables{
				"var": ast.Vector{
					Data: []any{ast.String("first")},
				},
			},
			invalid: true,
		},
		// $var[1]
		{
			input: makeSymbol("var", &ast.PathExpression{Steps: []ast.Expression{ast.Number{Value: 1}}}),
			variables: types.Variables{
				"var": ast.Vector{
					Data: []any{
						ast.String("first"),
						ast.String("second"),
					},
				},
			},
			expected: ast.String("second"),
		},
		// $aString[1]
		{
			input: makeSymbol("aString", &ast.PathExpression{Steps: []ast.Expression{ast.Number{Value: 1}}}),
			variables: types.Variables{
				"var": ast.String("bar"),
			},
			invalid: true,
		},
		// .
		{
			input:    makeSymbol("", &ast.PathExpression{}),
			expected: ast.Null{},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.input.String(), func(t *testing.T) {
			doc, err := eval.NewDocument(testcase.document)
			if err != nil {
				t.Fatalf("Failed to create test document: %v", err)
			}

			ctx := eval.NewContext(doc, testcase.variables, dummyFunctions)

			_, value, err := eval.EvalSymbol(ctx, testcase.input)
			if err != nil {
				if !testcase.invalid {
					t.Fatalf("Failed to run: %v", err)
				}

				return
			}

			if testcase.invalid {
				t.Fatalf("Should not have been able to run, but got: %v (%T)", value, value)
			}

			returned, ok := value.(ast.Literal)
			if !ok {
				t.Fatalf("EvalSymbol returned unexpected type %T", value)
			}

			equal, err := equality.StrictEqual(testcase.expected, returned)
			if err != nil {
				t.Fatalf("Could not compare result: %v", err)
			}

			if !equal {
				t.Fatal("Result does not match expectation.")
			}
		})
	}
}
