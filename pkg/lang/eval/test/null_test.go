// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval"
)

func TestEvalNull(t *testing.T) {
	testcases := []struct {
		input    ast.Null
		expected ast.Null
	}{
		{
			input:    ast.Null{},
			expected: ast.Null{},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.input.String(), func(t *testing.T) {
			doc, err := eval.NewDocument(nil)
			if err != nil {
				t.Fatalf("Failed to create test document: %v", err)
			}

			ctx := eval.NewContext(doc, nil, nil)

			_, value, err := eval.EvalNull(ctx, testcase.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			returned, ok := value.(ast.Null)
			if !ok {
				t.Fatalf("EvalNull returned unexpected type %T", value)
			}

			if !returned.Equal(testcase.expected) {
				t.Fatal("Result does not match expectation.")
			}
		})
	}
}
