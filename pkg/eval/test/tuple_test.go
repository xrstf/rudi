// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/testutil"
)

func TestEvalTuple(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			AST:     ast.Tuple{},
			Invalid: true,
		},
		// (true)
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Bool(true),
				},
			},
			Invalid: true,
		},
		// ("invalid")
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.String("invalid"),
				},
			},
			Invalid: true,
		},
		// (1)
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Number{Value: 1},
				},
			},
			Invalid: true,
		},
		// ((eval))
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Tuple{Expressions: []ast.Expression{ast.Identifier{Name: "eval"}}},
				},
			},
			Invalid: true,
		},
		// (unknownfunc)
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "unknownfunc"},
				},
			},
			Invalid: true,
		},
		// (eval "too" "many")
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "eval"},
					ast.String("too"),
					ast.String("many"),
				},
			},
			Invalid: true,
		},
		// (eval "foo")
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "eval"},
					ast.String("foo"),
				},
			},
			Expected: "foo",
		},
		// (eval {foo "bar"}).foo
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "eval"},
					ast.ObjectNode{
						Data: []ast.KeyValuePair{
							{
								Key:   ast.Identifier{Name: "foo"},
								Value: ast.String("bar"),
							},
						},
					},
				},
				PathExpression: &ast.PathExpression{
					Steps: []ast.Expression{
						ast.Identifier{Name: "foo"},
					},
				},
			},
			Expected: "bar",
		},
		// (eval {foo "bar"})[1]
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "eval"},
					ast.ObjectNode{
						Data: []ast.KeyValuePair{
							{
								Key:   ast.Identifier{Name: "foo"},
								Value: ast.String("bar"),
							},
						},
					},
				},
				PathExpression: &ast.PathExpression{
					Steps: []ast.Expression{
						ast.Number{Value: 1},
					},
				},
			},
			Invalid: true,
		},
		// (eval [1 2])[1]
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "eval"},
					ast.VectorNode{
						Expressions: []ast.Expression{
							ast.Number{Value: 1},
							ast.Number{Value: 2},
						},
					},
				},
				PathExpression: &ast.PathExpression{
					Steps: []ast.Expression{
						ast.Number{Value: 1},
					},
				},
			},
			Expected: 2,
		},
		// (eval [1 2]).invalid
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "eval"},
					ast.VectorNode{
						Expressions: []ast.Expression{
							ast.Number{Value: 1},
							ast.Number{Value: 2},
						},
					},
				},
				PathExpression: &ast.PathExpression{
					Steps: []ast.Expression{
						ast.String("invalid"),
					},
				},
			},
			Invalid: true,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = dummyFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestEvalTupleBangModifier(t *testing.T) {
	testcases := []testutil.Testcase{
		// (set!)
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
				},
			},
			Invalid: true,
		},
		// (set! "invalid" "value")
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					ast.String("invalid"),
					ast.String("value"),
				},
			},
			Invalid: true,
		},
		// (set! {} "value")
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					ast.ObjectNode{},
					ast.String("value"),
				},
			},
			Invalid: true,
		},
		// (set! .[true] "value")
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					ast.Symbol{
						PathExpression: &ast.PathExpression{
							Steps: []ast.Expression{
								ast.Bool(true),
							},
						},
					},
					ast.String("value"),
				},
			},
			Invalid: true,
		},
		// (set! .[-1] "value")
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					ast.Symbol{
						PathExpression: &ast.PathExpression{
							Steps: []ast.Expression{
								ast.Number{Value: -1},
							},
						},
					},
					ast.String("value"),
				},
			},
			Invalid: true,
		},
		// (set! . "value")
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					ast.Symbol{PathExpression: &ast.PathExpression{}},
					ast.String("value"),
				},
			},
			Expected:         "value",
			ExpectedDocument: "value",
		},
		// (set! . "value")
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					ast.Symbol{PathExpression: &ast.PathExpression{}},
					ast.String("value"),
				},
			},
			Expected:         "value",
			ExpectedDocument: "value",
		},
		// (set! $val "value")
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					makeVar("myvar", nil),
					ast.String("value"),
				},
			},
			Expected: "value",
			ExpectedVariables: types.Variables{
				"myvar": "value",
			},
		},
		// (set! .hello "value")
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					ast.Symbol{PathExpression: &ast.PathExpression{
						Steps: []ast.Expression{
							ast.Identifier{Name: "hello"},
						},
					}},
					ast.String("value"),
				},
			},
			Document:         map[string]any{"hello": "world", "hei": "verden"},
			Expected:         "value",
			ExpectedDocument: map[string]any{"hello": "value", "hei": "verden"},
		},
		// (set! $val.key "value")
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					makeVar("val", &ast.PathExpression{
						Steps: []ast.Expression{
							ast.Identifier{Name: "key"},
						},
					}),
					ast.String("value"),
				},
			},
			Variables: types.Variables{
				"val": map[string]any{"foo": "bar", "key": 42},
			},
			Expected: "value",
			ExpectedVariables: types.Variables{
				"val": map[string]any{"foo": "bar", "key": "value"},
			},
		},
		// (set! .hei (set! .hello "value"))
		{
			AST: ast.Tuple{
				Expressions: []ast.Expression{
					ast.Identifier{Name: "set", Bang: true},
					ast.Symbol{PathExpression: &ast.PathExpression{
						Steps: []ast.Expression{
							ast.Identifier{Name: "hei"},
						},
					}},
					ast.Tuple{
						Expressions: []ast.Expression{
							ast.Identifier{Name: "set", Bang: true},
							ast.Symbol{PathExpression: &ast.PathExpression{
								Steps: []ast.Expression{
									ast.Identifier{Name: "hello"},
								},
							}},
							ast.String("value"),
						},
					},
				},
			},
			Document:         map[string]any{"hello": "world", "hei": "verden"},
			Expected:         "value",
			ExpectedDocument: map[string]any{"hello": "value", "hei": "value"},
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = dummyFunctions
		t.Run(testcase.String(), testcase.Run)
	}
}
