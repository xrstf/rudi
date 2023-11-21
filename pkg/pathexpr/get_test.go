// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package pathexpr

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func TestGet(t *testing.T) {
	testcases := []struct {
		value    any
		path     Path
		expected any
		invalid  bool
	}{
		// basics

		{
			value:    nil,
			path:     Path{},
			expected: nil,
		},
		{
			value:    "hello world",
			path:     Path{},
			expected: "hello world",
		},
		{
			value:   nil,
			path:    Path{"foo"},
			invalid: true,
		},
		{
			value:   nil,
			path:    Path{0},
			invalid: true,
		},
		{
			value:   "scalar",
			path:    Path{"foo"},
			invalid: true,
		},
		{
			value:   func() {},
			path:    Path{"foo"},
			invalid: true,
		},

		// simply object access without recursion

		{
			value:    map[string]any{"foo": "bar"},
			path:     Path{"foo"},
			expected: "bar",
		},
		{
			value:   map[string]any{"foo": "bar"},
			path:    Path{"nonexisting"},
			invalid: true,
		},
		{
			value:   map[string]any{"foo": "bar"},
			path:    Path{0},
			invalid: true,
		},
		{
			value:    ast.Object{Data: map[string]any{"foo": "bar"}},
			path:     Path{"foo"},
			expected: "bar",
		},

		// simply slice access without recursion

		{
			value:    []any{"foo", "bar", "baz"},
			path:     Path{0},
			expected: "foo",
		},
		{
			value:    ast.Vector{Data: []any{"foo", "bar", "baz"}},
			path:     Path{0},
			expected: "foo",
		},
		{
			value:    []any{"foo", "bar", "baz"},
			path:     Path{int32(0)},
			expected: "foo",
		},
		{
			value:    []any{"foo", "bar", "baz"},
			path:     Path{int64(0)},
			expected: "foo",
		},
		{
			value:    []any{"foo", "bar", "baz"},
			path:     Path{2},
			expected: "baz",
		},
		{
			value:   []any{"foo", "bar", "baz"},
			path:    Path{-1},
			invalid: true,
		},
		{
			value:   []any{"foo", "bar", "baz"},
			path:    Path{3},
			invalid: true,
		},
		{
			value:   []any{"foo", "bar", "baz"},
			path:    Path{"string"},
			invalid: true,
		},

		// descend into deeper levels

		{
			value:    map[string]any{"foo": []any{"a", "b"}},
			path:     Path{"foo", 1},
			expected: "b",
		},
		{
			value:   map[string]any{"foo": []any{"a", "b"}},
			path:    Path{"foo", 2},
			invalid: true,
		},
		{
			value:   map[string]any{"foo": []any{"a", "b"}},
			path:    Path{"foo", 1, "deep"},
			invalid: true,
		},
		{
			value:    map[string]any{"foo": []any{"a", "b", map[string]any{"deep": "value"}}},
			path:     Path{"foo", 2, "deep"},
			expected: "value",
		},
		{
			value:   map[string]any{"foo": []any{"a", "b", map[string]any{"deep": "value"}}},
			path:    Path{"foo", 2, "missing"},
			invalid: true,
		},
	}

	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			result, err := Get(tc.value, tc.path)
			if err != nil {
				if !tc.invalid {
					t.Fatalf("Failed to run: %v", err)
				}

				return
			}

			if tc.invalid {
				t.Fatalf("Should not have been able to get value, but got: %v (%T)", result, result)
			}

			if !cmp.Equal(tc.expected, result) {
				t.Fatalf("Expected %v (%T), but got %v (%T)", tc.expected, tc.expected, result, result)
			}
		})
	}
}
