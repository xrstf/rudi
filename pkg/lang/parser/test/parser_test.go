// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"log"
	"strings"
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/lang/parser"
	"go.xrstf.de/rudi/pkg/printer"
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
			expected: `(tuple (identifier func+) (identifier -foo) (identifier bar (bang)))`,
		},
		{
			input:    `(1)`,
			expected: `(tuple (number (int64 1)))`,
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
			expected: `(tuple (vector))`,
		},
		{
			input:    `([1])`,
			expected: `(tuple (vector (number (int64 1))))`,
		},
		{
			input:    `([  1 2  ])`,
			expected: `(tuple (vector (number (int64 1)) (number (int64 2))))`,
		},
		{
			input:    `([  1, 2  ])`,
			expected: `(tuple (vector (number (int64 1)) (number (int64 2))))`,
		},
		{
			input:    `([  1, 2,3 ,4  ])`,
			expected: `(tuple (vector (number (int64 1)) (number (int64 2)) (number (int64 3)) (number (int64 4))))`,
		},

		/////////////////////////////////////////////////////
		// objects

		{
			input:    `({})`,
			expected: `(tuple (object))`,
		},
		{
			input:   `({foo})`,
			invalid: true,
		},
		{
			input:    `({foo bar})`,
			expected: `(tuple (object ((identifier foo) (identifier bar))))`,
		},
		{
			input:    `({ foo bar (bar) 42 })`,
			expected: `(tuple (object ((identifier foo) (identifier bar)) ((tuple (identifier bar)) (number (int64 42)))))`,
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
			expected: `(tuple (identifier +) (symbol (var foo) (path [(expr (string "bar"))])))`,
		},
		{
			input:    `(+ $foo[0])`,
			expected: `(tuple (identifier +) (symbol (var foo) (path [(expr (number (int64 0)))])))`,
		},
		{
			input:    `(+ $foo[(foo)])`,
			expected: `(tuple (identifier +) (symbol (var foo) (path [(expr (tuple (identifier foo)))])))`,
		},
		{
			input:   `(+ $)`,
			invalid: true,
		},
		{
			input:    `(+ .bar)`,
			expected: `(tuple (identifier +) (symbol (path [(expr (string "bar"))])))`,
		},
		{
			input:    `(+ .bar.foo)`,
			expected: `(tuple (identifier +) (symbol (path [(expr (string "bar")) (expr (string "foo"))])))`,
		},
		{
			input:    `(+ .bar[0])`,
			expected: `(tuple (identifier +) (symbol (path [(expr (string "bar")) (expr (number (int64 0)))])))`,
		},
		{
			input:    `(+ .bar[0].bar)`,
			expected: `(tuple (identifier +) (symbol (path [(expr (string "bar")) (expr (number (int64 0))) (expr (string "bar"))])))`,
		},
		{
			input:    `(+ .bar[0][1])`,
			expected: `(tuple (identifier +) (symbol (path [(expr (string "bar")) (expr (number (int64 0))) (expr (number (int64 1)))])))`,
		},
		{
			input:   `(+ .bar[0].[1])`,
			invalid: true,
		},

		/////////////////////////////////////////////////////
		// dot (document)

		{
			input:    `(+ . 42)`,
			expected: `(tuple (identifier +) (symbol (path [])) (number (int64 42)))`,
		},
		{
			input:    `(. 42)`,
			expected: `(tuple (symbol (path [])) (number (int64 42)))`,
		},
		{
			input:    `(+ .)`,
			expected: `(tuple (identifier +) (symbol (path [])))`,
		},
		{
			input:    `(+ .[0])`,
			expected: `(tuple (identifier +) (symbol (path [(expr (number (int64 0)))])))`,
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

		/////////////////////////////////////////////////////
		// can also use a single symbol as the top level element of a program

		{
			input:    `.`,
			expected: `(symbol (path []))`,
		},
		{
			input:    ` . `,
			expected: `(symbol (path []))`,
		},
		{
			input:   `..`,
			invalid: true,
		},
		{
			input:    `. .`,
			expected: `(symbol (path [])) (symbol (path []))`,
		},
		{
			input:    `(+ 1 2) .bar`,
			expected: `(tuple (identifier +) (number (int64 1)) (number (int64 2))) (symbol (path [(expr (string "bar"))]))`,
		},
		{
			input:    `.bar (+ 1 2)`,
			expected: `(symbol (path [(expr (string "bar"))])) (tuple (identifier +) (number (int64 1)) (number (int64 2)))`,
		},
		{
			input:   `.bar.`,
			invalid: true,
		},
		{
			input:   `(bar).`,
			invalid: true,
		},
		{
			input:    `.bar`,
			expected: `(symbol (path [(expr (string "bar"))]))`,
		},
		{
			input:    `.bar .bar`,
			expected: `(symbol (path [(expr (string "bar"))])) (symbol (path [(expr (string "bar"))]))`,
		},
		{
			input:    `.bar.foo`,
			expected: `(symbol (path [(expr (string "bar")) (expr (string "foo"))]))`,
		},
		{
			input:    `.bar[0]`,
			expected: `(symbol (path [(expr (string "bar")) (expr (number (int64 0)))]))`,
		},
		{
			input:    `.[0]`,
			expected: `(symbol (path [(expr (number (int64 0)))]))`,
		},
		{
			input:    `$var`,
			expected: `(symbol (var var))`,
		},
		{
			input:    `$var[0]`,
			expected: `(symbol (var var) (path [(expr (number (int64 0)))]))`,
		},
		{
			input:    `42`,
			expected: `(number (int64 42))`,
		},
		{
			input:    `true`,
			expected: `(bool true)`,
		},
		{
			input:    `null`,
			expected: `(null)`,
		},
		{
			input:    `"foobar"`,
			expected: `(string "foobar")`,
		},
		{
			input:    `[1 2]`,
			expected: `(vector (number (int64 1)) (number (int64 2)))`,
		},
		{
			input:    `{foo "bar"}`,
			expected: `(object ((identifier foo) (string "bar")))`,
		},

		/////////////////////////////////////////////////////
		// path expressions should work on tuples, vectors and objects too;
		// some of these examples look "obviously" wrong, but are still syntactically
		// correct; things like (add 1 2).foo will blow up during evaluation.

		{
			input:    `(add 1 2).foo`,
			expected: `(tuple (identifier add) (number (int64 1)) (number (int64 2))).(path [(expr (string "foo"))])`,
		},
		{
			input:    `(add 1 2)[1]`,
			expected: `(tuple (identifier add) (number (int64 1)) (number (int64 2))).(path [(expr (number (int64 1)))])`,
		},
		{
			input:    `(add (bar).foo)`,
			expected: `(tuple (identifier add) (tuple (identifier bar)).(path [(expr (string "foo"))]))`,
		},
		{
			input:    `(add (bar).foo[1])`,
			expected: `(tuple (identifier add) (tuple (identifier bar)).(path [(expr (string "foo")) (expr (number (int64 1)))]))`,
		},
		{
			input:    `(add (bar)[1])`,
			expected: `(tuple (identifier add) (tuple (identifier bar)).(path [(expr (number (int64 1)))]))`,
		},
		{
			input:    `(add (bar)[1].sub)`,
			expected: `(tuple (identifier add) (tuple (identifier bar)).(path [(expr (number (int64 1))) (expr (string "sub"))]))`,
		},
		{
			input:    `(add [1 2 3][0].foo)`,
			expected: `(tuple (identifier add) (vector (number (int64 1)) (number (int64 2)) (number (int64 3))).(path [(expr (number (int64 0))) (expr (string "foo"))]))`,
		},
		{
			input:   `(add [1 2 3].[0])`,
			invalid: true,
		},
		{
			input:   `(add [1 2 3].foo)`,
			invalid: true,
		},
		{
			input:    `(add {foo "bar"}.foo[1])`,
			expected: `(tuple (identifier add) (object ((identifier foo) (string "bar"))).(path [(expr (string "foo")) (expr (number (int64 1)))]))`,
		},
		{
			input:   `(add {foo "bar"}[1])`,
			invalid: true,
		},
		{
			input:    `(add {foo "bar"}.foo[{bla 1}.bla[0]])`,
			expected: `(tuple (identifier add) (object ((identifier foo) (string "bar"))).(path [(expr (string "foo")) (expr (object ((identifier bla) (number (int64 1)))).(path [(expr (string "bla")) (expr (number (int64 0)))]))]))`,
		},
		{
			input:    `(add {{foo "bar"}.foo "bar"})`,
			expected: `(tuple (identifier add) (object ((object ((identifier foo) (string "bar"))).(path [(expr (string "foo"))]) (string "bar"))))`,
		},
		{
			input:    `(add {[0 "foo"][1] "bar"})`,
			expected: `(tuple (identifier add) (object ((vector (number (int64 0)) (string "foo")).(path [(expr (number (int64 1)))]) (string "bar"))))`,
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
				var output strings.Builder
				p := printer.NewAstPrinter(&output)

				if err := p.Print(got); err != nil {
					t.Errorf("Failed to dump unexpected AST for %s: %v", testcase.input, err)
				}

				t.Fatalf("Should not have been able to parse %s, but got: %v", testcase.input, strings.TrimSpace(output.String()))
			}

			program, ok := got.(ast.Program)
			if !ok {
				log.Fatalf("Parsed result is not a ast.Program, but %T", got)
			}

			var output strings.Builder
			p := printer.NewAstPrinter(&output)

			if err := p.Program(&program); err != nil {
				t.Fatalf("Failed to dump AST for %s: %v", testcase.input, err)
			}

			result := strings.TrimSpace(output.String())
			if result != testcase.expected {
				t.Fatalf("Expected %s -- got %s", testcase.expected, result)
			}
		})
	}
}
