// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package test

import (
	"strings"
	"testing"

	"go.xrstf.de/rudi/pkg/lang/parser"
	"go.xrstf.de/rudi/pkg/printer"
)

func TestRudiRenderer(t *testing.T) {
	testcases := []struct {
		input  string
		output string
	}{
		{
			input:  `null`,
			output: `null`,
		},
		{
			input:  `.foo`,
			output: `.foo`,
		},
		{
			input:  `.foo[ 1 ]`,
			output: `.foo[1]`,
		},
		{
			input:  `.foo.bar[ ( + 2 23 ) ]`,
			output: `.foo.bar[(+ 2 23)]`,
		},
		{
			input:  `.["1"]`,
			output: `.["1"]`,
		},
		{
			input:  `.[1]`,
			output: `.[1]`,
		},
		{
			input:  `.[1][2].foo[3]`,
			output: `.[1][2].foo[3]`,
		},
		{
			input:  `.["foo"]`,
			output: `.foo`,
		},
		{
			input:  `.[" foo "]`,
			output: `.[" foo "]`,
		},
		{
			input:  `.["fo\"o"]`,
			output: `.["fo\"o"]`,
		},
		{
			// series of 3 statements
			input:  `1 (foo) [true]`,
			output: `1 (foo) [true]`,
		},
		{
			input:  `1`,
			output: `1`,
		},
		{
			input:  `true`,
			output: `true`,
		},
		{
			input:  `false`,
			output: `false`,
		},
		{
			input:  `[]`,
			output: `[]`,
		},
		{
			input:  `[ 1, 2  3]`,
			output: `[1 2 3]`,
		},
		{
			input:  `{}`,
			output: `{}`,
		},
		{
			input:  `{ key "value"  }`,
			output: `{key "value"}`,
		},
		{
			input:  `{ "key" 2  }`,
			output: `{key 2}`,
		},
		{
			input:  `(  foo )`,
			output: `(foo)`,
		},
		{
			input:  `(foo! bar 1 [true] "false")`,
			output: `(foo! bar 1 [true] "false")`,
		},
		{
			input:  `$foo`,
			output: `$foo`,
		},
		{
			input:  `$foo.bar`,
			output: `$foo.bar`,
		},
		{
			input:  `$foo[1]`,
			output: `$foo[1]`,
		},
		{
			input:  `$foo[1].foo`,
			output: `$foo[1].foo`,
		},
		{
			input:  `(foo).bar[1].foo`,
			output: `(foo).bar[1].foo`,
		},
		{
			input:  `(foo)[1]`,
			output: `(foo)[1]`,
		},
		{
			input:  `( map $foo.bar[ ( add 1.2)] [a  b] ( foo! ))`,
			output: `(map $foo.bar[(add 1.2)] [a b] (foo!))`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {
			prog, err := parser.ParseReader("test.go", strings.NewReader(tc.input))
			if err != nil {
				t.Fatalf("Program could not be parsed: %v", err)
			}

			var buf strings.Builder
			renderer := printer.NewRudiPrinter(&buf)

			if err = renderer.Print(prog); err != nil {
				t.Fatalf("Failed to render AST: %v", err)
			}

			rendered := buf.String()

			if rendered != tc.output {
				t.Fatalf("Expected %q, got %q", tc.output, rendered)
			}
		})
	}
}
