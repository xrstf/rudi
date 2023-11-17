// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval"
)

func TestEvalIdentifier(t *testing.T) {
	testcases := []struct {
		input ast.Identifier
	}{
		{
			input: ast.Identifier("foo"),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.input.String(), func(t *testing.T) {
			doc, err := eval.NewDocument(nil)
			if err != nil {
				t.Fatalf("Failed to create test document: %v", err)
			}

			ctx := eval.NewContext(doc, nil, nil)

			_, value, err := eval.EvalIdentifier(ctx, testcase.input)
			if err == nil {
				t.Fatalf("Evaluating identifiers should result in an error, but got %v.", value)
			}
		})
	}
}
