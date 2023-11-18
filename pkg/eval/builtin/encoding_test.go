// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"
)

type encodingTestcase struct {
	expr     string
	expected string
	invalid  bool
}

func (tc *encodingTestcase) Test(t *testing.T) {
	t.Helper()

	result, err := runExpression(t, tc.expr, nil, nil)
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

func TestToBase64Function(t *testing.T) {
	testcases := []encodingTestcase{
		{
			expr:    `(to-base64)`,
			invalid: true,
		},
		{
			expr:    `(to-base64 "too" "many")`,
			invalid: true,
		},
		{
			expr:    `(to-base64 true)`,
			invalid: true,
		},
		{
			expr:    `(to-base64 1)`,
			invalid: true,
		},
		{
			expr:    `(to-base64 null)`,
			invalid: true,
		},
		{
			expr:     `(to-base64 "")`,
			expected: "",
		},
		{
			expr:     `(to-base64 " ")`,
			expected: "IA==",
		},
		{
			expr:     `(to-base64 (concat "" "f" "o" "o"))`,
			expected: "Zm9v",
		},
		{
			expr:     `(to-base64 "test")`,
			expected: "dGVzdA==",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestFromBase64Function(t *testing.T) {
	testcases := []encodingTestcase{
		{
			expr:    `(from-base64)`,
			invalid: true,
		},
		{
			expr:    `(from-base64 "too" "many")`,
			invalid: true,
		},
		{
			expr:    `(from-base64 true)`,
			invalid: true,
		},
		{
			expr:    `(from-base64 1)`,
			invalid: true,
		},
		{
			expr:    `(from-base64 null)`,
			invalid: true,
		},
		{
			expr:    `(from-base64 "definitely-not-base64")`,
			invalid: true,
		},
		{
			// should be able to recover
			expr:     `(try (from-base64 "definitely-not-base64") "fallback")`,
			expected: "fallback",
		},
		{
			expr:     `(from-base64 "")`,
			expected: "",
		},
		{
			expr:     `(from-base64 "IA==")`,
			expected: " ",
		},
		{
			expr:     `(from-base64 (concat "" "Z" "m" "9" "v"))`,
			expected: "foo",
		},
		{
			expr:     `(from-base64 "dGVzdA==")`,
			expected: "test",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}
