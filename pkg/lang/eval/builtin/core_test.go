// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

type coreTestcase struct {
	expr      string
	expected  any
	document  any
	variables types.Variables
	invalid   bool
}

func (tc *coreTestcase) Test(t *testing.T) {
	t.Helper()

	result, err := runExpression(t, tc.expr, tc.document, tc.variables)
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

func TestIfFunction(t *testing.T) {
	testcases := []coreTestcase{
		{
			expr:    `(if)`,
			invalid: true,
		},
		{
			expr:    `(if true)`,
			invalid: true,
		},
		{
			expr:    `(if true "yes" "no" "extra")`,
			invalid: true,
		},
		{
			expr:     `(if true 3)`,
			expected: int64(3),
		},
		{
			expr:     `(if (eq 1 1) 3)`,
			expected: int64(3),
		},
		{
			expr:     `(if (eq 1 2) 3)`,
			expected: nil,
		},
		{
			expr:     `(if (eq 1 2) "yes" "else")`,
			expected: "else",
		},
		{
			expr:     `(if false "yes" (+ 1 4))`,
			expected: int64(5),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestDoFunction(t *testing.T) {
	testcases := []coreTestcase{
		{
			expr:    `(do)`,
			invalid: true,
		},
		{
			expr:     `(do 3)`,
			expected: int64(3),
		},

		// test that the runtime context is inherited from one step to another
		{
			expr:     `(do (set $var "foo") $var)`,
			expected: "foo",
		},
		{
			expr:     `(do (set $var "foo") $var (set $var "new") $var)`,
			expected: "new",
		},

		// test that the runtime context doesn't leak
		{
			expr:     `(set $var "outer") (do (set $var "inner")) (concat $var ["1" "2"])`,
			expected: "1outer2",
		},
		{
			expr:    `(do (set $var "inner")) (concat $var ["1" "2"])`,
			invalid: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestDefaultFunction(t *testing.T) {
	testcases := []coreTestcase{
		{
			expr:    `(default)`,
			invalid: true,
		},
		{
			expr:    `(default true)`,
			invalid: true,
		},
		{
			expr:     `(default null 3)`,
			expected: int64(3),
		},

		// coalescing should be applied

		{
			expr:     `(default false 3)`,
			expected: int64(3),
		},
		{
			expr:     `(default [] 3)`,
			expected: int64(3),
		},

		// errors are not swallowed

		{
			expr:    `(default (eq 3 "foo") 3)`,
			invalid: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestTryFunction(t *testing.T) {
	testcases := []coreTestcase{
		{
			expr:    `(try)`,
			invalid: true,
		},
		{
			expr:     `(try (+ 1 2))`,
			expected: int64(3),
		},

		// coalescing should be not applied

		{
			expr:     `(try false)`,
			expected: false,
		},
		{
			expr:     `(try null)`,
			expected: nil,
		},
		{
			expr:     `(try null "fallback")`,
			expected: nil,
		},

		// swallow errors

		{
			expr:     `(try (eq 3 "foo"))`,
			expected: nil,
		},
		{
			expr:     `(try (eq 3 "foo") "fallback")`,
			expected: "fallback",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestIsEmptyFunction(t *testing.T) {
	testcases := []coreTestcase{
		{
			expr:    `(empty?)`,
			invalid: true,
		},
		{
			expr:    `(empty? "too" "many")`,
			invalid: true,
		},
		{
			expr:    `(empty? ident)`,
			invalid: true,
		},
		{
			expr:     `(empty? null)`,
			expected: true,
		},
		{
			expr:     `(empty? true)`,
			expected: false,
		},
		{
			expr:     `(empty? false)`,
			expected: true,
		},
		{
			expr:     `(empty? 0)`,
			expected: true,
		},
		{
			expr:     `(empty? 0.0)`,
			expected: true,
		},
		{
			expr:     `(empty? (+ 0 0.0))`,
			expected: true,
		},
		{
			expr:     `(empty? (+ 1 0.0))`,
			expected: false,
		},
		{
			expr:     `(empty? [])`,
			expected: true,
		},
		{
			expr:     `(empty? [""])`,
			expected: false,
		},
		{
			expr:     `(empty? {})`,
			expected: true,
		},
		{
			expr:     `(empty? {foo "bar"})`,
			expected: false,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestHasFunction(t *testing.T) {
	testObjDocument := map[string]any{
		"aString": "foo",
		"aList":   []any{"first", 2, "third"},
		"aBool":   true,
		"anObject": map[string]any{
			"key1": true,
			"key2": nil,
			"key3": []any{9, map[string]any{"foo": "bar"}, 7},
		},
	}

	testArrDocument := []any{1, 2, map[string]any{"foo": "bar"}}

	testcases := []coreTestcase{
		{
			expr:    `(has?)`,
			invalid: true,
		},
		{
			expr:    `(has? "too" "many")`,
			invalid: true,
		},
		{
			expr:    `(has? true)`,
			invalid: true,
		},
		{
			expr:    `(has? (+ 1 2))`,
			invalid: true,
		},
		{
			expr:    `(has? "string")`,
			invalid: true,
		},
		{
			expr:    `(has? .[5.6])`,
			invalid: true,
		},

		// access the global document

		{
			expr:     `(has? .)`,
			expected: true,
			document: nil, // the . always matches, no matter what the document is
		},
		{
			expr:     `(has? .)`,
			expected: true,
			document: testObjDocument,
		},
		{
			expr:     `(has? .nonexistingKey)`,
			expected: false,
			document: testObjDocument,
		},
		{
			expr:     `(has? .[0])`,
			expected: false,
			document: testObjDocument,
		},
		{
			expr:     `(has? .aString)`,
			expected: true,
			document: testObjDocument,
		},
		{
			expr:     `(has? .["aString"])`,
			expected: true,
			document: testObjDocument,
		},
		{
			expr:     `(has? .[(concat "" "a" "String")])`,
			expected: true,
			document: testObjDocument,
		},
		{
			expr:     `(has? .aBool)`,
			expected: true,
			document: testObjDocument,
		},
		{
			expr:     `(has? .aList)`,
			expected: true,
			document: testObjDocument,
		},
		{
			expr:     `(has? .aList[0])`,
			expected: true,
			document: testObjDocument,
		},
		{
			expr:     `(has? .aList[99])`,
			expected: false,
			document: testObjDocument,
		},
		{
			expr:     `(has? .aList.invalidObjKey)`,
			expected: false,
			document: testObjDocument,
		},
		{
			expr:     `(has? .anObject)`,
			expected: true,
			document: testObjDocument,
		},
		{
			expr:     `(has? .anObject[99])`,
			expected: false,
			document: testObjDocument,
		},
		{
			expr:     `(has? .anObject.key1)`,
			expected: true,
			document: testObjDocument,
		},
		{
			expr:     `(has? .anObject.key99)`,
			expected: false,
			document: testObjDocument,
		},
		{
			expr:     `(has? .anObject.key3[1].foo)`,
			expected: true,
			document: testObjDocument,
		},
		{
			expr:     `(has? .anObject.key3[1].bar)`,
			expected: false,
			document: testObjDocument,
		},

		// global document is an array

		{
			expr:     `(has? .[1])`,
			expected: true,
			document: testArrDocument,
		},
		{
			expr:     `(has? .key)`,
			expected: false,
			document: testArrDocument,
		},
		{
			expr:     `(has? .[2].foo)`,
			expected: true,
			document: testArrDocument,
		},

		// global document is a scalar

		{
			expr:     `(has? .)`,
			expected: true,
			document: "testdata",
		},
		{
			expr:     `(has? .)`,
			expected: true,
			document: nil,
		},
		{
			expr:     `(has? .foo)`,
			expected: false,
			document: "testdata",
		},
		{
			expr:     `(has? .)`,
			expected: true,
			document: 64,
		},
		{
			expr:     `(has? .)`,
			expected: true,
			document: ast.String("foo"),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestRangeFunction(t *testing.T) {
	testcases := []coreTestcase{
		{
			// missing everything
			expr:    `(range)`,
			invalid: true,
		},
		{
			// missing naming vector
			expr:    `(range [1 2 3])`,
			invalid: true,
		},
		{
			// missing naming vector
			expr:    `(range [1 2 3] (+ 1 2))`,
			invalid: true,
		},
		{
			// naming vector must be 1 or 2 elements long
			expr:    `(range [1 2 3] [] (+ 1 2))`,
			invalid: true,
		},
		{
			// naming vector must be 1 or 2 elements long
			expr:    `(range [1 2 3] [a b c] (+ 1 2))`,
			invalid: true,
		},
		{
			// do not allow numbers in the naming vector
			expr:    `(range [1 2 3] [1 2] (+ 1 2))`,
			invalid: true,
		},
		{
			// do not allow strings in naming vector
			expr:    `(range [1 2 3] ["foo" "bar"] (+ 1 2))`,
			invalid: true,
		},
		{
			// cannot range over non-vectors/objects
			expr:    `(range "invalid" [a] (+ 1 2))`,
			invalid: true,
		},
		{
			// cannot range over non-vectors/objects
			expr:    `(range 5 [a] (+ 1 2))`,
			invalid: true,
		},
		{
			// single simple expression
			expr:     `(range [1 2 3] [a] (+ 1 2))`,
			expected: int64(3),
		},
		{
			// multiple expressions that use a common context
			expr:     `(range [1 2 3] [a] (set $foo $a) (+ $foo 3))`,
			expected: int64(6),
		},
		{
			// count iterations
			expr:     `(range [1 2 3] [loop-var] (set $counter (+ (default (try $counter) 0) 1)))`,
			expected: int64(3),
		},
		{
			// value is bound to desired variable
			expr:     `(range [1 2 3] [a] $a)`,
			expected: int64(3),
		},
		{
			// support loop index variable
			expr:     `(range [1 2 3] [idx var] $idx)`,
			expected: int64(2),
		},
		{
			// support loop index variable
			expr:     `(range [1 2 3] [idx var] $var)`,
			expected: int64(3),
		},
		{
			// variables do not leak outside the range
			expr:    `(range [1 2 3] [idx var] $idx) (+ $var 0)`,
			invalid: true,
		},
		{
			// variables do not leak outside the range
			expr:    `(range [1 2 3] [idx var] $idx) (+ $idx 0)`,
			invalid: true,
		},
		{
			// support ranging over objects
			expr:     `(range {} [key value] $key)`,
			expected: nil,
		},
		{
			expr:     `(range {foo "bar"} [key value] $key)`,
			expected: "foo",
		},
		{
			expr:     `(range {foo "bar"} [key value] $value)`,
			expected: "bar",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}
