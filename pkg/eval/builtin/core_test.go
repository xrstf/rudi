// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"

	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"

	"github.com/google/go-cmp/cmp"
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

	if _, ok := tc.expected.([]any); ok {
		if !cmp.Equal(tc.expected, result) {
			t.Fatalf("Expected %v (%T), but got %v (%T)", tc.expected, tc.expected, result, result)
		}
	} else if _, ok := tc.expected.(map[string]any); ok {
		if !cmp.Equal(tc.expected, result) {
			t.Fatalf("Expected %v (%T), but got %v (%T)", tc.expected, tc.expected, result, result)
		}
	} else {
		if result != tc.expected {
			t.Fatalf("Expected %v (%T), but got %v (%T)", tc.expected, tc.expected, result, result)
		}
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
			expr:    `(if identifier "yes")`,
			invalid: true,
		},
		{
			expr:    `(if {} "yes")`,
			invalid: true,
		},
		{
			expr:    `(if [] "yes")`,
			invalid: true,
		},
		{
			expr:    `(if 1 "yes")`,
			invalid: true,
		},
		{
			expr:    `(if 3.4 "yes")`,
			invalid: true,
		},
		{
			expr:    `(if (+ 1 1) "yes")`,
			invalid: true,
		},
		{
			expr:     `(if true 3)`,
			expected: int64(3),
		},
		{
			expr:     `(if (eq? 1 1) 3)`,
			expected: int64(3),
		},
		{
			expr:     `(if (eq? 1 2) 3)`,
			expected: nil,
		},
		{
			expr:     `(if (eq? 1 2) "yes" "else")`,
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

func TestSetFunction(t *testing.T) {
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

	testVecDocument := []any{1, 2, map[string]any{"foo": "bar"}}

	testVariables := types.Variables{
		"myvar":  42,
		"obj":    testObjDocument,
		"vec":    testVecDocument,
		"astVec": ast.Vector{Data: []any{ast.String("foo")}},
	}

	testcases := []coreTestcase{
		{
			expr:    `(set)`,
			invalid: true,
		},
		{
			expr:    `(set true)`,
			invalid: true,
		},
		{
			expr:    `(set "foo")`,
			invalid: true,
		},
		{
			expr:    `(set 42)`,
			invalid: true,
		},
		{
			expr:    `(set {foo "bar"})`,
			invalid: true,
		},
		{
			expr:    `(set $var)`,
			invalid: true,
		},
		{
			expr:    `(set $var "too" "many")`,
			invalid: true,
		},
		{
			expr:    `(set $var (unknown-func))`,
			invalid: true,
		},
		// return the value that was set
		{
			expr:     `(set $var "foo")`,
			expected: "foo",
		},
		{
			expr:     `(set $var 1)`,
			expected: int64(1),
		},
		// can overwrite variables on the top level
		{
			expr:      `(set $myvar 12)`,
			expected:  int64(12),
			variables: testVariables,
		},
		// can change the type
		{
			expr:      `(set $myvar "new value")`,
			expected:  "new value",
			variables: testVariables,
		},
		{
			expr: `(set $obj.aList[1] "new value")`,
			expected: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", "new value", "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{int64(9), map[string]any{"foo": "bar"}, int64(7)},
				},
			},
			variables: testVariables,
		},
		// set itself does not change the first argument
		{
			expr:      `(set $myvar "new value") $myvar`,
			expected:  int64(42),
			variables: testVariables,
		},
		{
			expr:      `(set $obj.aString "new value") $obj.aString`,
			expected:  "foo",
			variables: testVariables,
		},
		{
			expr:      `(set $obj.aList[1] "new value") $obj.aList`,
			expected:  []any{"first", int64(2), "third"},
			variables: testVariables,
		},
		// // ...but not leak into upper scopes
		// {
		// 	expr:     `(set $a 1) (if true (set $a 2)) $a`,
		// 	expected: int64(1),
		// },
		// {
		// 	expr:    `(set $a 1) (if true (set $b 2)) $b`,
		// 	invalid: true,
		// },
		// // do not accidentally set a key without creating a new context
		// {
		// 	expr:     `(set $a {foo "bar"}) (if true (set $a.foo "updated"))`,
		// 	expected: "updated",
		// },
		// {
		// 	expr:     `(set $a {foo "bar"}) (if true (set $a.foo "updated")) $a.foo`,
		// 	expected: "bar",
		// },
		// // handle bad paths
		// {
		// 	expr:    `(set $obj[5.6] "new value")`,
		// 	invalid: true,
		// },
		// // not a vector
		// {
		// 	expr:    `(set $obj[5] "new value")`,
		// 	invalid: true,
		// },
		// {
		// 	expr:    `(set $obj.aBool[5] "new value")`,
		// 	invalid: true,
		// },
		// // update a key within an object variable
		// {
		// 	expr:      `(set $obj.aString "new value")`,
		// 	expected:  "new value",
		// 	variables: testVariables,
		// },
		// {
		// 	expr:      `(set $obj.aString "new value") $obj.aString`,
		// 	expected:  "new value",
		// 	variables: testVariables,
		// },
		// // add a new sub key
		// {
		// 	expr:      `(set $obj.newKey "new value")`,
		// 	expected:  "new value",
		// 	variables: testVariables,
		// },
		// {
		// 	expr:      `(set $obj.newKey "new value") $obj.newKey`,
		// 	expected:  "new value",
		// 	variables: testVariables,
		// },
		// // runtime variables
		// {
		// 	expr:     `(set $vec [1]) (set $vec[0] 2) $vec[0]`,
		// 	expected: int64(2),
		// },
		// // replace the global document
		// {
		// 	expr:     `(set . 1) .`,
		// 	document: testObjDocument,
		// 	expected: int64(1),
		// },
		// // update keys in the global document
		// {
		// 	expr:     `(set .aString "new-value") .aString`,
		// 	document: testObjDocument,
		// 	expected: "new-value",
		// },
		// // add new keys
		// {
		// 	expr:     `(set .newKey "new-value") .newKey`,
		// 	document: testObjDocument,
		// 	expected: "new-value",
		// },
		// // update vectors
		// {
		// 	expr:     `(set .aList[1] "new-value") .aList[1]`,
		// 	document: testObjDocument,
		// 	expected: "new-value",
		// },
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}

