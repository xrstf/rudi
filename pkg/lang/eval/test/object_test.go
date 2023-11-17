// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/equality"
	"go.xrstf.de/otto/pkg/lang/eval"
)

func TestEvalObjectNode(t *testing.T) {
	testcases := []struct {
		input    ast.ObjectNode
		expected ast.Literal
	}{
		// {}
		{
			input:    ast.ObjectNode{},
			expected: ast.Object{},
		},
		// {foo "bar"}
		{
			input: ast.ObjectNode{
				Data: []ast.KeyValuePair{
					{
						Key:   ast.Identifier("foo"),
						Value: ast.String("bar"),
					},
				},
			},
			expected: ast.Object{
				Data: map[string]any{
					"foo": ast.String("bar"),
				},
			},
		},
		// {(eval "evaled") (eval "also evaled")}
		{
			input: ast.ObjectNode{
				Data: []ast.KeyValuePair{
					{
						Key: ast.Tuple{
							Expressions: []ast.Expression{
								ast.Identifier("eval"),
								ast.String("evaled"),
							},
						},
						Value: ast.Tuple{
							Expressions: []ast.Expression{
								ast.Identifier("eval"),
								ast.String("also evaled"),
							},
						},
					},
				},
			},
			expected: ast.Object{
				Data: map[string]any{
					"evaled": ast.String("also evaled"),
				},
			},
		},
		// {foo "bar"}.foo
		{
			input: ast.ObjectNode{
				Data: []ast.KeyValuePair{
					{
						Key:   ast.Identifier("foo"),
						Value: ast.String("bar"),
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
	}

	for _, testcase := range testcases {
		t.Run(testcase.input.String(), func(t *testing.T) {
			doc, err := eval.NewDocument(nil)
			if err != nil {
				t.Fatalf("Failed to create test document: %v", err)
			}

			ctx := eval.NewContext(doc, dummyFunctions, nil)

			_, value, err := eval.EvalObjectNode(ctx, testcase.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			returned, ok := value.(ast.Literal)
			if !ok {
				t.Fatalf("EvalObjectNode returned unexpected type %T", value)
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
