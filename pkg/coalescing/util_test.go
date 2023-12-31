// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package coalescing

import (
	"fmt"
	"strings"
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"

	"github.com/google/go-cmp/cmp"
)

type invalidConversion int

const invalid invalidConversion = iota

type testcase struct {
	value    any
	toNull   any
	toBool   any
	toInt    any
	toFloat  any
	toNumber any
	toString any
	toVector any
	toObject any
}

func newTestcase(value, toNull, toBool, toInt, toFloat, toNumber, toString, toVector, toObject any) testcase {
	return testcase{
		value:    value,
		toNull:   toNull,
		toBool:   toBool,
		toInt:    toInt,
		toFloat:  toFloat,
		toNumber: toNumber,
		toString: toString,
		toVector: toVector,
		toObject: toObject,
	}
}

func newNum(value any) ast.Number {
	return ast.Number{Value: value}
}

func testCoalescer(t *testing.T, coalescer Coalescer, testcases []testcase) {
	t.Helper()

	testConversion(
		t,
		coalescer,
		testcases,
		"null",
		func(val any) (any, error) { return coalescer.ToNull(val) },
		func(tc testcase) any { return tc.toNull },
	)

	testConversion(
		t,
		coalescer,
		testcases,
		"bool",
		func(val any) (any, error) { return coalescer.ToBool(val) },
		func(tc testcase) any { return tc.toBool },
	)

	testConversion(
		t,
		coalescer,
		testcases,
		"int",
		func(val any) (any, error) { return coalescer.ToInt64(val) },
		func(tc testcase) any { return tc.toInt },
	)

	testConversion(
		t,
		coalescer,
		testcases,
		"float",
		func(val any) (any, error) { return coalescer.ToFloat64(val) },
		func(tc testcase) any { return tc.toFloat },
	)

	testConversion(
		t,
		coalescer,
		testcases,
		"number",
		func(val any) (any, error) { return coalescer.ToNumber(val) },
		func(tc testcase) any { return tc.toNumber },
	)

	testConversion(
		t,
		coalescer,
		testcases,
		"string",
		func(val any) (any, error) { return coalescer.ToString(val) },
		func(tc testcase) any { return tc.toString },
	)

	testConversion(
		t,
		coalescer,
		testcases,
		"vector",
		func(val any) (any, error) { return coalescer.ToVector(val) },
		func(tc testcase) any { return tc.toVector },
	)

	testConversion(
		t,
		coalescer,
		testcases,
		"object",
		func(val any) (any, error) { return coalescer.ToObject(val) },
		func(tc testcase) any { return tc.toObject },
	)
}

func testConversion(
	t *testing.T,
	coalescer Coalescer,
	testcases []testcase,
	name string,
	convert func(any) (any, error),
	getExpected func(testcase) any,
) {
	t.Helper()

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v to %s", tc.value, name), func(t *testing.T) {
			expected := getExpected(tc)
			_, expectErr := expected.(invalidConversion)

			result, err := convert(tc.value)
			if err != nil {
				if !expectErr {
					t.Fatalf("Failed to do %s(%s): %v", name, printValue(tc.value), err)
				}

				return
			}

			if expectErr {
				t.Fatalf("Should not have been able to do %s(%s), but got: %v", name, printValue(tc.value), result)
			}

			if !cmp.Equal(result, expected) {
				t.Fatalf("expected %s(%s) => %s, but got %s", name, printValue(tc.value), printValue(expected), printValue(result))
			}
		})
	}
}

func printValue(val any) string {
	t := fmt.Sprintf("%T(%#v)", val, val)
	t = strings.ReplaceAll(t, "interface {}", "any")

	return t
}
