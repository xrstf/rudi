// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/equality"
	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"

	"github.com/google/go-cmp/cmp"
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
					ast.Tuple{Expressions: []ast.Expression{ast.Identifier{Name: "eval"}}},
				},
			},
			invalid: true,
		},
		// (unknownfunc)
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "unknownfunc"},
				},
			},
			invalid: true,
		},
		// (eval "too" "many")
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "eval"},
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
					ast.Identifier{Name: "eval"},
					ast.String("foo"),
				},
			},
			expected: ast.String("foo"),
		},
		// (eval {foo "bar"}).foo
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "eval"},
					ast.ObjectNode{
						Data: []ast.KeyValuePair{
							{
								Key:   ast.Identifier{Name: "foo"},
								Value: ast.String("bar"),
							},
						},
					},
				},
				PathExpression: &ast.PathExpression{
					Steps: []ast.Expression{
						ast.Identifier{Name: "foo"},
					},
				},
			},
			expected: ast.String("bar"),
		},
		// (eval {foo "bar"})[1]
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "eval"},
					ast.ObjectNode{
						Data: []ast.KeyValuePair{
							{
								Key:   ast.Identifier{Name: "foo"},
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
					ast.Identifier{Name: "eval"},
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
					ast.Identifier{Name: "eval"},
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

			ctx := eval.NewContext(doc, nil, dummyFunctions)

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

func TestEvalTupleBangModifier(t *testing.T) {
	testcases := []struct {
		input             ast.Tuple
		variables         types.Variables
		document          any
		expected          ast.Literal
		expectedDocument  any
		expectedVariables types.Variables
		invalid           bool
	}{
		// (set!)
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
				},
			},
			invalid: true,
		},
		// (set! "invalid" "value")
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					ast.String("invalid"),
					ast.String("value"),
				},
			},
			invalid: true,
		},
		// (set! {} "value")
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					ast.ObjectNode{},
					ast.String("value"),
				},
			},
			invalid: true,
		},
		// (set! .[true] "value")
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					ast.Symbol{
						PathExpression: &ast.PathExpression{
							Steps: []ast.Expression{
								ast.Bool(true),
							},
						},
					},
					ast.String("value"),
				},
			},
			invalid: true,
		},
		// (set! .[-1] "value")
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					ast.Symbol{
						PathExpression: &ast.PathExpression{
							Steps: []ast.Expression{
								ast.Number{Value: -1},
							},
						},
					},
					ast.String("value"),
				},
			},
			invalid: true,
		},
		// (set! . "value")
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					ast.Symbol{PathExpression: &ast.PathExpression{}},
					ast.String("value"),
				},
			},
			expected:         ast.String("value"),
			expectedDocument: "value",
		},
		// (set! . "value")
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					ast.Symbol{PathExpression: &ast.PathExpression{}},
					ast.String("value"),
				},
			},
			expected:         ast.String("value"),
			expectedDocument: "value",
		},
		// (set! $val "value")
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					makeVar("myvar", nil),
					ast.String("value"),
				},
			},
			expected: ast.String("value"),
			expectedVariables: types.Variables{
				"myvar": ast.String("value"),
			},
		},
		// (set! .hello "value")
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					ast.Symbol{PathExpression: &ast.PathExpression{
						Steps: []ast.Expression{
							ast.Identifier{Name: "hello"},
						},
					}},
					ast.String("value"),
				},
			},
			document:         map[string]any{"hello": "world", "hei": "verden"},
			expected:         ast.String("value"),
			expectedDocument: map[string]any{"hello": "value", "hei": "verden"},
		},
		// (set! $val.key "value")
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					makeVar("val", &ast.PathExpression{
						Steps: []ast.Expression{
							ast.Identifier{Name: "key"},
						},
					}),
					ast.String("value"),
				},
			},
			variables: types.Variables{
				"val": map[string]any{"foo": "bar", "key": 42},
			},
			expected: ast.String("value"),
			expectedVariables: types.Variables{
				"myvar": map[string]any{"foo": "bar", "key": "value"},
			},
		},
		// (set! .hei (set! .hello "value"))
		{
			input: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					ast.Symbol{PathExpression: &ast.PathExpression{
						Steps: []ast.Expression{
							ast.Identifier{Name: "hei"},
						},
					}},
					ast.Tuple{
						Expressions: []ast.Expression{
							ast.Identifier{Name: "set", Bang: true},
							ast.Symbol{PathExpression: &ast.PathExpression{
								Steps: []ast.Expression{
									ast.Identifier{Name: "hello"},
								},
							}},
							ast.String("value"),
						},
					},
				},
			},
			document:         map[string]any{"hello": "world", "hei": "verden"},
			expected:         ast.String("value"),
			expectedDocument: map[string]any{"hello": "value", "hei": "value"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.input.String(), func(t *testing.T) {
			doc, err := eval.NewDocument(tc.document)
			if err != nil {
				t.Fatalf("Failed to create test document: %v", err)
			}

			ctx := eval.NewContext(doc, tc.variables, dummyFunctions)

			resultCtx, value, err := eval.EvalTuple(ctx, tc.input)
			if err != nil {
				if !tc.invalid {
					t.Fatalf("Failed to run: %v", err)
				}

				return
			}

			if tc.invalid {
				t.Fatalf("Should not have been able to run, but got: %v (%T)", value, value)
			}

			// compare return value

			returned, ok := value.(ast.Literal)
			if !ok {
				t.Errorf("EvalTuple returned unexpected type %T", value)
			} else {
				equal, err := equality.StrictEqual(tc.expected, returned)
				if err != nil {
					t.Errorf("Could not compare result: %v", err)
				} else if !equal {
					t.Fatalf("Expected result value %v (%T), but got %v (%T)", tc.expected, tc.expected, value, value)
				}
			}

			// compare expected document

			resultDoc := resultCtx.GetDocument().Data()

			unwrappedDoc, err := types.UnwrapType(resultDoc)
			if err != nil {
				t.Errorf("Failed to unwrap document: %v", err)
			} else if !cmp.Equal(tc.expectedDocument, unwrappedDoc) {
				t.Fatalf("Expected document %v (%T), but got %v (%T)", tc.expectedDocument, tc.expectedDocument, unwrappedDoc, unwrappedDoc)
			}
		})
	}
}
