// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"log"
	"strings"
	"testing"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval"
	"go.xrstf.de/otto/pkg/lang/eval/types"
	"go.xrstf.de/otto/pkg/lang/parser"
)

func runExpression(t *testing.T, expr string, document any) (any, error) {
	prog := strings.NewReader(expr)

	got, err := parser.ParseReader("test.go", prog)
	if err != nil {
		t.Fatalf("Failed to parse %s: %v", expr, err)
	}

	program, ok := got.(ast.Program)
	if !ok {
		t.Fatalf("Parsed result is not a ast.Program, but %T", got)
	}

	doc, err := eval.NewDocument(document)
	if err != nil {
		log.Fatalf("Failed to create parser document: %v", err)
	}

	vars := eval.NewVariables().
		Set("global", types.Must(types.WrapNative(document)))

	progContext := eval.NewContext(doc, Functions, vars)

	return eval.Run(progContext, program)
}

func TestSumFunction(t *testing.T) {
	testcases := []struct {
		expr     string
		expected any
		invalid  bool
	}{
		{
			expr:    `(+)`,
			invalid: true,
		},
		{
			expr:    `(+ 1)`,
			invalid: true,
		},
		{
			expr:    `(+ 1 "foo")`,
			invalid: true,
		},
		{
			expr:    `(+ 1 [])`,
			invalid: true,
		},
		{
			expr:    `(+ 1 {})`,
			invalid: true,
		},
		{
			expr:     `(+ 1 2)`,
			expected: int64(3),
		},
		{
			expr:     `(+ 1 -2 5)`,
			expected: int64(4),
		},
		{
			expr:     `(+ 1 1.5)`,
			expected: float64(2.5),
		},
		{
			expr:     `(+ 1 1.5 (+ 1 2))`,
			expected: float64(5.5),
		},
		{
			expr:     `(+ 0 0.0 -5.6)`,
			expected: float64(-5.6),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, func(t *testing.T) {
			result, err := runExpression(t, testcase.expr, nil)
			if err != nil {
				if !testcase.invalid {
					t.Fatalf("Failed to run %s: %v", testcase.expr, err)
				}

				return
			}

			if testcase.invalid {
				t.Fatalf("Should not have been able to run %s, but got: %v", testcase.expr, result)
			}

			if result != testcase.expected {
				t.Fatalf("Expected %v (%T), but got %v (%T)", testcase.expected, testcase.expected, result, result)
			}
		})
	}
}
