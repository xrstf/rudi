// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"context"
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/types"
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

type cancellerFunc struct {
	cancelFn context.CancelFunc
	called   bool
}

var _ types.Function = &cancellerFunc{}

func (*cancellerFunc) Description() string {
	return ""
}

func (c *cancellerFunc) Evaluate(ctx types.Context, args []ast.Expression) (any, error) {
	if c.cancelFn != nil {
		c.cancelFn()
	}

	c.called = true

	return nil, nil
}

func TestEvalTupleContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// In this testcase, the context has already been cancelled and so this function is a NOP;
	// we need _some_ function to run this test and since we have cancellerFunc already…
	rudiCancelFunc := &cancellerFunc{}

	testcase := testutil.Testcase{
		AST: ast.Tuple{
			Expressions: []ast.Expression{
				ast.Identifier{Name: "cancelctx"},
			},
		},
		Context: ctx,
		Functions: types.Functions{
			"cancelctx": rudiCancelFunc,
		},
		Invalid: true,
	}

	testcase.Run(t)

	if rudiCancelFunc.called {
		t.Fatal("Should not have called cancelctx.")
	}
}

func TestEvalTupleSuddenContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// In this case, we rely on programs working as expected, which enables us to run multiple
	// tuples in succession, of which one will cancel the context, simulating that for example a
	// timeout kicked in.
	anyTuple := &cancellerFunc{}
	cancelsContext := &cancellerFunc{cancelFn: cancel}
	shouldNotBeCalledAnymore := &cancellerFunc{}

	testcase := testutil.Testcase{
		AST: ast.Program{
			Statements: []ast.Statement{
				{
					Expression: ast.Tuple{
						Expressions: []ast.Expression{
							ast.Identifier{Name: "anyTuple"},
						},
					},
				},
				{
					Expression: ast.Tuple{
						Expressions: []ast.Expression{
							ast.Identifier{Name: "cancelsContext"},
						},
					},
				},
				{
					Expression: ast.Tuple{
						Expressions: []ast.Expression{
							ast.Identifier{Name: "shouldNotBeCalledAnymore"},
						},
					},
				},
			},
		},
		Context: ctx,
		Functions: types.Functions{
			"anyTuple":                 anyTuple,
			"cancelsContext":           cancelsContext,
			"shouldNotBeCalledAnymore": shouldNotBeCalledAnymore,
		},
		Invalid: true,
	}

	testcase.Run(t)

	if !anyTuple.called {
		t.Fatal("Should have called anyTuple.")
	}

	if !cancelsContext.called {
		t.Fatal("Should have called cancelsContext.")
	}

	if shouldNotBeCalledAnymore.called {
		t.Fatal("Should not have called shouldNotBeCalledAnymore.")
	}
}
