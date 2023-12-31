// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/testutil"
)

func TestEvalObjectNode(t *testing.T) {
	testcases := []testutil.Testcase{
		// {}
		{
			AST:      ast.ObjectNode{},
			Expected: map[string]any{},
		},
		// {foo "bar"}
		{
			AST: ast.ObjectNode{
				Data: []ast.KeyValuePair{
					{
						Key:   ast.Identifier{Name: "foo"},
						Value: ast.String("bar"),
					},
				},
			},
			Expected: map[string]any{
				"foo": "bar",
			},
		},
		// {null "bar"}
		{
			AST: ast.ObjectNode{
				Data: []ast.KeyValuePair{
					{
						Key:   ast.Null{},
						Value: ast.String("bar"),
					},
				},
			},
			Expected: map[string]any{
				"": "bar",
			},
		},
		// {(eval "evaled") (eval "also evaled")}
		{
			AST: ast.ObjectNode{
				Data: []ast.KeyValuePair{
					{
						Key: ast.Tuple{
							Expressions: []ast.Expression{
								ast.Identifier{Name: "eval"},
								ast.String("evaled"),
							},
						},
						Value: ast.Tuple{
							Expressions: []ast.Expression{
								ast.Identifier{Name: "eval"},
								ast.String("also evaled"),
							},
						},
					},
				},
			},
			Expected: map[string]any{
				"evaled": "also evaled",
			},
		},
		// {{foo "bar"} "test"}
		{
			AST: ast.ObjectNode{
				Data: []ast.KeyValuePair{
					{
						Key: ast.ObjectNode{
							Data: []ast.KeyValuePair{
								{
									Key:   ast.Identifier{Name: "foo"},
									Value: ast.String("bar"),
								},
							},
						},
						Value: ast.String("test"),
					},
				},
			},
			Invalid: true,
		},
		// {{foo "bar"}.foo "test"}
		{
			AST: ast.ObjectNode{
				Data: []ast.KeyValuePair{
					{
						Key: ast.ObjectNode{
							Data: []ast.KeyValuePair{
								{
									Key:   ast.Identifier{Name: "foo"},
									Value: ast.String("bar"),
								},
							},
							PathExpression: &ast.PathExpression{
								Steps: []ast.Expression{
									ast.Identifier{Name: "foo"},
								},
							},
						},
						Value: ast.String("test"),
					},
				},
			},
			Expected: map[string]any{
				"bar": "test",
			},
		},
		// {foo "bar"}.foo
		{
			AST: ast.ObjectNode{
				Data: []ast.KeyValuePair{
					{
						Key:   ast.Identifier{Name: "foo"},
						Value: ast.String("bar"),
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
		// {foo bar}
		{
			AST: ast.ObjectNode{
				Data: []ast.KeyValuePair{
					{
						Key:   ast.Identifier{Name: "foo"},
						Value: ast.Identifier{Name: "bar"},
					},
				},
			},
			Invalid: true,
		},
		// {true "bar"}
		{
			AST: ast.ObjectNode{
				Data: []ast.KeyValuePair{
					{
						Key:   ast.Bool(true),
						Value: ast.String("bar"),
					},
				},
			},
			Invalid: true,
		},
		// {1 "bar"}
		{
			AST: ast.ObjectNode{
				Data: []ast.KeyValuePair{
					{
						Key:   ast.Number{Value: 1},
						Value: ast.String("bar"),
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
