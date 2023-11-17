// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval"
)

func TestEvalNumber(t *testing.T) {
	testcases := []struct {
		input    ast.Number
		expected ast.Number
	}{
		{
			input:    ast.Number{Value: 0},
			expected: ast.Number{Value: 0},
		},
		{
			input:    ast.Number{Value: 1},
			expected: ast.Number{Value: 1},
		},
		{
			input:    ast.Number{Value: 3.14},
			expected: ast.Number{Value: 3.14},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.input.String(), func(t *testing.T) {
			doc, err := eval.NewDocument(nil)
			if err != nil {
				t.Fatalf("Failed to create test document: %v", err)
			}

			ctx := eval.NewContext(doc, nil, nil)

			_, value, err := eval.EvalNumber(ctx, testcase.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			returned, ok := value.(ast.Number)
			if !ok {
				t.Fatalf("EvalNumber returned unexpected type %T", value)
			}

			if !returned.Equal(testcase.expected) {
				t.Fatal("Result does not match expectation.")
			}
		})
	}
}
