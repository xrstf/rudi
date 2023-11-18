// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/otto/pkg/eval"
	"go.xrstf.de/otto/pkg/lang/ast"
)

func TestEvalBool(t *testing.T) {
	testcases := []struct {
		input    ast.Bool
		expected ast.Bool
	}{
		{
			input:    ast.Bool(true),
			expected: ast.Bool(true),
		},
		{
			input:    ast.Bool(false),
			expected: ast.Bool(false),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.input.String(), func(t *testing.T) {
			doc, err := eval.NewDocument(nil)
			if err != nil {
				t.Fatalf("Failed to create test document: %v", err)
			}

			ctx := eval.NewContext(doc, nil, nil)

			_, value, err := eval.EvalBool(ctx, testcase.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			returned, ok := value.(ast.Bool)
			if !ok {
				t.Fatalf("EvalBool returned unexpected type %T", value)
			}

			if !returned.Equal(testcase.expected) {
				t.Fatal("Result does not match expectation.")
			}
		})
	}
}
