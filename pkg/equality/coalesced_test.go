// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package equality

import (
	"fmt"
	"testing"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

type invalidConversion int

const invalid invalidConversion = iota

type coalescedTestcase struct {
	left     any
	right    any
	pedantic any
	strict   any
	humane   any
}

func newCoalescedTest(left, right any, pedantic, strict, humane any) coalescedTestcase {
	return coalescedTestcase{
		left:     left,
		right:    right,
		pedantic: pedantic,
		strict:   strict,
		humane:   humane,
	}
}

// type checklist:
// null, bool, int64, float64, string, vector, object
// for brevity's sake, we know that int==int32==int64 internally, likewise for floats

func getEqualCoalescedTestcases() []coalescedTestcase {
	return []coalescedTestcase{
		///////////////////////////////////////////////////////////
		// test nil against all other types

		newCoalescedTest(nil, nil, true, true, true),
		newCoalescedTest(nil, ast.Null{}, true, true, true),

		newCoalescedTest(nil, true, invalid, invalid, invalid),
		newCoalescedTest(nil, false, invalid, invalid, true),
		newCoalescedTest(nil, ast.Bool(true), invalid, invalid, invalid),
		newCoalescedTest(nil, ast.Bool(false), invalid, invalid, true),

		newCoalescedTest(nil, int64(0), invalid, invalid, true),
		newCoalescedTest(nil, float64(0), invalid, invalid, true),
		newCoalescedTest(nil, float64(0.0), invalid, invalid, true),
		newCoalescedTest(nil, float64(0.1), invalid, invalid, invalid),
		newCoalescedTest(nil, int64(1), invalid, invalid, invalid),
		newCoalescedTest(nil, float64(1), invalid, invalid, invalid),
		newCoalescedTest(nil, int64(-1), invalid, invalid, invalid),
		newCoalescedTest(nil, float64(-1), invalid, invalid, invalid),
		newCoalescedTest(nil, ast.Number{Value: int64(0)}, invalid, invalid, true),
		newCoalescedTest(nil, ast.Number{Value: float64(0)}, invalid, invalid, true),
		newCoalescedTest(nil, ast.Number{Value: float64(0.0)}, invalid, invalid, true),
		newCoalescedTest(nil, ast.Number{Value: float64(0.1)}, invalid, invalid, invalid),
		newCoalescedTest(nil, ast.Number{Value: int64(1)}, invalid, invalid, invalid),
		newCoalescedTest(nil, ast.Number{Value: float64(1)}, invalid, invalid, invalid),
		newCoalescedTest(nil, ast.Number{Value: int64(-1)}, invalid, invalid, invalid),
		newCoalescedTest(nil, ast.Number{Value: float64(-1)}, invalid, invalid, invalid),

		newCoalescedTest(nil, "", invalid, invalid, true),
		newCoalescedTest(nil, " ", invalid, invalid, invalid),
		newCoalescedTest(nil, "test", invalid, invalid, invalid),
		newCoalescedTest(nil, ast.String(""), invalid, invalid, true),
		newCoalescedTest(nil, ast.String(" "), invalid, invalid, invalid),
		newCoalescedTest(nil, ast.String("test"), invalid, invalid, invalid),

		newCoalescedTest(nil, []any{}, invalid, invalid, true),
		newCoalescedTest(nil, []any{0}, invalid, invalid, invalid),
		newCoalescedTest(nil, []any{1}, invalid, invalid, invalid),
		newCoalescedTest(nil, []any{""}, invalid, invalid, invalid),
		newCoalescedTest(nil, ast.Vector{Data: nil}, invalid, invalid, true),
		newCoalescedTest(nil, ast.Vector{Data: []any{}}, invalid, invalid, true),
		newCoalescedTest(nil, ast.Vector{Data: []any{0}}, invalid, invalid, invalid),
		newCoalescedTest(nil, ast.Vector{Data: []any{1}}, invalid, invalid, invalid),
		newCoalescedTest(nil, ast.Vector{Data: []any{""}}, invalid, invalid, invalid),

		newCoalescedTest(nil, map[string]any{}, invalid, invalid, true),
		newCoalescedTest(nil, map[string]any{"": ""}, invalid, invalid, invalid),
		newCoalescedTest(nil, ast.Object{Data: nil}, invalid, invalid, true),
		newCoalescedTest(nil, ast.Object{Data: map[string]any{}}, invalid, invalid, true),
		newCoalescedTest(nil, ast.Object{Data: map[string]any{"": ""}}, invalid, invalid, invalid),

		///////////////////////////////////////////////////////////
		// test bool against all other types, except nils

		newCoalescedTest(true, true, true, true, true),
		newCoalescedTest(false, false, true, true, true),
		newCoalescedTest(true, false, false, false, false),
		newCoalescedTest(true, ast.Bool(true), true, true, true),
		newCoalescedTest(false, ast.Bool(false), true, true, true),
		newCoalescedTest(true, ast.Bool(false), false, false, false),

		newCoalescedTest(true, int64(0), invalid, invalid, false),
		newCoalescedTest(true, float64(0), invalid, invalid, false),
		newCoalescedTest(true, int64(1), invalid, invalid, true),
		newCoalescedTest(true, float64(1), invalid, invalid, true),
		newCoalescedTest(true, float64(1.0), invalid, invalid, true),
		newCoalescedTest(true, int64(-1), invalid, invalid, true),
		newCoalescedTest(true, float64(-1), invalid, invalid, true),
		newCoalescedTest(false, int64(0), invalid, invalid, true),
		newCoalescedTest(false, float64(0), invalid, invalid, true),
		newCoalescedTest(false, int64(1), invalid, invalid, false),
		newCoalescedTest(false, float64(1), invalid, invalid, false),
		newCoalescedTest(false, float64(1.0), invalid, invalid, false),
		newCoalescedTest(false, int64(-1), invalid, invalid, false),
		newCoalescedTest(false, float64(-1), invalid, invalid, false),
		newCoalescedTest(true, ast.Number{Value: int64(0)}, invalid, invalid, false),
		newCoalescedTest(true, ast.Number{Value: float64(0)}, invalid, invalid, false),
		newCoalescedTest(true, ast.Number{Value: int64(1)}, invalid, invalid, true),
		newCoalescedTest(true, ast.Number{Value: float64(1)}, invalid, invalid, true),
		newCoalescedTest(true, ast.Number{Value: float64(1.0)}, invalid, invalid, true),
		newCoalescedTest(true, ast.Number{Value: int64(-1)}, invalid, invalid, true),
		newCoalescedTest(true, ast.Number{Value: float64(-1)}, invalid, invalid, true),
		newCoalescedTest(false, ast.Number{Value: int64(0)}, invalid, invalid, true),
		newCoalescedTest(false, ast.Number{Value: float64(0)}, invalid, invalid, true),
		newCoalescedTest(false, ast.Number{Value: int64(1)}, invalid, invalid, false),
		newCoalescedTest(false, ast.Number{Value: float64(1)}, invalid, invalid, false),
		newCoalescedTest(false, ast.Number{Value: float64(1.0)}, invalid, invalid, false),
		newCoalescedTest(false, ast.Number{Value: int64(-1)}, invalid, invalid, false),
		newCoalescedTest(false, ast.Number{Value: float64(-1)}, invalid, invalid, false),

		newCoalescedTest(true, "", invalid, invalid, false),
		newCoalescedTest(true, " ", invalid, invalid, true),
		newCoalescedTest(true, "test", invalid, invalid, true),
		newCoalescedTest(false, "", invalid, invalid, true),
		newCoalescedTest(false, " ", invalid, invalid, false),
		newCoalescedTest(false, "test", invalid, invalid, false),
		newCoalescedTest(true, ast.String(""), invalid, invalid, false),
		newCoalescedTest(true, ast.String(" "), invalid, invalid, true),
		newCoalescedTest(true, ast.String("test"), invalid, invalid, true),
		newCoalescedTest(false, ast.String(""), invalid, invalid, true),
		newCoalescedTest(false, ast.String(" "), invalid, invalid, false),
		newCoalescedTest(false, ast.String("test"), invalid, invalid, false),

		newCoalescedTest(true, []any{}, invalid, invalid, false),
		newCoalescedTest(true, []any{0}, invalid, invalid, true),
		newCoalescedTest(true, []any{1}, invalid, invalid, true),
		newCoalescedTest(true, []any{""}, invalid, invalid, true),
		newCoalescedTest(false, []any{}, invalid, invalid, true),
		newCoalescedTest(false, []any{0}, invalid, invalid, false),
		newCoalescedTest(false, []any{1}, invalid, invalid, false),
		newCoalescedTest(false, []any{""}, invalid, invalid, false),
		newCoalescedTest(true, ast.Vector{Data: []any{}}, invalid, invalid, false),
		newCoalescedTest(true, ast.Vector{Data: []any{0}}, invalid, invalid, true),
		newCoalescedTest(true, ast.Vector{Data: []any{1}}, invalid, invalid, true),
		newCoalescedTest(true, ast.Vector{Data: []any{""}}, invalid, invalid, true),
		newCoalescedTest(false, ast.Vector{Data: []any{}}, invalid, invalid, true),
		newCoalescedTest(false, ast.Vector{Data: []any{0}}, invalid, invalid, false),
		newCoalescedTest(false, ast.Vector{Data: []any{1}}, invalid, invalid, false),
		newCoalescedTest(false, ast.Vector{Data: []any{""}}, invalid, invalid, false),

		newCoalescedTest(true, map[string]any{}, invalid, invalid, false),
		newCoalescedTest(true, map[string]any{"": ""}, invalid, invalid, true),
		newCoalescedTest(false, map[string]any{}, invalid, invalid, true),
		newCoalescedTest(false, map[string]any{"": ""}, invalid, invalid, false),
		newCoalescedTest(true, ast.Object{Data: nil}, invalid, invalid, false),
		newCoalescedTest(true, ast.Object{Data: map[string]any{}}, invalid, invalid, false),
		newCoalescedTest(true, ast.Object{Data: map[string]any{"": ""}}, invalid, invalid, true),
		newCoalescedTest(false, ast.Object{Data: nil}, invalid, invalid, true),
		newCoalescedTest(false, ast.Object{Data: map[string]any{}}, invalid, invalid, true),
		newCoalescedTest(false, ast.Object{Data: map[string]any{"": ""}}, invalid, invalid, false),

		///////////////////////////////////////////////////////////
		// test numbers against all other types, except nils and bools

		newCoalescedTest(0, 0, true, true, true),
		newCoalescedTest(0, 1, false, false, false),
		newCoalescedTest(0, -1, false, false, false),
		newCoalescedTest(0, ast.Number{Value: -1}, false, false, false),
		newCoalescedTest(2, 2, true, true, true),
		newCoalescedTest(0, 0.0, invalid, true, true),
		newCoalescedTest(0, 1.0, invalid, false, false),
		newCoalescedTest(2, 2.0, invalid, true, true),
		newCoalescedTest(2, ast.Number{Value: 2.0}, invalid, true, true),
		newCoalescedTest(-3.14, -3.14, true, true, true),
		newCoalescedTest(-3.14, ast.Number{Value: -3.14}, true, true, true),
		newCoalescedTest(ast.Number{Value: -3.14}, ast.Number{Value: -3.14}, true, true, true),

		newCoalescedTest(0, "", invalid, invalid, true),
		newCoalescedTest(0.0, "", invalid, invalid, true),
		newCoalescedTest(0, " ", invalid, invalid, true),
		newCoalescedTest(0.0, " ", invalid, invalid, true),
		newCoalescedTest(0, "0", invalid, invalid, true),
		newCoalescedTest(0.0, "0", invalid, invalid, true),
		newCoalescedTest(0, "0000", invalid, invalid, true),
		newCoalescedTest(0, "1", invalid, invalid, false),
		newCoalescedTest(1, "1", invalid, invalid, true),
		newCoalescedTest(1, " 1 ", invalid, invalid, true),
		newCoalescedTest(3, "3", invalid, invalid, true),
		newCoalescedTest(3.1, "3.1", invalid, invalid, true),
		newCoalescedTest(3.1, ast.String("3.1"), invalid, invalid, true),

		newCoalescedTest(0, []any{}, invalid, invalid, invalid),
		newCoalescedTest(0, []any{0}, invalid, invalid, invalid),
		newCoalescedTest(0.0, []any{}, invalid, invalid, invalid),
		newCoalescedTest(0.0, []any{0}, invalid, invalid, invalid),
		newCoalescedTest(1, []any{}, invalid, invalid, invalid),
		newCoalescedTest(1, []any{0}, invalid, invalid, invalid),
		newCoalescedTest(3.14, []any{}, invalid, invalid, invalid),
		newCoalescedTest(3.14, []any{0}, invalid, invalid, invalid),
		newCoalescedTest(0, ast.Vector{Data: []any{}}, invalid, invalid, invalid),
		newCoalescedTest(0, ast.Vector{Data: []any{0}}, invalid, invalid, invalid),
		newCoalescedTest(0.0, ast.Vector{Data: []any{}}, invalid, invalid, invalid),
		newCoalescedTest(0.0, ast.Vector{Data: []any{0}}, invalid, invalid, invalid),
		newCoalescedTest(1, ast.Vector{Data: []any{}}, invalid, invalid, invalid),
		newCoalescedTest(1, ast.Vector{Data: []any{0}}, invalid, invalid, invalid),
		newCoalescedTest(3.14, ast.Vector{Data: []any{}}, invalid, invalid, invalid),
		newCoalescedTest(3.14, ast.Vector{Data: []any{0}}, invalid, invalid, invalid),

		newCoalescedTest(0, map[string]any{}, invalid, invalid, invalid),
		newCoalescedTest(0, map[string]any{"": ""}, invalid, invalid, invalid),
		newCoalescedTest(0.0, map[string]any{}, invalid, invalid, invalid),
		newCoalescedTest(0.0, map[string]any{"": ""}, invalid, invalid, invalid),
		newCoalescedTest(1, map[string]any{}, invalid, invalid, invalid),
		newCoalescedTest(1, map[string]any{"": ""}, invalid, invalid, invalid),
		newCoalescedTest(3.14, map[string]any{}, invalid, invalid, invalid),
		newCoalescedTest(3.14, map[string]any{"": ""}, invalid, invalid, invalid),
		newCoalescedTest(0, ast.Object{Data: map[string]any{}}, invalid, invalid, invalid),
		newCoalescedTest(0, ast.Object{Data: map[string]any{"": ""}}, invalid, invalid, invalid),
		newCoalescedTest(0.0, ast.Object{Data: map[string]any{}}, invalid, invalid, invalid),
		newCoalescedTest(0.0, ast.Object{Data: map[string]any{"": ""}}, invalid, invalid, invalid),
		newCoalescedTest(1, ast.Object{Data: map[string]any{}}, invalid, invalid, invalid),
		newCoalescedTest(1, ast.Object{Data: map[string]any{"": ""}}, invalid, invalid, invalid),
		newCoalescedTest(3.14, ast.Object{Data: map[string]any{}}, invalid, invalid, invalid),
		newCoalescedTest(3.14, ast.Object{Data: map[string]any{"": ""}}, invalid, invalid, invalid),

		///////////////////////////////////////////////////////////
		// test strings against all other types, except nils, bools and numbers

		newCoalescedTest("", "", true, true, true),
		newCoalescedTest("", " ", false, false, false),
		newCoalescedTest("", "0", false, false, false),
		newCoalescedTest("", "a", false, false, false),
		newCoalescedTest("a", "a", true, true, true),
		newCoalescedTest("a", "A", false, false, false),
		newCoalescedTest("a", " a ", false, false, false),

		newCoalescedTest("", []any{}, invalid, invalid, invalid),
		newCoalescedTest("", []any{0}, invalid, invalid, invalid),
		newCoalescedTest("0", []any{}, invalid, invalid, invalid),
		newCoalescedTest("0", []any{0}, invalid, invalid, invalid),
		newCoalescedTest("", ast.Vector{Data: []any{}}, invalid, invalid, invalid),
		newCoalescedTest("", ast.Vector{Data: []any{0}}, invalid, invalid, invalid),
		newCoalescedTest("0", ast.Vector{Data: []any{}}, invalid, invalid, invalid),
		newCoalescedTest("0", ast.Vector{Data: []any{0}}, invalid, invalid, invalid),

		newCoalescedTest("", map[string]any{}, invalid, invalid, invalid),
		newCoalescedTest("", map[string]any{"": ""}, invalid, invalid, invalid),
		newCoalescedTest("0", map[string]any{}, invalid, invalid, invalid),
		newCoalescedTest("0", map[string]any{"": ""}, invalid, invalid, invalid),
		newCoalescedTest("", ast.Object{Data: map[string]any{}}, invalid, invalid, invalid),
		newCoalescedTest("", ast.Object{Data: map[string]any{"": ""}}, invalid, invalid, invalid),
		newCoalescedTest("0", ast.Object{Data: map[string]any{}}, invalid, invalid, invalid),
		newCoalescedTest("0", ast.Object{Data: map[string]any{"": ""}}, invalid, invalid, invalid),

		///////////////////////////////////////////////////////////
		// test vectors against all other types, except nils, bools, numbers and strings

		newCoalescedTest([]any{}, []any{}, true, true, true),
		newCoalescedTest([]any{}, []any{0}, false, false, false),
		newCoalescedTest([]any{0}, []any{0}, true, true, true),
		newCoalescedTest([]any{0}, []any{"0"}, invalid, invalid, true),
		newCoalescedTest([]any{false}, []any{0}, invalid, invalid, true),
		newCoalescedTest([]any{}, ast.Vector{Data: []any{}}, true, true, true),
		newCoalescedTest([]any{}, ast.Vector{Data: []any{0}}, false, false, false),
		newCoalescedTest([]any{0}, ast.Vector{Data: []any{0}}, true, true, true),
		newCoalescedTest([]any{0}, ast.Vector{Data: []any{"0"}}, invalid, invalid, true),
		newCoalescedTest([]any{false}, ast.Vector{Data: []any{0}}, invalid, invalid, true),

		newCoalescedTest([]any{}, map[string]any{}, invalid, invalid, true),
		newCoalescedTest([]any{}, map[string]any{"": ""}, invalid, invalid, invalid),
		newCoalescedTest([]any{0}, map[string]any{"": ""}, invalid, invalid, invalid),
		newCoalescedTest([]any{1}, map[string]any{}, invalid, invalid, invalid),
		newCoalescedTest([]any{nil}, map[string]any{"": ""}, invalid, invalid, invalid),
		newCoalescedTest([]any{}, ast.Object{Data: map[string]any{}}, invalid, invalid, true),
		newCoalescedTest([]any{}, ast.Object{Data: map[string]any{"": ""}}, invalid, invalid, invalid),
		newCoalescedTest([]any{0}, ast.Object{Data: map[string]any{"": ""}}, invalid, invalid, invalid),
		newCoalescedTest([]any{nil}, ast.Object{Data: map[string]any{"": ""}}, invalid, invalid, invalid),

		///////////////////////////////////////////////////////////
		// test objects

		newCoalescedTest(map[string]any{}, map[string]any{}, true, true, true),
		newCoalescedTest(map[string]any{}, map[string]any{"foo": "bar"}, false, false, false),
		newCoalescedTest(map[string]any{"foo": false}, map[string]any{"foo": ""}, invalid, invalid, true),
		newCoalescedTest(map[string]any{"foo": "bar"}, map[string]any{"foo": "bar"}, true, true, true),
		newCoalescedTest(map[string]any{"foo": "bar"}, map[string]any{"foo": "X"}, false, false, false),
		newCoalescedTest(map[string]any{}, ast.Object{Data: map[string]any{}}, true, true, true),
		newCoalescedTest(map[string]any{}, ast.Object{Data: map[string]any{"foo": "bar"}}, false, false, false),
		newCoalescedTest(map[string]any{"foo": false}, ast.Object{Data: map[string]any{"foo": ""}}, invalid, invalid, true),
		newCoalescedTest(map[string]any{"foo": "bar"}, ast.Object{Data: map[string]any{"foo": "bar"}}, true, true, true),
		newCoalescedTest(map[string]any{"foo": "bar"}, ast.Object{Data: map[string]any{"foo": "X"}}, false, false, false),
	}
}

func TestEqualCoalesced(t *testing.T) {
	pedanticCoalescer := coalescing.NewPedantic()
	strictCoalescer := coalescing.NewStrict()
	humaneCoalescer := coalescing.NewHumane()

	type subtest struct {
		left     any
		right    any
		coal     coalescing.Coalescer
		expected any
	}

	for _, testcase := range getEqualCoalescedTestcases() {
		t.Run(fmt.Sprintf("%v %v", testcase.left, testcase.right), func(t *testing.T) {
			subtests := []subtest{
				{
					left:     testcase.left,
					right:    testcase.right,
					coal:     pedanticCoalescer,
					expected: testcase.pedantic,
				},
				{
					left:     testcase.left,
					right:    testcase.right,
					coal:     strictCoalescer,
					expected: testcase.strict,
				},
				{
					left:     testcase.left,
					right:    testcase.right,
					coal:     humaneCoalescer,
					expected: testcase.humane,
				},
			}

			for _, subtest := range subtests {
				_, expectErr := subtest.expected.(invalidConversion)

				equal, err := EqualCoalesced(subtest.coal, subtest.left, subtest.right)
				if err != nil {
					if !expectErr {
						t.Errorf("%T unexpectedly failed: %v (%T) == %v (%T): %v", subtest.coal, subtest.left, subtest.left, subtest.right, subtest.right, err)
					}
				} else {
					if expectErr {
						t.Errorf("Expected %T to fail on %v (%T) == %v (%T), but got %v", subtest.coal, subtest.left, subtest.left, subtest.right, subtest.right, equal)
					} else if equal != subtest.expected {
						t.Errorf("Expected %T to return %v (%T) == %v (%T) => %v", subtest.coal, subtest.left, subtest.left, subtest.right, subtest.right, subtest.expected)
					}
				}

				// comparisons must be associated (a == b means b == a)
				flippedEqual, err := EqualCoalesced(subtest.coal, subtest.right, subtest.left)
				if err != nil {
					if !expectErr {
						t.Errorf("%T unexpectedly failed on reverse test: %v (%T) == %v (%T): %v", subtest.coal, subtest.right, subtest.right, subtest.left, subtest.left, err)
					}
				} else {
					if expectErr {
						t.Errorf("Expected %T to fail on %v (%T) == %v (%T), but got %v", subtest.coal, subtest.left, subtest.left, subtest.right, subtest.right, flippedEqual)
					} else if equal != flippedEqual {
						t.Errorf("Expected %T to be associative, but is not", subtest.coal)
					}
				}
			}
		})
	}
}
