// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/equality"
	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func makeProgram(exprs ...ast.Expression) ast.Program {
	prog := ast.Program{
		Statements: []ast.Statement{},
	}

	for _, expr := range exprs {
		prog.Statements = append(prog.Statements, ast.Statement{Expression: expr})
	}

	return prog
}

func makeTuple(exprs ...ast.Expression) ast.Tuple {
	return ast.Tuple{
		Expressions: exprs,
	}
}

func makeVar(name string) ast.Symbol {
	variable := ast.Variable(name)

	return ast.Symbol{
		Variable: &variable,
	}
}

func TestEvalProgram(t *testing.T) {
	testcases := []struct {
		input    ast.Program
		expected ast.Literal
		invalid  bool
	}{
		// (empty program)
		{
			input:    makeProgram(),
			expected: ast.Null{},
		},
		// single statement
		// "foo"
		{
			input: makeProgram(
				ast.String("foo"),
			),
			expected: ast.String("foo"),
		},
		// program result should be the result from the last statement
		// "foo" "bar"
		{
			input: makeProgram(
				ast.String("foo"),
				ast.String("bar"),
			),
			expected: ast.String("bar"),
		},
		// context changes from one statement should affect the next
		// (set $foo 1) $foo (set $bar $foo) $bar
		{
			input: makeProgram(
				makeTuple(
					ast.Identifier{Name: "set"},
					makeVar("foo"),
					ast.Number{Value: 1},
				),
				makeVar("foo"),
				makeTuple(
					ast.Identifier{Name: "set"},
					makeVar("bar"),
					makeVar("foo"),
				),
				makeVar("bar"),
			),
			expected: ast.Number{Value: 1},
		},
		// context changes from inner statements should not leak
		// (set $foo (set $bar 1)) $bar
		{
			input: makeProgram(
				makeTuple(
					ast.Identifier{Name: "set"},
					makeVar("foo"),
					makeTuple(
						ast.Identifier{Name: "set"},
						makeVar("bar"),
						ast.Number{Value: 1},
					),
				),
				makeVar("bar"),
			),
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

			_, value, err := eval.EvalProgram(ctx, &testcase.input)
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
				t.Fatalf("EvalProgram returned unexpected type %T", value)
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
