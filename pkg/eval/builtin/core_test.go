// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"testing"

	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/testutil"
)

func TestIfFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(if)`,
			Invalid:    true,
		},
		{
			Expression: `(if true)`,
			Invalid:    true,
		},
		{
			Expression: `(if true "yes" "no" "extra")`,
			Invalid:    true,
		},
		{
			Expression: `(if identifier "yes")`,
			Invalid:    true,
		},
		{
			Expression: `(if {} "yes")`,
			Invalid:    true,
		},
		{
			Expression: `(if [] "yes")`,
			Invalid:    true,
		},
		{
			Expression: `(if 1 "yes")`,
			Invalid:    true,
		},
		{
			Expression: `(if 3.4 "yes")`,
			Invalid:    true,
		},
		{
			Expression: `(if (+ 1 1) "yes")`,
			Invalid:    true,
		},
		{
			Expression: `(if true 3)`,
			Expected:   ast.Number{Value: int64(3)},
		},
		{
			Expression: `(if (eq? 1 1) 3)`,
			Expected:   ast.Number{Value: int64(3)},
		},
		{
			Expression: `(if (eq? 1 2) 3)`,
			Expected:   ast.Null{},
		},
		{
			Expression: `(if (eq? 1 2) "yes" "else")`,
			Expected:   ast.String("else"),
		},
		{
			Expression: `(if false "yes" (+ 1 4))`,
			Expected:   ast.Number{Value: int64(5)},
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestSetFunction(t *testing.T) {
	testObjDocument := func() any {
		return map[string]any{
			"aString": "foo",
			"aList":   []any{"first", int64(2), "third"},
			"aBool":   true,
			"anObject": map[string]any{
				"key1": true,
				"key2": nil,
				"key3": []any{int64(9), map[string]any{"foo": "bar"}, int64(7)},
			},
		}
	}

	testVecDocument := func() any {
		return []any{int64(1), int64(2), map[string]any{"foo": "bar"}}
	}

	testVariables := func() types.Variables {
		return types.Variables{
			"myvar":  int64(42),
			"obj":    testObjDocument(),
			"vec":    testVecDocument(),
			"astVec": ast.Vector{Data: []any{ast.String("foo")}},
		}
	}

	testcases := []testutil.Testcase{
		{
			Expression: `(set)`,
			Invalid:    true,
		},
		{
			Expression: `(set true)`,
			Invalid:    true,
		},
		{
			Expression: `(set "foo")`,
			Invalid:    true,
		},
		{
			Expression: `(set 42)`,
			Invalid:    true,
		},
		{
			Expression: `(set {foo "bar"})`,
			Invalid:    true,
		},
		{
			Expression: `(set $var)`,
			Invalid:    true,
		},
		{
			Expression: `(set $var "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(set $var (unknown-func))`,
			Invalid:    true,
		},
		// return the value that was set (without a bang modifier, this doesn't actually modify anything)
		{
			Expression: `(set $var "foo")`,
			Expected:   ast.String("foo"),
		},
		{
			Expression: `(set $var "foo") $var`,
			Invalid:    true,
		},
		{
			Expression: `(set $var 1)`,
			Expected:   ast.Number{Value: int64(1)},
		},
		// with bang it works as expected
		{
			Expression: `(set! $var 1)`,
			Expected:   ast.Number{Value: int64(1)},
		},
		{
			Expression: `(set! $var 1) $var`,
			Expected:   ast.Number{Value: int64(1)},
		},
		// can overwrite variables on the top level
		{
			Expression: `(set! $myvar 12) $myvar`,
			Variables:  testVariables(),
			Expected:   ast.Number{Value: int64(12)},
		},
		// can change the type
		{
			Expression: `(set! $myvar "new value") $myvar`,
			Variables:  testVariables(),
			Expected:   ast.String("new value"),
		},
		{
			Expression: `(set! $obj.aList[1] "new value")`,
			Variables:  testVariables(),
			Expected:   ast.String("new value"),
		},
		{
			Expression: `(set! $obj.aList[1] "new value") $obj`,
			Variables:  testVariables(),
			Expected: ast.Object{Data: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", ast.String("new value"), "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{int64(9), map[string]any{"foo": "bar"}, int64(7)},
				},
			}},
		},
		// set itself does not change the first argument
		{
			Expression: `(set $myvar "new value") $myvar`,
			Variables:  testVariables(),
			Expected:   ast.Number{Value: int64(42)},
		},
		{
			Expression: `(set $obj.aString "new value") $obj.aString`,
			Variables:  testVariables(),
			Expected:   ast.String("foo"),
		},
		{
			Expression: `(set $obj.aList[1] "new value") $obj.aList`,
			Variables:  testVariables(),
			Expected:   ast.Vector{Data: []any{"first", int64(2), "third"}},
		},
		// ...but not leak into upper scopes
		{
			Expression: `(set! $a 1) (if true (set! $a 2)) $a`,
			Expected:   ast.Number{Value: int64(1)},
		},
		{
			Expression: `(set! $a 1) (if true (set! $b 2)) $b`,
			Invalid:    true,
		},
		// do not accidentally set a key without creating a new context
		{
			Expression: `(set! $a {foo "bar"}) (if true (set! $a.foo "updated"))`,
			Expected:   ast.String("updated"),
		},
		{
			Expression: `(set! $a {foo "bar"}) (if true (set! $a.foo "updated")) $a.foo`,
			Expected:   ast.String("bar"),
		},
		// handle bad paths
		{
			Expression: `(set! $obj[5.6] "new value")`,
			Invalid:    true,
		},
		// not a vector
		{
			Expression: `(set! $obj[5] "new value")`,
			Invalid:    true,
		},
		{
			Expression: `(set! $obj.aBool[5] "new value")`,
			Invalid:    true,
		},
		// update a key within an object variable
		{
			Expression: `(set! $obj.aString "new value")`,
			Expected:   ast.String("new value"),
			Variables:  testVariables(),
		},
		{
			Expression: `(set! $obj.aString "new value") $obj.aString`,
			Expected:   ast.String("new value"),
			Variables:  testVariables(),
		},
		// add a new sub key
		{
			Expression: `(set! $obj.newKey "new value")`,
			Expected:   ast.String("new value"),
			Variables:  testVariables(),
		},
		{
			Expression: `(set! $obj.newKey "new value") $obj.newKey`,
			Expected:   ast.String("new value"),
			Variables:  testVariables(),
		},
		// runtime variables
		{
			Expression: `(set! $vec [1]) (set! $vec[0] 2) $vec[0]`,
			Expected:   ast.Number{Value: int64(2)},
		},
		// replace the global document
		{
			Expression:       `(set! . 1) .`,
			Document:         testObjDocument(),
			Expected:         ast.Number{Value: int64(1)},
			ExpectedDocument: int64(1),
		},
		// update keys in the global document
		{
			Expression: `(set! .aString "new-value") .aString`,
			Document:   testObjDocument(),
			Expected:   ast.String("new-value"),
			ExpectedDocument: map[string]any{
				"aString": ast.String("new-value"),
				"aList":   []any{"first", int64(2), "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{int64(9), map[string]any{"foo": "bar"}, int64(7)},
				},
			},
		},
		// add new keys
		{
			Expression: `(set! .newKey "new-value") .newKey`,
			Document:   testObjDocument(),
			Expected:   ast.String("new-value"),
			ExpectedDocument: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", int64(2), "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{int64(9), map[string]any{"foo": "bar"}, int64(7)},
				},
				// TODO: Should we generally try to put native values in objects and vectors?
				"newKey": ast.String("new-value"),
			},
		},
		// update vectors
		{
			Expression: `(set! .aList[1] "new-value") .aList[1]`,
			Document:   testObjDocument(),
			Expected:   ast.String("new-value"),
			ExpectedDocument: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", ast.String("new-value"), "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{int64(9), map[string]any{"foo": "bar"}, int64(7)},
				},
			},
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestDeleteFunction(t *testing.T) {
	testObjDocument := func() map[string]any {
		return map[string]any{
			"aString": "foo",
			"aList":   []any{"first", int64(2), "third"},
			"aBool":   true,
			"anObject": map[string]any{
				"key1": true,
				"key2": nil,
				"key3": []any{int64(9), map[string]any{"foo": "bar"}, int64(7)},
			},
		}
	}

	testcases := []testutil.Testcase{
		{
			Expression: `(delete)`,
			Invalid:    true,
		},
		{
			Expression: `(delete "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(delete $var)`,
			Invalid:    true,
		},
		{
			// TODO: This should be valid.
			Expression: `(delete [1 2 3][1])`,
			Invalid:    true,
		},
		{
			// TODO: This should be valid.
			Expression: `(delete {foo "bar"}.foo)`,
			Invalid:    true,
		},
		// allow removing everything
		{
			Expression:       `(delete .)`,
			Document:         map[string]any{"foo": "bar"},
			Expected:         ast.Null{},
			ExpectedDocument: map[string]any{"foo": "bar"},
		},
		{
			Expression: `(delete! .) .`,
			Document:   map[string]any{"foo": "bar"},
			Expected:   ast.Null{},
		},
		// delete does not update the target
		{
			Expression:       `(delete .) .`,
			Document:         map[string]any{"foo": "bar"},
			Expected:         ast.Object{Data: map[string]any{"foo": "bar"}},
			ExpectedDocument: map[string]any{"foo": "bar"},
		},
		// can remove a key
		{
			Expression:       `(delete .foo)`,
			Document:         map[string]any{"foo": "bar"},
			Expected:         ast.Object{Data: map[string]any{}},
			ExpectedDocument: map[string]any{"foo": "bar"},
		},
		{
			Expression:       `(delete .foo) .`,
			Document:         map[string]any{"foo": "bar"},
			Expected:         ast.Object{Data: map[string]any{"foo": "bar"}},
			ExpectedDocument: map[string]any{"foo": "bar"},
		},
		{
			Expression:       `(delete! .foo) .`,
			Document:         map[string]any{"foo": "bar"},
			Expected:         ast.Object{Data: map[string]any{}},
			ExpectedDocument: map[string]any{},
		},
		// non-existent key is okay
		{
			Expression:       `(delete .bar)`,
			Document:         map[string]any{"foo": "bar"},
			Expected:         ast.Object{Data: map[string]any{"foo": "bar"}},
			ExpectedDocument: map[string]any{"foo": "bar"},
		},
		// path must be sane though
		{
			Expression: `(delete .[1])`,
			Document:   map[string]any{"foo": "bar"},
			Invalid:    true,
		},
		// can delete from array
		{
			Expression:       `(delete .[1])`,
			Document:         []any{"a", "b", "c"},
			Expected:         ast.Vector{Data: []any{"a", "c"}},
			ExpectedDocument: []any{"a", "b", "c"},
		},
		{
			Expression:       `(delete .[1]) .`,
			Document:         []any{"a", "b", "c"},
			Expected:         ast.Vector{Data: []any{"a", "b", "c"}},
			ExpectedDocument: []any{"a", "b", "c"},
		},
		{
			Expression:       `(delete! .[1]) .`,
			Document:         []any{"a", "b", "c"},
			Expected:         ast.Vector{Data: []any{"a", "c"}},
			ExpectedDocument: []any{"a", "c"},
		},
		// vector bounds are checked
		{
			Expression: `(delete .[-1])`,
			Document:   []any{"a", "b", "c"},
			Invalid:    true,
		},
		{
			Expression: `(delete .[3])`,
			Document:   []any{"a", "b", "c"},
			Invalid:    true,
		},
		// can delete sub keys
		{
			Expression: `(delete .aList[1])`,
			Document:   testObjDocument(),
			Expected: ast.Object{Data: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{int64(9), map[string]any{"foo": "bar"}, int64(7)},
				},
			}},
			ExpectedDocument: testObjDocument(),
		},
		{
			Expression:       `(delete .aList[1]) .`,
			Document:         testObjDocument(),
			Expected:         ast.Object{Data: testObjDocument()},
			ExpectedDocument: testObjDocument(),
		},
		{
			Expression: `(delete! .aList[1]) .`,
			Document:   testObjDocument(),
			Expected: ast.Object{Data: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{int64(9), map[string]any{"foo": "bar"}, int64(7)},
				},
			}},
			ExpectedDocument: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{int64(9), map[string]any{"foo": "bar"}, int64(7)},
				},
			},
		},
		{
			Expression: `(delete .anObject.key3[1].foo)`,
			Document:   testObjDocument(),
			Expected: ast.Object{Data: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", int64(2), "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{int64(9), map[string]any{}, int64(7)},
				},
			}},
			ExpectedDocument: testObjDocument(),
		},
		{
			Expression:       `(delete .anObject.key3[1].foo) .`,
			Document:         testObjDocument(),
			Expected:         ast.Object{Data: testObjDocument()},
			ExpectedDocument: testObjDocument(),
		},
		{
			Expression: `(delete! .anObject.key3[1].foo) .`,
			Document:   testObjDocument(),
			Expected: ast.Object{Data: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", int64(2), "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{int64(9), map[string]any{}, int64(7)},
				},
			}},
			ExpectedDocument: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", int64(2), "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{int64(9), map[string]any{}, int64(7)},
				},
			},
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestDoFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(do)`,
			Invalid:    true,
		},
		{
			Expression: `(do identifier)`,
			Invalid:    true,
		},
		{
			Expression: `(do 3)`,
			Expected:   ast.Number{Value: int64(3)},
		},

		// test that the runtime context is inherited from one step to another
		{
			Expression: `(do (set! $var "foo") $var)`,
			Expected:   ast.String("foo"),
		},
		{
			Expression: `(do (set! $var "foo") $var (set! $var "new") $var)`,
			Expected:   ast.String("new"),
		},

		// test that the runtime context doesn't leak
		{
			Expression: `(set! $var "outer") (do (set! $var "inner")) (concat $var ["1" "2"])`,
			Expected:   ast.String("1outer2"),
		},
		{
			Expression: `(do (set! $var "inner")) (concat $var ["1" "2"])`,
			Invalid:    true,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestDefaultFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(default)`,
			Invalid:    true,
		},
		{
			Expression: `(default true)`,
			Invalid:    true,
		},
		{
			Expression: `(default null 3)`,
			Expected:   ast.Number{Value: int64(3)},
		},

		// coalescing should be applied

		{
			Expression: `(default false 3)`,
			Expected:   ast.Number{Value: int64(3)},
		},
		{
			Expression: `(default [] 3)`,
			Expected:   ast.Number{Value: int64(3)},
		},

		// errors are not swallowed

		{
			Expression: `(default (eq? 3 "foo") 3)`,
			Invalid:    true,
		},

		{
			Expression: `(default false (eq? 3 "foo"))`,
			Invalid:    true,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestTryFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(try)`,
			Invalid:    true,
		},
		{
			Expression: `(try (+ 1 2))`,
			Expected:   ast.Number{Value: int64(3)},
		},

		// coalescing should be not applied

		{
			Expression: `(try false)`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(try null)`,
			Expected:   ast.Null{},
		},
		{
			Expression: `(try null "fallback")`,
			Expected:   ast.Null{},
		},

		// swallow errors

		{
			Expression: `(try (eq? 3 "foo"))`,
			Expected:   ast.Null{},
		},
		{
			Expression: `(try (eq? 3 "foo") "fallback")`,
			Expected:   ast.String("fallback"),
		},

		// not in the fallback though

		{
			Expression: `(try (eq? 3 "foo") (eq? 3 "foo"))`,
			Invalid:    true,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestIsEmptyFunction(t *testing.T) {
	testcases := []testutil.Testcase{
		{
			Expression: `(empty?)`,
			Invalid:    true,
		},
		{
			Expression: `(empty? "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(empty? ident)`,
			Invalid:    true,
		},
		{
			Expression: `(empty? null)`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(empty? true)`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(empty? false)`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(empty? 0)`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(empty? 0.0)`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(empty? (+ 0 0.0))`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(empty? (+ 1 0.0))`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(empty? [])`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(empty? [""])`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(empty? {})`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(empty? {foo "bar"})`,
			Expected:   ast.Bool(false),
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}

func TestHasFunction(t *testing.T) {
	testObjDocument := map[string]any{
		"aString": "foo",
		"aList":   []any{"first", int64(2), "third"},
		"aBool":   true,
		"anObject": map[string]any{
			"key1": true,
			"key2": nil,
			"key3": []any{int64(9), map[string]any{"foo": "bar"}, int64(7)},
		},
	}

	testVecDocument := []any{int64(1), int64(2), map[string]any{"foo": "bar"}}

	testVariables := types.Variables{
		// value does not matter here, but this testcase is still meant
		// to ensure the missing path is detected, not detect an unknown variable
		"myvar":  int64(42),
		"obj":    testObjDocument,
		"vec":    testVecDocument,
		"astVec": ast.Vector{Data: []any{ast.String("foo")}},
	}

	testcases := []testutil.Testcase{
		{
			Expression: `(has?)`,
			Invalid:    true,
		},
		{
			Expression: `(has? "too" "many")`,
			Invalid:    true,
		},
		{
			Expression: `(has? true)`,
			Invalid:    true,
		},
		{
			Expression: `(has? (+ 1 2))`,
			Invalid:    true,
		},
		{
			Expression: `(has? "string")`,
			Invalid:    true,
		},
		{
			Expression: `(has? .[5.6])`,
			Invalid:    true,
		},
		{
			Expression: `(has? (unknown-func).bar)`,
			Invalid:    true,
		},

		// access the global document

		{
			Expression: `(has? .)`,
			Expected:   ast.Bool(true),
			Document:   nil, // the . always matches, no matter what the document is
		},
		{
			Expression:       `(has? .)`,
			Expected:         ast.Bool(true),
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .nonexistingKey)`,
			Expected:         ast.Bool(false),
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .[0])`,
			Expected:         ast.Bool(false),
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .aString)`,
			Expected:         ast.Bool(true),
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .["aString"])`,
			Expected:         ast.Bool(true),
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .[(concat "" "a" "String")])`,
			Expected:         ast.Bool(true),
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .aBool)`,
			Expected:         ast.Bool(true),
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .aList)`,
			Expected:         ast.Bool(true),
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .aList[0])`,
			Expected:         ast.Bool(true),
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .aList[99])`,
			Expected:         ast.Bool(false),
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .aList.invalidObjKey)`,
			Expected:         ast.Bool(false),
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .anObject)`,
			Expected:         ast.Bool(true),
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .anObject[99])`,
			Expected:         ast.Bool(false),
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .anObject.key1)`,
			Expected:         ast.Bool(true),
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .anObject.key99)`,
			Expected:         ast.Bool(false),
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .anObject.key3[1].foo)`,
			Expected:         ast.Bool(true),
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .anObject.key3[1].bar)`,
			Expected:         ast.Bool(false),
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},

		// global document is an array

		{
			Expression:       `(has? .[1])`,
			Expected:         ast.Bool(true),
			Document:         testVecDocument,
			ExpectedDocument: testVecDocument,
		},
		{
			Expression:       `(has? .key)`,
			Expected:         ast.Bool(false),
			Document:         testVecDocument,
			ExpectedDocument: testVecDocument,
		},
		{
			Expression:       `(has? .[2].foo)`,
			Expected:         ast.Bool(true),
			Document:         testVecDocument,
			ExpectedDocument: testVecDocument,
		},

		// global document is a scalar

		{
			Expression:       `(has? .)`,
			Expected:         ast.Bool(true),
			Document:         "testdata",
			ExpectedDocument: "testdata",
		},
		{
			Expression:       `(has? .)`,
			Expected:         ast.Bool(true),
			Document:         nil,
			ExpectedDocument: nil,
		},
		{
			Expression:       `(has? .foo)`,
			Expected:         ast.Bool(false),
			Document:         "testdata",
			ExpectedDocument: "testdata",
		},
		{
			Expression:       `(has? .)`,
			Expected:         ast.Bool(true),
			Document:         64,
			ExpectedDocument: int64(64),
		},
		{
			Expression:       `(has? .)`,
			Expected:         ast.Bool(true),
			Document:         ast.String("foo"),
			ExpectedDocument: "foo",
		},

		// follow a path expression on a variable

		{
			Expression: `(has? $myvar)`,
			Invalid:    true,
		},
		{
			// missing path expression (TODO: should this be valid?)
			Expression: `(has? $myvar)`,
			Invalid:    true,
			Variables:  testVariables,
		},
		{
			Expression: `(has? $myvar.foo)`,
			Expected:   ast.Bool(false),
			Variables:  testVariables,
		},
		{
			Expression: `(has? $myvar[0])`,
			Expected:   ast.Bool(false),
			Variables:  testVariables,
		},
		{
			Expression: `(has? $obj.aString)`,
			Expected:   ast.Bool(true),
			Variables:  testVariables,
		},
		{
			Expression: `(has? $obj.aList[1])`,
			Expected:   ast.Bool(true),
			Variables:  testVariables,
		},
		{
			Expression: `(has? $vec[1])`,
			Expected:   ast.Bool(true),
			Variables:  testVariables,
		},
		{
			Expression: `(has? $astVec[0])`,
			Expected:   ast.Bool(true),
			Variables:  testVariables,
		},

		// follow a path expression on a vector node

		{
			Expression: `(has? [1 2 3][1])`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(has? [1 2 3][4])`,
			Expected:   ast.Bool(false),
		},

		// follow a path expression on an object node

		{
			Expression: `(has? {foo "bar"}.foo)`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(has? {foo "bar"}.bar)`,
			Expected:   ast.Bool(false),
		},

		// follow a path expression on a tuple node
		// (don't even need "set!" here)

		{
			Expression: `(has? (set $foo {foo "bar"}).foo)`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(has? (set $foo {foo "bar"}).bar)`,
			Expected:   ast.Bool(false),
		},
		{
			Expression: `(has? (set $foo [1])[0])`,
			Expected:   ast.Bool(true),
		},
		{
			Expression: `(has? (set $foo {foo "bar"})[0])`,
			Expected:   ast.Bool(false),
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}
