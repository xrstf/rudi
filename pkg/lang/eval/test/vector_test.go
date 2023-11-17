// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/equality"
	"go.xrstf.de/otto/pkg/lang/eval"
)

func TestEvalVectorNode(t *testing.T) {
	testcases := []struct {
		input    ast.VectorNode
		expected ast.Literal
	}{
		// []
		{
			input:    ast.VectorNode{},
			expected: ast.Vector{},
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
				t.Fatalf("Unexpected error: %v", err)
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
