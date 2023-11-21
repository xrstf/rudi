// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/equality"
	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func TestEvalVectorNode(t *testing.T) {
	testcases := []struct {
		input    ast.VectorNode
		expected ast.Literal
		invalid  bool
	}{
		// []
		{
			input:    ast.VectorNode{},
			expected: ast.Vector{},
		},
		// [identifier]
		{
			input: ast.VectorNode{
				Expressions: []ast.Expression{
					ast.Identifier("identifier"),
				},
			},
			invalid: true,
		},
		// [true "foo" (eval "evaled")]
		{
			input: ast.VectorNode{
				Expressions: []ast.Expression{
					ast.Bool(true),
					ast.String("foo"),
					ast.Tuple{
						Expressions: []ast.Expression{
							ast.Identifier("eval"),
							ast.String("evaled"),
						},
					},
				},
			},
			expected: ast.Vector{
				Data: []any{
					true,
					ast.String("foo"),
					ast.String("evaled"),
				},
			},
		},
		// [true "foo"][1]
		{
			input: ast.VectorNode{
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
			expected: ast.String("foo"),
		},
		// ["foo"][1]
		{
			input: ast.VectorNode{
				Expressions: []ast.Expression{
					ast.String("foo"),
				},
				PathExpression: &ast.PathExpression{
					Steps: []ast.Expression{
						ast.Number{Value: 1},
					},
				},
			},
			invalid: true,
		},
		// ["foo"].ident
		{
			input: ast.VectorNode{
				Expressions: []ast.Expression{
					ast.String("foo"),
				},
				PathExpression: &ast.PathExpression{
					Steps: []ast.Expression{
						ast.Identifier("foo"),
					},
				},
			},
			invalid: true,
		},
		// ["foo"][1.2]
		{
			input: ast.VectorNode{
				Expressions: []ast.Expression{
					ast.String("foo"),
				},
				PathExpression: &ast.PathExpression{
					Steps: []ast.Expression{
						ast.Number{Value: 1.2},
					},
				},
			},
			invalid: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.input.String(), func(t *testing.T) {
			doc, err := eval.NewDocument(nil)
			if err != nil {
				t.Fatalf("Failed to create test document: %v", err)
			}

			ctx := eval.NewContext(doc, dummyFunctions, nil)

			_, value, err := eval.EvalVectorNode(ctx, testcase.input)
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
				t.Fatalf("EvalVectorNode returned unexpected type %T", value)
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
