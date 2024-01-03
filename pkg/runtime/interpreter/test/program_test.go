// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/testutil"
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

func makeVar(name string, pathExpr *ast.PathExpression) ast.Symbol {
	variable := ast.Variable(name)

	return ast.Symbol{
		Variable:       &variable,
		PathExpression: pathExpr,
	}
}

func TestEvalProgram(t *testing.T) {
	testcases := []testutil.Testcase{
		// (empty program)
		{
			AST:      makeProgram(),
			Expected: nil,
		},
		// single statement
		// "foo"
		{
			AST: makeProgram(
				ast.String("foo"),
			),
			Expected: "foo",
		},
		// program result should be the result from the last statement
		// "foo" "bar"
		{
			AST: makeProgram(
				ast.String("foo"),
				ast.String("bar"),
			),
			Expected: "bar",
		},
		// context changes from one statement should affect the next
		// (set! $foo 1) $foo (set! $bar $foo) $bar
		{
			AST: makeProgram(
				makeTuple(
					ast.Identifier{Name: "set", Bang: true},
					makeVar("foo", nil),
					ast.Number{Value: 1},
				),
				makeVar("foo", nil),
				makeTuple(
					ast.Identifier{Name: "set", Bang: true},
					makeVar("bar", nil),
					makeVar("foo", nil),
				),
				makeVar("bar", nil),
			),
			Expected: 1,
		},
		// all variables share one scope, even in sub expressions
		// (set! $foo (set! $bar 1)) $bar
		{
			AST: makeProgram(
				makeTuple(
					ast.Identifier{Name: "set", Bang: true},
					makeVar("foo", nil),
					makeTuple(
						ast.Identifier{Name: "set", Bang: true},
						makeVar("bar", nil),
						ast.Number{Value: 1},
					),
				),
				makeVar("bar", nil),
			),
			Expected: 1,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = dummyFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}
