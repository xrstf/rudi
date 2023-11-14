// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"log"
	"strings"
	"testing"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/debug"
	"go.xrstf.de/otto/pkg/lang/parser"
)

func TestParseProgram(t *testing.T) {
	testcases := []struct {
		input    string
		expected string
		invalid  bool
	}{
		/////////////////////////////////////////////////////
		// basics

		{
			input:   `(`,
			invalid: true,
		},
		{
			input:   `(foo`,
			invalid: true,
		},
		{
			input:   `foo)`,
			invalid: true,
		},
		{
			input:   `()`,
			invalid: true,
		},
		{
			input:    `(add)`,
			expected: `(tuple (identifier add))`,
		},
		{
			input:    `(+ - * /)`,
			expected: `(tuple (identifier +) (identifier -) (identifier *) (identifier /))`,
		},
		{
			input:    `(func+ -foo bar!)`,
			expected: `(tuple (identifier func+) (identifier -foo) (identifier bar!))`,
		},
		{
			input:    `(1)`,
			expected: `(tuple (number 1))`,
		},
		{
			input:    `("foo")`,
			expected: `(tuple (string "foo"))`,
		},
		{
			input:    `("f\\o\"o")`,
			expected: `(tuple (string "f\\o\"o"))`,
		},
		{
			input:   `("fo\o")`,
			invalid: true,
		},
		{
			input:   `("fo\")`,
			invalid: true,
		},
		{
			input:    `(null true false)`,
			expected: `(tuple (null) (bool true) (bool false))`,
		},
		{
			input:    `(((foo)))`,
			expected: `(tuple (tuple (tuple (identifier foo))))`,
		},

		/////////////////////////////////////////////////////
		// vectors

		{
			input:    `([])`,
			expected: `(tuple [])`,
		},
		{
			input:    `([1])`,
			expected: `(tuple [(number 1)])`,
		},
		{
			input:    `([  1 2  ])`,
			expected: `(tuple [(number 1) (number 2)])`,
		},
		{
			input:    `([  1, 2  ])`,
			expected: `(tuple [(number 1) (number 2)])`,
		},
		{
			input:    `([  1, 2,3 ,4  ])`,
			expected: `(tuple [(number 1) (number 2) (number 3) (number 4)])`,
		},

		/////////////////////////////////////////////////////
		// objects

		{
			input:    `({})`,
			expected: `(tuple {})`,
		},
		{
			input:   `({foo})`,
			invalid: true,
		},
		{
			input:    `({foo bar})`,
			expected: `(tuple {(identifier foo) (identifier bar)})`,
		},
		{
			input:    `({ foo bar (bar) 42 })`,
			expected: `(tuple {(identifier foo) (identifier bar) (tuple (identifier bar)) (number 42)})`,
		},

		/////////////////////////////////////////////////////
		// variables

		{
			input:    `($foo)`,
			expected: `(tuple (symbol (var foo)))`,
		},
		{
			input:    `(+ $foo)`,
			expected: `(tuple (identifier +) (symbol (var foo)))`,
		},
		{
			input:    `(+ $foo.bar)`,
			expected: `(tuple (identifier +) (symbol (var foo) (path [(identifier bar)])))`,
		},
		{
			input:    `(+ $foo[0])`,
			expected: `(tuple (identifier +) (symbol (var foo) (path [(number 0)])))`,
		},
		{
			input:    `(+ $foo[(foo)])`,
			expected: `(tuple (identifier +) (symbol (var foo) (path [(tuple (identifier foo))])))`,
		},
		{
			input:   `(+ $)`,
			invalid: true,
		},
		{
			input:    `(+ .bar)`,
			expected: `(tuple (identifier +) (symbol (path [(identifier bar)])))`,
		},
		{
			input:    `(+ .bar.foo)`,
			expected: `(tuple (identifier +) (symbol (path [(identifier bar) (identifier foo)])))`,
		},
		{
			input:    `(+ .bar[0])`,
			expected: `(tuple (identifier +) (symbol (path [(identifier bar) (number 0)])))`,
		},
		{
			input:    `(+ .bar[0].bar)`,
			expected: `(tuple (identifier +) (symbol (path [(identifier bar) (number 0) (identifier bar)])))`,
		},
		{
			input:    `(+ .bar[0][1])`,
			expected: `(tuple (identifier +) (symbol (path [(identifier bar) (number 0) (number 1)])))`,
		},
		{
			input:   `(+ .bar[0].[1])`,
			invalid: true,
		},

		/////////////////////////////////////////////////////
		// dot (document)

		{
			input:    `(+ . 42)`,
			expected: `(tuple (identifier +) (symbol (path [])) (number 42))`,
		},
		{
			input:    `(. 42)`,
			expected: `(tuple (symbol (path [])) (number 42))`,
		},
		{
			input:    `(+ .)`,
			expected: `(tuple (identifier +) (symbol (path [])))`,
		},
		{
			input:    `(+ .[0])`,
			expected: `(tuple (identifier +) (symbol (path [(number 0)])))`,
		},
		{
			input:   `(..)`,
			invalid: true,
		},
		{
			input:   `(.foo.)`,
			invalid: true,
		},
		{
			input:   `(.[0].)`,
			invalid: true,
		},
		{
			input:   `(foo .[0]. bar)`,
			invalid: true,
		},
		{
			input:   `( foo ..)`,
			invalid: true,
		},
		{
			input:   `(.. foo)`,
			invalid: true,
		},
		{
			input:   `(."foo")`,
			invalid: true,
		},

		// These would be nice, but the current grammar does not support it; use a temp variable
		// instead; maybe add support for this later? :)

		{
			input:   `(add (bar).foo)`,
			invalid: true,
		},
		{
			input:   `(add (bar)[1])`,
			invalid: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.input, func(t *testing.T) {
			prog := strings.NewReader(testcase.input)

			got, err := parser.ParseReader("test.go", prog)
			if err != nil {
				if !testcase.invalid {
					t.Fatalf("Failed to parse %s: %v", testcase.input, err)
				}

				return
			}

			if testcase.invalid {
				t.Fatalf("Should not have been able to parse %s, but got: %v", testcase.input, got)
			}

			program, ok := got.(ast.Program)
			if !ok {
				log.Fatalf("Parsed result is not a ast.Program, but %T", got)
			}

			var output strings.Builder
			if err := debug.DumpSingleline(&program, &output); err != nil {
				t.Fatalf("Failed to dump AST for %s: %v", testcase.input, err)
			}

			result := strings.TrimSpace(output.String())
			if result != testcase.expected {
				t.Fatalf("Expected %s -- got %s", testcase.expected, result)
			}
		})
	}
}
