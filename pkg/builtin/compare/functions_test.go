// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package compare

import (
	"fmt"
	"testing"

	"go.xrstf.de/rudi/pkg/eval/types"
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
			// strict coalescing allows lossless float->int conversion
			left:     `1`,
			right:    `1.0`,
			expected: true,
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
			left:    `{}`,
			right:   `[1]`,
			invalid: true,
		},
		{
			left:    `{foo "bar"}`,
			right:   `[]`,
			invalid: true,
		},
	})

	for _, testcase := range append(syntax, testcases...) {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

type comparisonTestcase struct {
	left  any
	right any
	lt    bool
	lte   bool
	gt    bool
	gte   bool
}

func TestInvalidComparisonFunctions(t *testing.T) {
	testcases := []comparisonTestcase{
		{
			left:  3,
			right: []any{},
		},
		{
			left:  []any{},
			right: 3,
		},
	}

	funcs := []func(ctx types.Context, left, right any) (any, error){
		ltFunction,
		lteFunction,
		gtFunction,
		gteFunction,
	}

	ctx := types.NewContext(types.Document{}, nil, nil, nil)

	for _, tc := range testcases {
		for _, f := range funcs {
			_, err := f(ctx, tc.left, tc.right)
			if err == nil {
				t.Errorf("Should have errored on %v <-> %v", tc.left, tc.right)
			}
		}
	}
}

func TestComparisonFunctions(t *testing.T) {
	testcases := []comparisonTestcase{
		{
			left:  0,
			right: 0,
			lt:    false,
			lte:   true,
			gt:    false,
			gte:   true,
		},
		{
			left:  0,
			right: 1,
			lt:    true,
			lte:   true,
			gt:    false,
			gte:   false,
		},
		{
			left:  0,
			right: -1,
			lt:    false,
			lte:   false,
			gt:    true,
			gte:   true,
		},
		{
			left:  -3,
			right: 4.1,
			lt:    true,
			lte:   true,
			gt:    false,
			gte:   false,
		},
		{
			left:  "0",
			right: "-1",
			lt:    false,
			lte:   false,
			gt:    true,
			gte:   true,
		},
		{
			left:  true,
			right: false,
			lt:    false,
			lte:   false,
			gt:    true,
			gte:   true,
		},
		{
			left:  "foo",
			right: "bar",
			lt:    false,
			lte:   false,
			gt:    true,
			gte:   true,
		},
	}

	ctx := types.NewContext(types.Document{}, nil, nil, nil)

	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			lt, err := ltFunction(ctx, tc.left, tc.right)
			if err != nil {
				t.Errorf("lt returned error: %v", err)
			} else if lt != tc.lt {
				t.Errorf("Expected %v < %v, but didn't get that result", tc.left, tc.right)
			}

			lte, err := lteFunction(ctx, tc.left, tc.right)
			if err != nil {
				t.Errorf("lte returned error: %v", err)
			} else if lte != tc.lte {
				t.Errorf("Expected %v <= %v, but didn't get that result", tc.left, tc.right)
			}

			gt, err := gtFunction(ctx, tc.left, tc.right)
			if err != nil {
				t.Errorf("gt returned error: %v", err)
			} else if gt != tc.gt {
				t.Errorf("Expected %v > %v, but didn't get that result", tc.left, tc.right)
			}

			gte, err := gteFunction(ctx, tc.left, tc.right)
			if err != nil {
				t.Errorf("gte returned error: %v", err)
			} else if gte != tc.gte {
				t.Errorf("Expected %v >= %v, but didn't get that result", tc.left, tc.right)
			}
		})
	}
}

func TestComparisonRudiFunctions(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(%s?)`,
		},
		{
			Expression: `(%s? true)`,
		},
		{
			Expression: `(%s? "too" "many" "args")`,
		},
		{
			Expression: `(%s? identifier "foo")`,
		},
		{
			Expression: `(%s? "foo" identifier)`,
		},
		{
			Expression: `(%s? 3 "strings")`,
		},
		{
			Expression: `(%s? 3 [1 2 3])`,
		},
		{
			Expression: `(%s? 3 {foo "bar"})`,
		},
	}

	for _, fun := range []string{"lt", "lte", "gt", "gte"} {
		for _, tc := range testcases {
			test := testutil.Testcase{
				Expression: fmt.Sprintf(tc.Expression, fun),
				Functions:  Functions,
				Invalid:    true,
			}

			t.Run(test.String(), test.Run)
		}
	}
}
