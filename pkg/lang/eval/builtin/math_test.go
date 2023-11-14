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

func evalExpression(t *testing.T, expr string, document any) (any, error) {
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
			expr:     "(add 1 2)",
			expected: 3,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, func(t *testing.T) {

		})
	}
}
