// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/testutil"
)

type flippedTestcases struct {
	left     string
	right    string
	document any
	expected any
	invalid  bool
}

func genFlippedExpressions(fun string, testcases []flippedTestcases) []testutil.Testcase {
	result := []testutil.Testcase{}

	for _, tc := range testcases {
		result = append(
			result,
			testutil.Testcase{
				Expression:       fmt.Sprintf(`(%s %s %s)`, fun, tc.left, tc.right),
				Invalid:          tc.invalid,
				Expected:         tc.expected,
				Document:         tc.document,
				ExpectedDocument: tc.document,
			},
			testutil.Testcase{
				Expression:       fmt.Sprintf(`(%s %s %s)`, fun, tc.right, tc.left),
				Invalid:          tc.invalid,
				Expected:         tc.expected,
				Document:         tc.document,
				ExpectedDocument: tc.document,
			},
		)
	}

	return result
}

func TestEqFunction(t *testing.T) {
	testDoc := map[string]any{
		"int":    int64(4),
		"float":  float64(1.2),
		"bool":   true,
		"string": "foo",
		"null":   nil,
		"vector": []any{int64(1)},
		"object": map[string]any{
			"key": "value",
		},
	}

	syntax := []testutil.Testcase{
		{
			Expression: `(eq?)`,
			Invalid:    true,
		},
		{
			Expression: `(eq? true)`,
			Invalid:    true,
		},
		{
			Expression: `(eq? "too" "many" "args")`,
			Invalid:    true,
		},
		{
			Expression: `(eq? identifier "foo")`,
			Invalid:    true,
		},
		{
			Expression: `(eq? "foo" identifier)`,
			Invalid:    true,
		},
	}

	flipped := genFlippedExpressions("eq?", []flippedTestcases{
		{
			left:     `true`,
			right:    `true`,
			expected: ast.Bool(true),
		},
		{
			left:     `true`,
			right:    `false`,
			expected: ast.Bool(false),
		},
		{
			left:     `true`,
			right:    `.bool`,
			expected: ast.Bool(true),
			document: testDoc,
		},
		{
			left:     `1`,
			right:    `1`,
			expected: ast.Bool(true),
		},
		{
			left:     `1`,
			right:    `2`,
			expected: ast.Bool(false),
		},
		{
			left:     `.int`,
			right:    `4`,
			expected: ast.Bool(true),
			document: testDoc,
		},
		{
			left:     `1`,
			right:    `1.0`,
			expected: ast.Bool(false),
		},
		{
			left:    `1`,
			right:   `"foo"`,
			invalid: true,
		},
		{
			left:    `1`,
			right:   `true`,
			invalid: true,
		},
		{
			left:     `"foo"`,
			right:    `"foo"`,
			expected: ast.Bool(true),
		},
		{
			left:     `.string`,
			right:    `"foo"`,
			expected: ast.Bool(true),
			document: testDoc,
		},
		{
			left:     `"foo"`,
			right:    `"bar"`,
			expected: ast.Bool(false),
		},
		{
			left:     `"foo"`,
			right:    `"Foo"`,
			expected: ast.Bool(false),
		},
		{
			left:     `"foo"`,
			right:    `" foo"`,
			expected: ast.Bool(false),
		},
		{
			left:     `[]`,
			right:    `[]`,
			expected: ast.Bool(true),
		},
		{
			left:     `[]`,
			right:    `[1]`,
			expected: ast.Bool(false),
		},
		{
			left:     `[1]`,
			right:    `[1]`,
			expected: ast.Bool(true),
		},
		{
			left:     `[1 [2] {foo "bar"}]`,
			right:    `[1 [2] {foo "bar"}]`,
			expected: ast.Bool(true),
		},
		{
			left:     `[1 [2] {foo "bar"}]`,
			right:    `[1 [2] {foo "baz"}]`,
			expected: ast.Bool(false),
		},
		{
			left:     `{}`,
			right:    `{}`,
			expected: ast.Bool(true),
		},
		{
			left:     `{}`,
			right:    `{foo "bar"}`,
			expected: ast.Bool(false),
		},
		{
			left:     `{foo "bar"}`,
			right:    `{foo "bar"}`,
			expected: ast.Bool(true),
		},
		{
			left:     `{foo "bar"}`,
			right:    `{foo "baz"}`,
			expected: ast.Bool(false),
		},
		{
			left:     `{foo "bar" l [1 2]}`,
			right:    `{foo "bar" l [1 2]}`,
			expected: ast.Bool(true),
		},
	})

	for _, testcase := range append(syntax, flipped...) {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestLikeFunction(t *testing.T) {
	syntax := []testutil.Testcase{
		{
			Expression: `(like?)`,
			Invalid:    true,
		},
		{
			Expression: `(like? true)`,
			Invalid:    true,
		},
		{
			Expression: `(like? "too" "many" "args")`,
			Invalid:    true,
		},
		{
			Expression: `(like? identifier "foo")`,
			Invalid:    true,
		},
		{
			Expression: `(like? "foo" identifier)`,
			Invalid:    true,
		},
	}

	testcases := genFlippedExpressions("like?", []flippedTestcases{
		{
			left:     `1`,
			right:    `1`,
			expected: ast.Bool(true),
		},
		{
			left:     `1`,
			right:    `"1"`,
			expected: ast.Bool(true),
		},
		{
			left:     `1`,
			right:    `1.0`,
			expected: ast.Bool(true),
		},
		{
			left:     `1`,
			right:    `"1.0"`,
			expected: ast.Bool(true),
		},
		{
			left:     `1`,
			right:    `"2.0"`,
			expected: ast.Bool(false),
		},
		{
			left:     `1`,
			right:    `true`,
			expected: ast.Bool(true),
		},
		{
			left:     `1`,
			right:    `2`,
			expected: ast.Bool(false),
		},
		{
			left:     `false`,
			right:    `null`,
			expected: ast.Bool(true),
		},
		{
			left:     `false`,
			right:    `"null"`,
			expected: ast.Bool(false),
		},
		{
			left:     `0`,
			right:    `null`,
			expected: ast.Bool(true),
		},
		{
			left:     `0.0`,
			right:    `null`,
			expected: ast.Bool(true),
		},
		{
			left:     `""`,
			right:    `null`,
			expected: ast.Bool(true),
		},
		{
			left:     `false`,
			right:    `0`,
			expected: ast.Bool(true),
		},
		{
			left:     `false`,
			right:    `""`,
			expected: ast.Bool(true),
		},
		{
			left:     `false`,
			right:    `[]`,
			expected: ast.Bool(true),
		},
		{
			left:     `false`,
			right:    `{}`,
			expected: ast.Bool(true),
		},
		{
			left:     `false`,
			right:    `"false"`,
			expected: ast.Bool(true),
		},
		{
			left:     `true`,
			right:    `{foo "bar"}`,
			expected: ast.Bool(true),
		},
		{
			left:     `"foo"`,
			right:    `"bar"`,
			expected: ast.Bool(false),
		},
		{
			left:     `true`,
			right:    `[""]`,
			expected: ast.Bool(true),
		},
		{
			left:     `{}`,
			right:    `[]`,
			expected: ast.Bool(true),
		},
		{
			left:     `{}`,
			right:    `[1]`,
			expected: ast.Bool(false),
		},
		{
			left:     `{foo "bar"}`,
			right:    `[]`,
			expected: ast.Bool(false),
		},
	})

	for _, testcase := range append(syntax, testcases...) {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestLtFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(lt?)`,
			Invalid:    true,
		},
		{
			Expression: `(lt? true)`,
			Invalid:    true,
		},
		{
			Expression: `(lt? "too" "many" "args")`,
			Invalid:    true,
		},
		{
			Expression: `(lt? identifier "foo")`,
			Invalid:    true,
		},
		{
			Expression: `(lt? "foo" identifier)`,
			Invalid:    true,
		},
		{
			Expression: `(lt? 3 "strings")`,
			Invalid:    true,
		},
		{
			Expression: `(lt? 3 3.1)`,
			Invalid:    true,
		},
		{
			Expression: `(lt? 3 [1 2 3])`,
			Invalid:    true,
		},
		{
			Expression: `(lt? 3 {foo "bar"})`,
			Invalid:    true,
		},
		{
			Expression: `(lt? 3 3)`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(lt? 2 (+ 1 2))`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(lt? 2 3)`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(lt? -3 2)`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(lt? -3 -5)`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(lt? 3.4 3.4)`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(lt? 2.4 (+ 1.4 2))`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(lt? 2.4 3.4)`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(lt? -3.4 2.4)`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(lt? -3.4 -5.4)`,
			Expected:   ast.Bool(false),
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}
