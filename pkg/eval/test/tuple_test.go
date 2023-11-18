// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/otto/pkg/equality"
	"go.xrstf.de/otto/pkg/eval"
	"go.xrstf.de/otto/pkg/lang/ast"
)

func TestEvalTuple(t *testing.T) {
	testcases := []struct {
		input    ast.Tuple
		expected ast.Literal
		invalid  bool
	}{
		{
			input:   ast.Tuple{},
			invalid: true,
		},
		// (true)
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Bool(true),
				},
			},
			invalid: true,
		},
		// ("invalid")
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.String("invalid"),
				},
			},
			invalid: true,
		},
		// (1)
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Number{Value: 1},
				},
			},
			invalid: true,
		},
		// ((eval))
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Tuple{Expressions: []ast.Expression{ast.Identifier("eval")}},
				},
			},
			invalid: true,
		},
		// (eval "too" "many")
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier("eval"),
					ast.String("too"),
					ast.String("many"),
				},
			},
			invalid: true,
		},
		// (eval "foo")
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier("eval"),
					ast.String("foo"),
				},
			},
			expected: ast.String("foo"),
		},
		// (eval {foo "bar"}).foo
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier("eval"),
					ast.ObjectNode{
						Data: []ast.KeyValuePair{
							{
								Key:   ast.Identifier("foo"),
								Value: ast.String("bar"),
							},
						},
					},
				},
				PathExpression: &ast.PathExpression{
					Steps: []ast.Expression{
						ast.Identifier("foo"),
					},
				},
			},
			expected: ast.String("bar"),
		},
		// (eval {foo "bar"})[1]
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier("eval"),
					ast.ObjectNode{
						Data: []ast.KeyValuePair{
							{
								Key:   ast.Identifier("foo"),
								Value: ast.String("bar"),
							},
						},
					},
				},
				PathExpression: &ast.PathExpression{
					Steps: []ast.Expression{
						ast.Number{Value: 1},
					},
				},
			},
			invalid: true,
		},
		// (eval [1 2])[1]
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier("eval"),
					ast.VectorNode{
						Expressions: []ast.Expression{
							ast.Number{Value: 1},
							ast.Number{Value: 2},
						},
					},
				},
				PathExpression: &ast.PathExpression{
					Steps: []ast.Expression{
						ast.Number{Value: 1},
					},
				},
			},
			expected: ast.Number{Value: 2},
		},
		// (eval [1 2]).invalid
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier("eval"),
					ast.VectorNode{
						Expressions: []ast.Expression{
							ast.Number{Value: 1},
							ast.Number{Value: 2},
						},
					},
				},
				PathExpression: &ast.PathExpression{
					Steps: []ast.Expression{
						ast.String("invalid"),
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

			_, value, err := eval.EvalTuple(ctx, testcase.input)
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
				t.Fatalf("EvalTuple returned unexpected type %T", value)
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
