// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
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

	testcases := []comparisonsTestcase{
		{
			expr:    `(eq)`,
			invalid: true,
		},
		{
			expr:    `(eq true)`,
			invalid: true,
		},
		{
			expr:     `(eq true true)`,
			expected: true,
		},
		{
			expr:     `(eq false true)`,
			expected: false,
		},
		{
			expr:     `(eq .bool true)`,
			expected: true,
			document: testDoc,
		},
		{
			expr:     `(eq 1 1)`,
			expected: true,
		},
		{
			expr:     `(eq 1 2)`,
			expected: false,
		},
		{
			expr:     `(eq .int 4)`,
			expected: true,
			document: testDoc,
		},
		{
			expr:    `(eq 1 "foo")`,
			invalid: true,
		},
		{
			expr:    `(eq 1 true)`,
			invalid: true,
		},
		{
			expr:     `(eq "foo" "foo")`,
			expected: true,
		},
		{
			expr:     `(eq .string "foo")`,
			expected: true,
			document: testDoc,
		},
		{
			expr:     `(eq "foo" "bar")`,
			expected: false,
		},
		{
			expr:     `(eq "foo" "Foo")`,
			expected: false,
		},
		{
			expr:     `(eq "foo" " foo")`,
			expected: false,
		},
		{
			expr:     `(eq [] [])`,
			expected: true,
		},
		{
			expr:     `(eq [1] [])`,
			expected: false,
		},
		{
			expr:     `(eq [] [1])`,
			expected: false,
		},
		{
			expr:     `(eq [1] [1])`,
			expected: true,
		},
		{
			expr:     `(eq [1 [2] {foo "bar"}] [1 [2] {foo "bar"}])`,
			expected: true,
		},
		{
			expr:     `(eq [1 [2] {foo "bar"}] [1 [2] {foo "baz"}])`,
			expected: false,
		},
		{
			expr:     `(eq {} {})`,
			expected: true,
		},
		{
			expr:     `(eq {foo "bar"} {foo "bar"})`,
			expected: true,
		},
		{
			expr:     `(eq {foo "bar"} {foo "baz"})`,
			expected: false,
		},
		{
			expr:     `(eq {foo "bar"} {})`,
			expected: false,
		},
		{
			expr:     `(eq {} {foo "bar"})`,
			expected: false,
		},
		{
			expr:     `(eq {foo "bar" l [1 2]} {foo "bar" l [1 2]})`,
			expected: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}
