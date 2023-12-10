// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package core

import (
	"testing"

	"go.xrstf.de/rudi/pkg/eval/types"
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
			Expression: `(if true 3)`,
			Expected:   int64(3),
		},
		{
			Expression: `(if false 3)`,
			Expected:   nil,
		},
		{
			Expression: `(if false "yes" "else")`,
			Expected:   "else",
		},
		{
			// strict coalescing allows to turn null into false
			Expression: `(if null "yes" "else")`,
			Expected:   "else",
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
			"aList":   []any{"first", 2, "third"},
			"aBool":   true,
			"anObject": map[string]any{
				"key1": true,
				"key2": nil,
				"key3": []any{9, map[string]any{"foo": "bar"}, 7},
			},
		}
	}

	testVecDocument := func() any {
		return []any{1, 2, map[string]any{"foo": "bar"}}
	}

	testVariables := func() types.Variables {
		return types.Variables{
			"myvar": 42,
			"obj":   testObjDocument(),
			"vec":   testVecDocument(),
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
			Expected:   "foo",
		},
		{
			Expression: `(set $var "foo") $var`,
			Invalid:    true,
		},
		{
			Expression: `(set $var 1)`,
			Expected:   int64(1),
		},
		// with bang it works as expected
		{
			Expression: `(set! $var 1)`,
			Expected:   int64(1),
		},
		{
			Expression: `(set! $var 1) $var`,
			Expected:   int64(1),
		},
		// can overwrite variables on the top level
		{
			Expression: `(set! $myvar 12) $myvar`,
			Variables:  testVariables(),
			Expected:   int64(12),
			ExpectedVariables: types.Variables{
				"myvar": int64(12),
			},
		},
		// can change the type
		{
			Expression: `(set! $myvar "new value") $myvar`,
			Variables:  testVariables(),
			Expected:   "new value",
		},
		{
			Expression: `(set! $obj.aList[1] "new value")`,
			Variables:  testVariables(),
			Expected:   "new value",
		},
		{
			Expression: `(set! $obj.aList[1] "new value") $obj`,
			Variables:  testVariables(),
			Expected: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", "new value", "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{9, map[string]any{"foo": "bar"}, 7},
				},
			},
		},
		// set itself does not change the first argument
		{
			Expression: `(set $myvar "new value") $myvar`,
			Variables:  testVariables(),
			Expected:   42,
		},
		{
			Expression: `(set $obj.aString "new value") $obj.aString`,
			Variables:  testVariables(),
			Expected:   "foo",
		},
		{
			Expression: `(set $obj.aList[1] "new value") $obj.aList`,
			Variables:  testVariables(),
			Expected:   []any{"first", 2, "third"},
		},
		// ...but not leak into upper scopes
		{
			Expression: `(set! $a 1) (if true (set! $a 2)) $a`,
			Expected:   int64(1),
		},
		{
			Expression: `(set! $a 1) (if true (set! $b 2)) $b`,
			Invalid:    true,
		},
		// do not accidentally set a key without creating a new context
		{
			Expression: `(set! $a {foo "bar"}) (if true (set! $a.foo "updated"))`,
			Expected:   "updated",
		},
		{
			Expression: `(set! $a {foo "bar"}) (if true (set! $a.foo "updated")) $a.foo`,
			Expected:   "bar",
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
			Expected:   "new value",
			Variables:  testVariables(),
		},
		{
			Expression: `(set! $obj.aString "new value") $obj.aString`,
			Expected:   "new value",
			Variables:  testVariables(),
		},
		// add a new sub key
		{
			Expression: `(set! $obj.newKey "new value")`,
			Expected:   "new value",
			Variables:  testVariables(),
		},
		{
			Expression: `(set! $obj.newKey "new value") $obj.newKey`,
			Expected:   "new value",
			Variables:  testVariables(),
		},
		// runtime variables
		{
			Expression: `(set! $vec [1]) (set! $vec[0] 2) $vec[0]`,
			Expected:   int64(2),
		},
		// replace the global document
		{
			Expression:       `(set! . 1) .`,
			Document:         testObjDocument(),
			Expected:         int64(1),
			ExpectedDocument: int64(1),
		},
		// update keys in the global document
		{
			Expression: `(set! .aString "new-value") .aString`,
			Document:   testObjDocument(),
			Expected:   "new-value",
			ExpectedDocument: map[string]any{
				"aString": "new-value",
				"aList":   []any{"first", 2, "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{9, map[string]any{"foo": "bar"}, 7},
				},
			},
		},
		// add new keys
		{
			Expression: `(set! .newKey "new-value") .newKey`,
			Document:   testObjDocument(),
			Expected:   "new-value",
			ExpectedDocument: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", 2, "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{9, map[string]any{"foo": "bar"}, 7},
				},
				"newKey": "new-value",
			},
		},
		// update vectors
		{
			Expression: `(set! .aList[1] "new-value") .aList[1]`,
			Document:   testObjDocument(),
			Expected:   "new-value",
			ExpectedDocument: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", "new-value", "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{9, map[string]any{"foo": "bar"}, 7},
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
			"aList":   []any{"first", 2, "third"},
			"aBool":   true,
			"anObject": map[string]any{
				"key1": true,
				"key2": nil,
				"key3": []any{9, map[string]any{"foo": "bar"}, 7},
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
			Expected:         nil,
			ExpectedDocument: map[string]any{"foo": "bar"},
		},
		{
			Expression: `(delete! .) .`,
			Document:   map[string]any{"foo": "bar"},
			Expected:   nil,
		},
		// delete does not update the target
		{
			Expression:       `(delete .) .`,
			Document:         map[string]any{"foo": "bar"},
			Expected:         map[string]any{"foo": "bar"},
			ExpectedDocument: map[string]any{"foo": "bar"},
		},
		// can remove a key
		{
			Expression:       `(delete .foo)`,
			Document:         map[string]any{"foo": "bar"},
			Expected:         map[string]any{},
			ExpectedDocument: map[string]any{"foo": "bar"},
		},
		{
			Expression:       `(delete .foo) .`,
			Document:         map[string]any{"foo": "bar"},
			Expected:         map[string]any{"foo": "bar"},
			ExpectedDocument: map[string]any{"foo": "bar"},
		},
		{
			Expression:       `(delete! .foo) .`,
			Document:         map[string]any{"foo": "bar"},
			Expected:         map[string]any{},
			ExpectedDocument: map[string]any{},
		},
		// non-existent key is okay
		{
			Expression:       `(delete .bar)`,
			Document:         map[string]any{"foo": "bar"},
			Expected:         map[string]any{"foo": "bar"},
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
			Expected:         []any{"a", "c"},
			ExpectedDocument: []any{"a", "b", "c"},
		},
		{
			Expression:       `(delete .[1]) .`,
			Document:         []any{"a", "b", "c"},
			Expected:         []any{"a", "b", "c"},
			ExpectedDocument: []any{"a", "b", "c"},
		},
		{
			Expression:       `(delete! .[1]) .`,
			Document:         []any{"a", "b", "c"},
			Expected:         []any{"a", "c"},
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
			Expected: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{9, map[string]any{"foo": "bar"}, 7},
				},
			},
			ExpectedDocument: testObjDocument(),
		},
		{
			Expression:       `(delete .aList[1]) .`,
			Document:         testObjDocument(),
			Expected:         testObjDocument(),
			ExpectedDocument: testObjDocument(),
		},
		{
			Expression: `(delete! .aList[1]) .`,
			Document:   testObjDocument(),
			Expected: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{9, map[string]any{"foo": "bar"}, 7},
				},
			},
			ExpectedDocument: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{9, map[string]any{"foo": "bar"}, 7},
				},
			},
		},
		{
			Expression: `(delete .anObject.key3[1].foo)`,
			Document:   testObjDocument(),
			Expected: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", 2, "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{9, map[string]any{}, 7},
				},
			},
			ExpectedDocument: testObjDocument(),
		},
		{
			Expression:       `(delete .anObject.key3[1].foo) .`,
			Document:         testObjDocument(),
			Expected:         testObjDocument(),
			ExpectedDocument: testObjDocument(),
		},
		{
			Expression: `(delete! .anObject.key3[1].foo) .`,
			Document:   testObjDocument(),
			Expected: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", 2, "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{9, map[string]any{}, 7},
				},
			},
			ExpectedDocument: map[string]any{
				"aString": "foo",
				"aList":   []any{"first", 2, "third"},
				"aBool":   true,
				"anObject": map[string]any{
					"key1": true,
					"key2": nil,
					"key3": []any{9, map[string]any{}, 7},
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
			Expected:   int64(3),
		},

		// test that the runtime context is inherited from one step to another
		{
			Expression: `(do (set! $var "foo") $var)`,
			Expected:   "foo",
		},
		{
			Expression: `(do (set! $var "foo") $var (set! $var "new") $var)`,
			Expected:   "new",
		},

		// test that the runtime context doesn't leak
		{
			Expression: `(set! $var "outer") (do (set! $var "inner")) $var`,
			Expected:   "outer",
		},
		{
			Expression: `(do (set! $var "inner")) $var`,
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
			Expected:   int64(3),
		},

		// coalescing should be applied

		{
			Expression: `(default false 3)`,
			Expected:   int64(3),
		},
		{
			Expression: `(default [] 3)`,
			Expected:   int64(3),
		},

		// errors are not swallowed

		{
			Expression: `(default (error "foo") 3)`,
			Invalid:    true,
		},

		{
			Expression: `(default false (error "foo"))`,
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
			Expression: `(try 2)`,
			Expected:   int64(2),
		},

		// coalescing should be not applied

		{
			Expression: `(try false)`,
			Expected:   false,
		},
		{
			Expression: `(try null)`,
			Expected:   nil,
		},
		{
			Expression: `(try null "fallback")`,
			Expected:   nil,
		},

		// swallow errors

		{
			Expression: `(try (error "foo"))`,
			Expected:   nil,
		},
		{
			Expression: `(try (error "foo") "fallback")`,
			Expected:   "fallback",
		},

		// not in the fallback though

		{
			Expression: `(try (error "foo") (error "foo"))`,
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
			Expected:   true,
		},
		{
			Expression: `(empty? true)`,
			Expected:   false,
		},
		{
			Expression: `(empty? false)`,
			Expected:   true,
		},
		{
			Expression: `(empty? 0)`,
			Expected:   true,
		},
		{
			Expression: `(empty? 0.0)`,
			Expected:   true,
		},
		{
			Expression: `(empty? 1)`,
			Expected:   false,
		},
		{
			Expression: `(empty? -1.0)`,
			Expected:   false,
		},
		{
			Expression: `(empty? [])`,
			Expected:   true,
		},
		{
			Expression: `(empty? [""])`,
			Expected:   false,
		},
		{
			Expression: `(empty? {})`,
			Expected:   true,
		},
		{
			Expression: `(empty? {foo "bar"})`,
			Expected:   false,
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
		"aList":   []any{"first", 2, "third"},
		"aBool":   true,
		"anObject": map[string]any{
			"key1": true,
			"key2": nil,
			"key3": []any{9, map[string]any{"foo": "bar"}, 77},
		},
	}

	testVecDocument := []any{1, 2, map[string]any{"foo": "bar"}}

	testVariables := types.Variables{
		// value does not matter here, but this testcase is still meant
		// to ensure the missing path is detected, not detect an unknown variable
		"myvar": 42,
		"obj":   testObjDocument,
		"vec":   testVecDocument,
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
			Expected:   true,
			Document:   nil, // the . always matches, no matter what the document is
		},
		{
			Expression:       `(has? .)`,
			Expected:         true,
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .nonexistingKey)`,
			Expected:         false,
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .[0])`,
			Expected:         false,
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .aString)`,
			Expected:         true,
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .["aString"])`,
			Expected:         true,
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			// (set) is just a dummy function that in this case simply returns "aString"
			Expression:       `(has? .[(set $foo "aString")])`,
			Expected:         true,
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .aBool)`,
			Expected:         true,
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .aList)`,
			Expected:         true,
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .aList[0])`,
			Expected:         true,
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .aList[99])`,
			Expected:         false,
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .aList.invalidObjKey)`,
			Expected:         false,
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .anObject)`,
			Expected:         true,
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .anObject[99])`,
			Expected:         false,
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .anObject.key1)`,
			Expected:         true,
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .anObject.key99)`,
			Expected:         false,
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .anObject.key3[1].foo)`,
			Expected:         true,
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},
		{
			Expression:       `(has? .anObject.key3[1].bar)`,
			Expected:         false,
			Document:         testObjDocument,
			ExpectedDocument: testObjDocument,
		},

		// global document is an array

		{
			Expression:       `(has? .[1])`,
			Expected:         true,
			Document:         testVecDocument,
			ExpectedDocument: testVecDocument,
		},
		{
			Expression:       `(has? .key)`,
			Expected:         false,
			Document:         testVecDocument,
			ExpectedDocument: testVecDocument,
		},
		{
			Expression:       `(has? .[2].foo)`,
			Expected:         true,
			Document:         testVecDocument,
			ExpectedDocument: testVecDocument,
		},

		// global document is a scalar

		{
			Expression:       `(has? .)`,
			Expected:         true,
			Document:         "testdata",
			ExpectedDocument: "testdata",
		},
		{
			Expression:       `(has? .)`,
			Expected:         true,
			Document:         nil,
			ExpectedDocument: nil,
		},
		{
			Expression:       `(has? .foo)`,
			Expected:         false,
			Document:         "testdata",
			ExpectedDocument: "testdata",
		},
		{
			Expression:       `(has? .)`,
			Expected:         true,
			Document:         64,
			ExpectedDocument: 64,
		},
		{
			Expression:       `(has? .)`,
			Expected:         true,
			Document:         "foo",
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
			Expected:   false,
			Variables:  testVariables,
		},
		{
			Expression: `(has? $myvar[0])`,
			Expected:   false,
			Variables:  testVariables,
		},
		{
			Expression: `(has? $obj.aString)`,
			Expected:   true,
			Variables:  testVariables,
		},
		{
			Expression: `(has? $obj.aList[1])`,
			Expected:   true,
			Variables:  testVariables,
		},
		{
			Expression: `(has? $vec[1])`,
			Expected:   true,
			Variables:  testVariables,
		},

		// follow a path expression on a vector node

		{
			Expression: `(has? [1 2 3][1])`,
			Expected:   true,
		},
		{
			Expression: `(has? [1 2 3][4])`,
			Expected:   false,
		},

		// follow a path expression on an object node

		{
			Expression: `(has? {foo "bar"}.foo)`,
			Expected:   true,
		},
		{
			Expression: `(has? {foo "bar"}.bar)`,
			Expected:   false,
		},

		// follow a path expression on a tuple node
		// (don't even need "set!" here)

		{
			Expression: `(has? (set $foo {foo "bar"}).foo)`,
			Expected:   true,
		},
		{
			Expression: `(has? (set $foo {foo "bar"}).bar)`,
			Expected:   false,
		},
		{
			Expression: `(has? (set $foo [1])[0])`,
			Expected:   true,
		},
		{
			Expression: `(has? (set $foo {foo "bar"})[0])`,
			Expected:   false,
		},
	}

	for _, testcase := range testcases {
		testcase.Functions = Functions
		t.Run(testcase.String(), testcase.Run)
	}
}
