// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package testutil

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/lang/parser"
	"go.xrstf.de/rudi/pkg/runtime/interpreter"
	"go.xrstf.de/rudi/pkg/runtime/types"

	"github.com/google/go-cmp/cmp"
)

type Testcase struct {
	// use either Expression or AST
	Expression string
	AST        ast.Expression

	Context   context.Context
	Document  any
	Variables types.Variables
	Functions types.Functions
	Coalescer coalescing.Coalescer
	Runtime   types.Runtime

	Expected          any
	ExpectedDocument  any
	ExpectedVariables types.Variables

	Invalid bool
}

func (tc *Testcase) String() string {
	if tc.Expression != "" {
		return tc.Expression
	}

	if tc.AST != nil {
		return tc.AST.String()
	}

	return "<invalid testcase>"
}

func (tc *Testcase) Run(t *testing.T) {
	t.Helper()

	ctx, result, err := tc.eval(t)
	if err != nil {
		if !tc.Invalid {
			t.Fatalf("Failed to eval %s: %v", tc.Expression, err)
		}

		return
	}

	if tc.Invalid {
		t.Fatalf("Should not have been able to eval %s, but got: %v", tc.Expression, result)
	}

	assertResultValue(t, tc.Expected, result)
	assertVariables(t, tc.ExpectedVariables, ctx)
	assertDocument(t, tc.ExpectedDocument, ctx)
}

func (tc *Testcase) eval(t *testing.T) (types.Context, any, error) {
	if (tc.Expression == "") == (tc.AST == nil) {
		t.Fatal("Must use either AST or Expression as test input.")
	}

	doc, err := types.NewDocument(tc.Document)
	if err != nil {
		log.Fatalf("Failed to create parser document: %v", err)
	}

	if tc.Runtime == nil {
		tc.Runtime = interpreter.New()
	}

	progContext, err := types.NewContext(tc.Runtime, tc.Context, doc, tc.Variables, tc.Functions, tc.Coalescer)
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	if tc.Expression != "" {
		prog := strings.NewReader(tc.Expression)

		got, err := parser.ParseReader("test.go", prog)
		if err != nil {
			t.Fatalf("Failed to parse %s: %v", tc.Expression, err)
		}

		program, ok := got.(ast.Program)
		if !ok {
			t.Fatalf("Parsed result is not a ast.Program, but %T", got)
		}

		return tc.Runtime.EvalProgram(progContext, &program)
	}

	// To enable tests for programs and statements, we handle them explicitly
	// instead of extending EvalExpression() to handle them, as that would not
	// fit the language structure.

	switch asserted := tc.AST.(type) {
	case ast.Program:
		return tc.Runtime.EvalProgram(progContext, &asserted)
	case ast.Statement:
		return tc.Runtime.EvalStatement(progContext, asserted)
	default:
		return tc.Runtime.EvalExpression(progContext, tc.AST)
	}
}

func renderDiff(expected any, actual any) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Expected type...: %T\n", expected))
	builder.WriteString(fmt.Sprintf("Expected value..: %#v\n", expected))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("Actual type.....: %T\n", actual))
	builder.WriteString(fmt.Sprintf("Actual value....: %#v\n", actual))

	return builder.String()
}

func assertResultValue(t *testing.T, expected any, actual any) {
	if !cmp.Equal(expected, actual) {
		t.Errorf("Resulting value does not match expectation:\n\n%s\n", renderDiff(expected, actual))
	}
}

func assertDocument(t *testing.T, expected any, ctx types.Context) {
	resultDoc := ctx.GetDocument().Data()

	if !cmp.Equal(expected, resultDoc) {
		t.Errorf("Resulting document does not match expectation:\n\n%s\n", renderDiff(expected, resultDoc))
	}
}

func assertVariables(t *testing.T, expected types.Variables, ctx types.Context) {
	if expected == nil {
		return
	}

	for varName, value := range expected {
		actualValue, ok := ctx.GetVariable(varName)
		if !ok {
			t.Errorf("Variable $%s does not exist anymore.", varName)
			continue
		}

		if !cmp.Equal(value, actualValue) {
			t.Errorf("Variable $%s dooes not match expectation:\n\n%s\n", varName, renderDiff(value, actualValue))
		}
	}
}
