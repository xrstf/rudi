// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"
	"testing"
)

type comparisonsTestcase struct {
	expr     string
	expected any
	document any
	invalid  bool
}

func (tc *comparisonsTestcase) Test(t *testing.T) {
	t.Helper()

	result, err := runExpression(t, tc.expr, tc.document, nil)
	if err != nil {
		if !tc.invalid {
			t.Fatalf("Failed to run %s: %v", tc.expr, err)
		}

		return
	}

	if tc.invalid {
		t.Fatalf("Should not have been able to run %s, but got: %v", tc.expr, result)
	}

	if result != tc.expected {
		t.Fatalf("Expected %v (%T), but got %v (%T)", tc.expected, tc.expected, result, result)
	}
}

type flippedTestcases struct {
	left     string
	right    string
	document any
	expected any
	invalid  bool
}

func genFlippedExpressions(fun string, testcases []flippedTestcases) []comparisonsTestcase {
	result := []comparisonsTestcase{}

	for _, tc := range testcases {
		result = append(
			result,
			comparisonsTestcase{
				expr:     fmt.Sprintf(`(%s %s %s)`, fun, tc.left, tc.right),
				invalid:  tc.invalid,
				expected: tc.expected,
				document: tc.document,
			},
			comparisonsTestcase{
				expr:     fmt.Sprintf(`(%s %s %s)`, fun, tc.right, tc.left),
				invalid:  tc.invalid,
				expected: tc.expected,
				document: tc.document,
			},
		)
	}

	return result
}

func TestEqFunction(t *testing.T) {
	testDoc := map[string]any{
		"int":    int(4),
		"float":  float64(1.2),
		"bool":   true,
		"string": "foo",
		"null":   nil,
		"vector": []any{int64(1)},
		"object": map[string]any{
			"key": "value",
		},
	}

	syntax := []comparisonsTestcase{
		{
			expr:    `(eq?)`,
			invalid: true,
		},
		{
			expr:    `(eq? true)`,
			invalid: true,
		},
		{
			expr:    `(eq? "too" "many" "args")`,
			invalid: true,
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
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestLikeFunction(t *testing.T) {
	syntax := []comparisonsTestcase{
		{
			expr:    `(like?)`,
			invalid: true,
		},
		{
			expr:    `(like? true)`,
			invalid: true,
		},
		{
			expr:    `(like? "too" "many" "args")`,
			invalid: true,
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
		t.Run(testcase.expr, testcase.Test)
	}
}
