// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/equality"
	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func TestEvalExpression(t *testing.T) {
	testcases := []struct {
		input    ast.Expression
		expected ast.Literal
		invalid  bool
	}{
		{
			input:    ast.Null{},
			expected: ast.Null{},
		},
		{
			input:    ast.Bool(true),
			expected: ast.Bool(true),
		},
		{
			input:    ast.String("foo"),
			expected: ast.String("foo"),
		},
		{
			input:    ast.Number{Value: 1},
			expected: ast.Number{Value: 1},
		},
		{
			input:    ast.Object{Data: map[string]any{"foo": "bar"}},
			expected: ast.Object{Data: map[string]any{"foo": "bar"}},
		},
		{
			input: ast.ObjectNode{
				Data: []ast.KeyValuePair{
					{
						Key:   ast.Identifier{Name: "foo"},
						Value: ast.String("bar"),
					},
				},
			},
			expected: ast.Object{Data: map[string]any{"foo": "bar"}},
		},
		{
			input:    ast.Vector{Data: []any{"foo", 1}},
			expected: ast.Vector{Data: []any{"foo", 1}},
		},
		{
			input: ast.VectorNode{
				Expressions: []ast.Expression{
					ast.String("foo"),
					ast.Number{Value: 1},
				},
			},
			expected: ast.Vector{Data: []any{"foo", 1}},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.input.String(), func(t *testing.T) {
			doc, err := eval.NewDocument(nil)
			if err != nil {
				t.Fatalf("Failed to create test document: %v", err)
			}

			ctx := eval.NewContext(doc, nil, dummyFunctions)

			_, value, err := eval.EvalExpression(ctx, testcase.input)
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
				t.Fatalf("EvalExpression returned unexpected type %T", value)
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
