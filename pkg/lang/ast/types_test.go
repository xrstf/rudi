// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package ast

import (
	"testing"
)

func makeVar(name string) *Variable {
	v := Variable(name)
	return &v
}

func ptrTo[T any](v T) *T {
	return &v
}

func TestExpressionNames(t *testing.T) {
	testcases := []struct {
		expr     Expression
		expected string
	}{
		{
			// technically invalid
			expr:     Program{},
			expected: "Program(?)",
		},
		{
			expr: Program{
				Statements: []Statement{{
					Expression: Null{},
				}},
			},
			expected: "Program(Statements)",
		},
		{
			expr:     Statement{},
			expected: "Statement",
		},
		{
			expr:     Null{},
			expected: "Null",
		},
		{
			expr:     Bool(true),
			expected: "Bool",
		},
		{
			expr:     String("foo"),
			expected: "String",
		},
		{
			expr:     Number{},
			expected: "Number",
		},
		{
			expr:     Symbol{},
			expected: "Symbol(?)",
		},
		{
			expr:     Vector{},
			expected: "Vector",
		},
		{
			expr:     VectorNode{}, // ...Node suffix is only an internal distinction
			expected: "Vector",
		},
		{
			expr:     Object{},
			expected: "Object",
		},
		{
			expr:     ObjectNode{},
			expected: "Object",
		},
		{
			expr:     KeyValuePair{},
			expected: "KeyValuePair",
		},
		{
			expr:     Variable("foo"),
			expected: "Variable",
		},
		{
			expr:     Identifier("foo"),
			expected: "Identifier",
		},
		{
			expr: Symbol{
				Variable: makeVar("foo"),
			},
			expected: "Symbol(Variable)",
		},
		{
			expr: Symbol{
				Variable:       makeVar("foo"),
				PathExpression: &PathExpression{Steps: []Expression{String("foo")}},
			},
			expected: "Symbol(Variable)",
		},
		{
			expr: Symbol{
				PathExpression: &PathExpression{Steps: []Expression{String("foo")}},
			},
			expected: "Symbol(PathExpression)",
		},
		{
			expr:     Tuple{},
			expected: "Tuple",
		},
		{
			expr:     PathExpression{},
			expected: "PathExpression",
		},
		{
			expr:     EvaluatedPathExpression{},
			expected: "PathExpression",
		},
		{
			expr:     EvaluatedPathStep{StringValue: ptrTo("foo")},
			expected: "PathStep(String)",
		},
		{
			expr:     EvaluatedPathStep{IntegerValue: ptrTo(int64(42))},
			expected: "PathStep(Number)",
		},
	}

	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			result := tc.expr.ExpressionName()

			if result != tc.expected {
				t.Fatalf("Expected %q, got %q.", tc.expected, result)
			}
		})
	}
}

func TestStringers(t *testing.T) {
	testcases := []struct {
		expr     Expression
		expected string
	}{
		{
			// technically invalid
			expr:     Program{},
			expected: "",
		},
		{
			// single statement
			expr: Program{
				Statements: []Statement{{
					Expression: Bool(true),
				}},
			},
			expected: `true`,
		},
		{
			// multiple statements are separated by a space
			expr: Program{
				Statements: []Statement{{
					Expression: Bool(true),
				}, {
					Expression: String("foo"),
				}},
			},
			expected: `true "foo"`,
		},
		{
			expr:     Statement{},
			expected: `<invalid Statement>`,
		},
		{
			expr:     Statement{Expression: Bool(true)},
			expected: `true`,
		},
		{
			expr:     Symbol{},
			expected: `<invalid Symbol>`,
		},
		{
			expr: Symbol{
				Variable: makeVar("foo"),
			},
			expected: `$foo`,
		},
		{
			expr: Symbol{
				Variable: makeVar("foo"),
				PathExpression: &PathExpression{
					Steps: []Expression{
						Number{Value: 1},
					},
				},
			},
			expected: `$foo[1]`,
		},
		{
			expr: Symbol{
				Variable: makeVar("foo"),
				PathExpression: &PathExpression{
					Steps: []Expression{
						Number{Value: 1},
						Number{Value: 2},
					},
				},
			},
			expected: `$foo[1][2]`,
		},
		{
			expr: Symbol{
				Variable: makeVar("foo"),
				PathExpression: &PathExpression{
					Steps: []Expression{
						Number{Value: 1},
						Identifier("foo"),
					},
				},
			},
			expected: `$foo[1].foo`,
		},
		{
			expr: Symbol{
				Variable: makeVar("foo"),
				PathExpression: &PathExpression{
					Steps: []Expression{
						Identifier("foo"),
					},
				},
			},
			expected: `$foo.foo`,
		},
		{
			expr: Symbol{
				Variable: makeVar("foo"),
				PathExpression: &PathExpression{
					Steps: []Expression{
						Identifier("foo"),
						Identifier("bar"),
					},
				},
			},
			expected: `$foo.foo.bar`,
		},
		{
			expr: Symbol{
				Variable: makeVar("foo"),
				PathExpression: &PathExpression{
					Steps: []Expression{
						Identifier("foo"),
						Number{Value: 1},
					},
				},
			},
			expected: `$foo.foo[1]`,
		},
		{
			expr: Symbol{
				Variable: makeVar("foo"),
				PathExpression: &PathExpression{
					Steps: []Expression{
						String("foo"),
					},
				},
			},
			expected: `$foo["foo"]`,
		},
		{
			expr: Symbol{
				Variable: makeVar("foo"),
				PathExpression: &PathExpression{
					Steps: []Expression{
						Symbol{
							Variable: makeVar("bla"),
							PathExpression: &PathExpression{
								Steps: []Expression{
									String("sub"),
								},
							},
						},
					},
				},
			},
			expected: `$foo[$bla["sub"]]`,
		},
		{
			expr: Symbol{
				Variable: makeVar("foo"),
				PathExpression: &PathExpression{
					Steps: []Expression{
						Bool(true),
					},
				},
			},
			expected: `$foo[true]`,
		},
		{
			expr: Symbol{
				Variable: makeVar("foo"),
				PathExpression: &PathExpression{
					Steps: []Expression{
						Null{},
					},
				},
			},
			expected: `$foo[null]`,
		},
		{
			expr: Symbol{
				Variable: makeVar("foo"),
				PathExpression: &PathExpression{
					Steps: []Expression{
						Tuple{
							Expressions: []Expression{
								Identifier("ident"),
							},
						},
					},
				},
			},
			expected: `$foo[(ident)]`,
		},
		{
			expr: Symbol{
				Variable: makeVar("foo"),
				PathExpression: &PathExpression{
					Steps: []Expression{
						ObjectNode{
							Data: []KeyValuePair{
								{
									Key:   Identifier("k"),
									Value: String("v"),
								},
							},
						},
					},
				},
			},
			// technically invalid without a path expr on that object, but still
			// printable
			expected: `$foo[{k "v"}]`,
		},
		{
			expr: Symbol{
				Variable: makeVar("foo"),
				PathExpression: &PathExpression{
					Steps: []Expression{
						VectorNode{
							Expressions: []Expression{
								String("foo"),
							},
						},
					},
				},
			},
			// technically invalid without a path expr on that vector, but still
			// printable
			expected: `$foo[["foo"]]`,
		},
	}

	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			result := tc.expr.String()

			if result != tc.expected {
				t.Fatalf("Expected %q, got %q.", tc.expected, result)
			}
		})
	}
}