func TestDeleteFunction(t *testing.T) {
	testcases := []coreTestcase{
		{
			expr:    `(delete)`,
			invalid: true,
		},
		{
			expr:    `(delete "too" "many")`,
			invalid: true,
		},
		{
			expr:    `(delete $var)`,
			invalid: true,
		},
		{
			// TODO: This should be valid.
			expr:    `(delete [1 2 3][1])`,
			invalid: true,
		},
		{
			// TODO: This should be valid.
			expr:    `(delete {foo "bar"}.foo)`,
			invalid: true,
		},
		// allow removing everything
		{
			expr:     `(delete .)`,
			document: map[string]any{"foo": "bar"},
			expected: nil,
		},
		{
			expr:     `(delete .)`,
			document: map[string]any{"foo": "bar"},
			expected: nil,
		},
		// delete does not update the target
		{
			expr:     `(delete .) .`,
			document: map[string]any{"foo": "bar"},
			expected: map[string]any{"foo": "bar"},
		},
		// can remove a key
		{
			expr:     `(delete .foo)`,
			document: map[string]any{"foo": "bar"},
			expected: map[string]any{},
		},
		// non-existent key is okay
		{
			expr:     `(delete .bar)`,
			document: map[string]any{"foo": "bar"},
			expected: map[string]any{"foo": "bar"},
		},
		// path must be sane though
		{
			expr:     `(delete .[1])`,
			document: map[string]any{"foo": "bar"},
			invalid:  true,
		},
		// can delete from array
		{
			expr:     `(delete .[1])`,
			document: []any{"a", "b", "c"},
			expected: []any{"a", "c"},
		},
		// vector bounds are checked
		{
			expr:     `(delete .[-1])`,
			document: []any{"a", "b", "c"},
			invalid:  true,
		},
		{
			expr:     `(delete .[3])`,
			document: []any{"a", "b", "c"},
			invalid:  true,
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
			expr:    `(do identifier)`,
			invalid: true,
		},
		{
			expr:     `(do 3)`,
			expected: int64(3),
		},

		// test that the runtime context is inherited from one step to another
		{
			expr:     `(do (set! $var "foo") $var)`,
			expected: "foo",
		},
		{
			expr:     `(do (set! $var "foo") $var (set! $var "new") $var)`,
			expected: "new",
		},

		// test that the runtime context doesn't leak
		{
			expr:     `(set! $var "outer") (do (set! $var "inner")) (concat $var ["1" "2"])`,
			expected: "1outer2",
		},
		{
			expr:    `(do (set! $var "inner")) (concat $var ["1" "2"])`,
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
			expr:    `(default (eq? 3 "foo") 3)`,
			invalid: true,
		},

		{
			expr:    `(default false (eq? 3 "foo"))`,
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
			expr:     `(try (eq? 3 "foo"))`,
			expected: nil,
		},
		{
			expr:     `(try (eq? 3 "foo") "fallback")`,
			expected: "fallback",
		},

		// not in the fallback though

		{
			expr:    `(try (eq? 3 "foo") (eq? 3 "foo"))`,
			invalid: true,
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

	testVecDocument := []any{1, 2, map[string]any{"foo": "bar"}}

	testVariables := types.Variables{
		// value does not matter here, but this testcase is still meant
		// to ensure the missing path is detected, not detect an unknown variable
		"myvar":  42,
		"obj":    testObjDocument,
		"vec":    testVecDocument,
		"astVec": ast.Vector{Data: []any{ast.String("foo")}},
	}

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
		{
			expr:    `(has? (unknown-func).bar)`,
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
			document: testVecDocument,
		},
		{
			expr:     `(has? .key)`,
			expected: false,
			document: testVecDocument,
		},
		{
			expr:     `(has? .[2].foo)`,
			expected: true,
			document: testVecDocument,
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

		// follow a path expression on a variable

		{
			expr:    `(has? $myvar)`,
			invalid: true,
		},
		{
			// missing path expression (TODO: should this be valid?)
			expr:      `(has? $myvar)`,
			invalid:   true,
			variables: testVariables,
		},
		{
			expr:      `(has? $myvar.foo)`,
			expected:  false,
			variables: testVariables,
		},
		{
			expr:      `(has? $myvar[0])`,
			expected:  false,
			variables: testVariables,
		},
		{
			expr:      `(has? $obj.aString)`,
			expected:  true,
			variables: testVariables,
		},
		{
			expr:      `(has? $obj.aList[1])`,
			expected:  true,
			variables: testVariables,
		},
		{
			expr:      `(has? $vec[1])`,
			expected:  true,
			variables: testVariables,
		},
		{
			expr:      `(has? $astVec[0])`,
			expected:  true,
			variables: testVariables,
		},

		// follow a path expression on a vector node

		{
			expr:     `(has? [1 2 3][1])`,
			expected: true,
		},
		{
			expr:     `(has? [1 2 3][4])`,
			expected: false,
		},

		// follow a path expression on an object node

		{
			expr:     `(has? {foo "bar"}.foo)`,
			expected: true,
		},
		{
			expr:     `(has? {foo "bar"}.bar)`,
			expected: false,
		},

		// follow a path expression on a tuple node
		// (don't even need "set!" here)

		{
			expr:     `(has? (set $foo {foo "bar"}).foo)`,
			expected: true,
		},
		{
			expr:     `(has? (set $foo {foo "bar"}).bar)`,
			expected: false,
		},
		{
			expr:     `(has? (set $foo [1])[0])`,
			expected: true,
		},
		{
			expr:     `(has? (set $foo {foo "bar"})[0])`,
			expected: false,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.expr, testcase.Test)
	}
}
