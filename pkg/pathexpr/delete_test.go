// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package pathexpr

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDelete(t *testing.T) {
	testcases := []struct {
		name     string
		dest     any
		path     Path
		expected any
		invalid  bool
	}{
		{
			name:    "invalid root value",
			dest:    func() {},
			path:    Path{"foo"},
			invalid: true,
		},
		{
			name:    "invalid step",
			dest:    "value",
			path:    Path{true},
			invalid: true,
		},
		{
			name:    "invalid step",
			dest:    "value",
			path:    Path{"foo"},
			invalid: true,
		},
		{
			name:     "delete root key",
			dest:     map[string]any{"foo": "bar", "other": "value"},
			path:     Path{"other"},
			expected: map[string]any{"foo": "bar"},
		},
		{
			name:     "accept key already gone",
			dest:     map[string]any{"foo": "bar"},
			path:     Path{"other"},
			expected: map[string]any{"foo": "bar"},
		},
		{
			name:     "can result in empty objects",
			dest:     map[string]any{"foo": "bar"},
			path:     Path{"foo"},
			expected: map[string]any{},
		},
		{
			name:     "nils are values, as they make keys exists",
			dest:     map[string]any{"foo": nil},
			path:     Path{"foo"},
			expected: map[string]any{},
		},
		{
			name:     "delete deeper key",
			dest:     map[string]any{"foo": "bar", "other": map[string]any{"deeper": "value"}},
			path:     Path{"other", "deeper"},
			expected: map[string]any{"foo": "bar", "other": map[string]any{}},
		},
		{
			name:     "accept missing sub key",
			dest:     map[string]any{"foo": "bar", "other": map[string]any{"deeper": "value"}},
			path:     Path{"other", "missing"},
			expected: map[string]any{"foo": "bar", "other": map[string]any{"deeper": "value"}},
		},
		{
			name:    "path must still make sense",
			dest:    map[string]any{"foo": "bar", "other": map[string]any{"deeper": "value"}},
			path:    Path{"other", 1},
			invalid: true,
		},
		{
			name:     "object in slice",
			dest:     []any{"foo", map[string]any{"foo": "bar"}, "bar"},
			path:     Path{1, "foo"},
			expected: []any{"foo", map[string]any{}, "bar"},
		},
		{
			name:    "out of bounds",
			dest:    []any{"foo", map[string]any{"foo": "bar"}, "bar"},
			path:    Path{-1},
			invalid: true,
		},
		{
			name:    "out of bounds",
			dest:    []any{"foo", map[string]any{"foo": "bar"}, "bar"},
			path:    Path{3},
			invalid: true,
		},
		{
			name:    "out of bounds",
			dest:    []any{"foo", map[string]any{"foo": "bar"}, "bar"},
			path:    Path{-1, "list"},
			invalid: true,
		},
		{
			name:    "out of bounds",
			dest:    []any{"foo", map[string]any{"foo": "bar"}, "bar"},
			path:    Path{3, "list"},
			invalid: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Delete(tc.dest, tc.path)
			if err != nil {
				if !tc.invalid {
					t.Fatalf("Failed to run: %v", err)
				}

				return
			}

			if tc.invalid {
				t.Fatalf("Should not have been able to delete path, but got: %v (%T)", result, result)
			}

			if !cmp.Equal(tc.expected, result) {
				t.Fatalf("Expected %v (%T), but got %v (%T)", tc.expected, tc.expected, result, result)
			}
		})
	}
}
