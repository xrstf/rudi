// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package rudifunc_test

import (
	"testing"

	"go.xrstf.de/rudi/pkg/builtin/core"
	"go.xrstf.de/rudi/pkg/builtin/math"
	"go.xrstf.de/rudi/pkg/builtin/rudifunc"
	"go.xrstf.de/rudi/pkg/builtin/strings"
	"go.xrstf.de/rudi/pkg/runtime/types"
	"go.xrstf.de/rudi/pkg/testutil"
)

func TestUserDefinedFunctions(t *testing.T) {
	// testDoc := map[string]any{
	// 	"int":    int64(4),
	// 	"float":  float64(1.2),
	// 	"bool":   true,
	// 	"string": "foo",
	// 	"null":   nil,
	// 	"vector": []any{int64(1)},
	// 	"object": map[string]any{
	// 		"key": "value",
	// 	},
	// }

	testcases := []testutil.Testcase{
		// syntax checks

		{
			Expression: `(func!)`,
			Invalid:    true,
		},
		{
			Expression: `(func! foo)`,
			Invalid:    true,
		},
		{
			Expression: `(func! foo bar)`,
			Invalid:    true,
		},
		{
			Expression: `(func! foo [])`,
			Invalid:    true,
		},
		{
			Expression: `(func! foo [1] (bar))`,
			Invalid:    true,
		},

		// defining functions, but not calling them

		{
			Expression: `(func! foo [] (bar))`,
			Expected:   nil,
		},
		{
			Expression: `(func! foo [] (bar) (bar) (bar))`,
			Expected:   nil,
		},
		{
			Expression: `(func! foo [a b c] (bar))`,
			Expected:   nil,
		},

		// functions can be redefined

		{
			Expression: `(func! foo [] (bar)) (func! foo [] (other))`,
			Expected:   nil,
		},

		// functions can be defined at any place

		{
			Expression: `(if true (func! foo [] 12)) (foo)`,
			Expected:   int64(12),
		},
		{
			Expression: `(if false (func! foo [] 12)) (foo)`,
			Invalid:    true,
		},

		// calling functions

		{
			Expression: `(func! foo [] 12) (foo)`,
			Expected:   int64(12),
		},
		{
			Expression: `(func! foo [] (append "foo" "bar")) (foo)`,
			Expected:   "foobar",
		},
		{
			Expression: `(func! foo [] "foo" 12) (foo)`,
			Expected:   int64(12),
		},

		// .. with arguments

		{
			Expression: `(func! foo [a] (+ $a 1)) (foo)`,
			Invalid:    true,
		},
		{
			Expression: `(func! foo [a] (+ $a 1)) (foo 1)`,
			Expected:   int64(2),
		},
		{
			Expression: `(func! foo [a] (+ $a 1)) (foo 1 2)`,
			Invalid:    true,
		},
		{
			Expression: `(func! foo [a b] (+ $a $b)) (foo 1 2)`,
			Expected:   int64(3),
		},

		// functions form a singular scope

		{
			Expression: `(func! foo [] (set! $a 1) (+ $a 1)) (foo)`,
			Expected:   int64(2),
		},

		// variables from within a function do not leak outside

		{
			Expression: `(func! foo [] (set! $a 1)) (foo) $a`,
			Invalid:    true,
		},

		// arguments only live inside their function as well

		{
			Expression: `(func! foo [arg] (+ $arg 1)) (foo 1)`,
			Expected:   int64(2),
		},
		{
			Expression: `(func! foo [arg] (+ $arg 1)) (foo 1) $arg`,
			Invalid:    true,
		},
		{
			Expression: `(func! foo [arg] (set! $arg 2) (+ $arg 1)) (foo 1)`,
			Expected:   int64(3),
		},
		{
			Expression: `(func! foo [arg] (set! $arg 2) (+ $arg 1)) (foo 1) $arg`,
			Invalid:    true,
		},

		// functions only see their arguments, no other variables

		{
			Expression: `(func! foo [] (+ $a 1)) (set! $a 1) (foo)`,
			Invalid:    true,
		},
		{
			Expression: `(set! $a 1) (func! foo [] (+ $a 1)) (foo)`,
			Invalid:    true,
		},
		{
			Expression: `(func! foo [a] (+ $a 1)) (set! $a 1) (foo $a)`,
			Expected:   int64(2),
		},
		{
			Expression: `(set! $a 1) (func! foo [a] (+ $a 1)) (foo $a)`,
			Expected:   int64(2),
		},
		{
			Expression: `(set! $a 1) (func! foo [b] (+ $b $a)) (foo $a)`,
			Invalid:    true,
		},

		// cannot call before it's defined

		{
			Expression: `(foo) (func! foo [] 12)`,
			Invalid:    true,
		},
	}

	funcs := types.NewFunctions()
	funcs.Add(rudifunc.Functions)
	funcs.Add(core.Functions)
	funcs.Add(math.Functions)
	funcs.Add(strings.Functions)

	for _, testcase := range testcases {
		testcase.Functions = funcs
		t.Run(testcase.String(), testcase.Run)
	}
}
