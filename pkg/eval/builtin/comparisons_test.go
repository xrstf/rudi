// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"
	"testing"

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
			expected: true,
		},
		{
			left:     `true`,
			right:    `false`,
			expected: false,
		},
		{
			left:     `true`,
			right:    `.bool`,
			expected: true,
			document: testDoc,
		},
		{
			left:     `1`,
			right:    `1`,
			expected: true,
		},
		{
			left:     `1`,
			right:    `2`,
			expected: false,
		},
		{
			left:     `.int`,
			right:    `4`,
			expected: true,
			document: testDoc,
		},
		{
			left:     `1`,
			right:    `1.0`,
			expected: false,
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
			expected: true,
		},
		{
			left:     `.string`,
			right:    `"foo"`,
			expected: true,
			document: testDoc,
		},
		{
			left:     `"foo"`,
			right:    `"bar"`,
			expected: false,
		},
		{
			left:     `"foo"`,
			right:    `"Foo"`,
			expected: false,
		},
		{
			left:     `"foo"`,
			right:    `" foo"`,
			expected: false,
		},
		{
			left:     `[]`,
			right:    `[]`,
			expected: true,
		},
		{
			left:     `[]`,
			right:    `[1]`,
			expected: false,
		},
		{
			left:     `[1]`,
			right:    `[1]`,
			expected: true,
		},
		{
			left:     `[1 [2] {foo "bar"}]`,
			right:    `[1 [2] {foo "bar"}]`,
			expected: true,
		},
		{
			left:     `[1 [2] {foo "bar"}]`,
			right:    `[1 [2] {foo "baz"}]`,
			expected: false,
		},
		{
			left:     `{}`,
			right:    `{}`,
			expected: true,
		},
		{
			left:     `{}`,
			right:    `{foo "bar"}`,
			expected: false,
		},
		{
			left:     `{foo "bar"}`,
			right:    `{foo "bar"}`,
			expected: true,
		},
		{
			left:     `{foo "bar"}`,
			right:    `{foo "baz"}`,
			expected: false,
		},
		{
			left:     `{foo "bar" l [1 2]}`,
			right:    `{foo "bar" l [1 2]}`,
			expected: true,
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
			expected: true,
		},
		{
			left:     `1`,
			right:    `"1"`,
			expected: true,
		},
		{
			left:     `1`,
			right:    `1.0`,
			expected: true,
		},
		{
			left:     `1`,
			right:    `"1.0"`,
			expected: true,
		},
		{
			left:     `1`,
			right:    `"2.0"`,
			expected: false,
		},
		{
			left:     `1`,
			right:    `true`,
			expected: true,
		},
		{
			left:     `1`,
			right:    `2`,
			expected: false,
		},
		{
			left:     `false`,
			right:    `null`,
			expected: true,
		},
		{
			left:     `false`,
			right:    `"null"`,
			expected: false,
		},
		{
			left:     `0`,
			right:    `null`,
			expected: true,
		},
		{
			left:     `0.0`,
			right:    `null`,
			expected: true,
		},
		{
			left:     `""`,
			right:    `null`,
			expected: true,
		},
		{
			left:     `false`,
			right:    `0`,
			expected: true,
		},
		{
			left:     `false`,
			right:    `""`,
			expected: true,
		},
		{
			left:     `false`,
			right:    `[]`,
			expected: true,
		},
		{
			left:     `false`,
			right:    `{}`,
			expected: true,
		},
		{
			left:     `false`,
			right:    `"false"`,
			expected: true,
		},
		{
			left:     `true`,
			right:    `{foo "bar"}`,
			expected: true,
		},
		{
			left:     `"foo"`,
			right:    `"bar"`,
			expected: false,
		},
		{
			left:     `true`,
			right:    `[""]`,
			expected: true,
		},
		{
			left:     `{}`,
			right:    `[]`,
			expected: true,
		},
		{
			left:     `{}`,
			right:    `[1]`,
			expected: false,
		},
		{
			left:     `{foo "bar"}`,
			right:    `[]`,
			expected: false,
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
			Expected:   false,
		},
		{
			Expression: `(lt? 2 (+ 1 2))`,
			Expected:   true,
		},
		{
			Expression: `(lt? 2 3)`,
			Expected:   true,
		},
		{
			Expression: `(lt? -3 2)`,
			Expected:   true,
		},
		{
			Expression: `(lt? -3 -5)`,
			Expected:   false,
		},
		{
			Expression: `(lt? 3.4 3.4)`,
			Expected:   false,
		},
		{
			Expression: `(lt? 2.4 (+ 1.4 2))`,
			Expected:   true,
		},
		{
			Expression: `(lt? 2.4 3.4)`,
			Expected:   true,
		},
		{
			Expression: `(lt? -3.4 2.4)`,
			Expected:   true,
		},
		{
			Expression: `(lt? -3.4 -5.4)`,
			Expected:   false,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}
