// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package pathexpr

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSet(t *testing.T) {
	testcases := []struct {
		name     string
		dest     any
		path     Path
		newValue any
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
			name:     "scalar root value can simply be changed",
			dest:     nil,
			path:     Path{},
			newValue: "foo",
			expected: "foo",
		},
		{
			name:     "scalar root value can simply be changed",
			dest:     "hello world",
			path:     Path{},
			newValue: "new value",
			expected: "new value",
		},
		{
			name:     "nils can be turned into objects",
			dest:     nil,
			path:     Path{"foo"},
			newValue: "bar",
			expected: map[string]any{"foo": "bar"},
		},
		{
			name:     "nils cannot turn into slices",
			dest:     nil,
			path:     Path{0},
			newValue: "bar",
			invalid:  true,
		},
		{
			name:     "only nils can type shift",
			dest:     "a string",
			path:     Path{"foo"},
			newValue: "bar",
			invalid:  true,
		},
		{
			name:     "only nils can type shift",
			dest:     42,
			path:     Path{"foo"},
			newValue: "bar",
			invalid:  true,
		},
		{
			name:     "root object key can be updated",
			dest:     map[string]any{"foo": "bar"},
			path:     Path{"foo"},
			newValue: "new-value",
			expected: map[string]any{"foo": "new-value"},
		},
		{
			name:     "root object key can be added",
			dest:     map[string]any{"foo": "bar"},
			path:     Path{"test"},
			newValue: "new-value",
			expected: map[string]any{"foo": "bar", "test": "new-value"},
		},
		{
			name:     "root slice can be updated",
			dest:     []any{int64(1), 2, int64(3)},
			path:     Path{1},
			newValue: "new-value",
			expected: []any{int64(1), "new-value", int64(3)},
		},
		{
			name:     "handle out of bounds",
			dest:     []any{1, 2, 3},
			path:     Path{-1},
			newValue: "new-value",
			invalid:  true,
		},
		{
			name:     "handle out of bounds",
			dest:     []any{1, 2, 3},
			path:     Path{3},
			newValue: "new-value",
			invalid:  true,
		},
		{
			name:     "sub object key can be updated",
			dest:     map[string]any{"foo": "bar", "deeper": map[string]any{"deep": "value", "other": "value"}},
			path:     Path{"deeper", "deep"},
			newValue: "new-value",
			expected: map[string]any{"foo": "bar", "deeper": map[string]any{"deep": "new-value", "other": "value"}},
		},
		{
			name:     "sub slice key can be updated",
			dest:     map[string]any{"foo": "bar", "deeper": []any{1, 2, map[string]any{"deep": "value"}}},
			path:     Path{"deeper", 2, "deep"},
			newValue: "new-value",
			expected: map[string]any{"foo": "bar", "deeper": []any{1, 2, map[string]any{"deep": "new-value"}}},
		},
		{
			name:     "sub slice key can be updated",
			dest:     map[string]any{"foo": "bar", "deeper": []any{1, 2, map[string]any{"deep": "value"}}},
			path:     Path{"deeper", "whoops"},
			newValue: "new-value",
			invalid:  true,
		},
		{
			name:     "can change value types",
			dest:     map[string]any{"foo": "bar", "deeper": []any{1, 2, map[string]any{"deep": "value"}}},
			path:     Path{"deeper", 2},
			newValue: "new-value",
			expected: map[string]any{"foo": "bar", "deeper": []any{1, 2, "new-value"}},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Set(tc.dest, tc.path, tc.newValue)
			if err != nil {
				if !tc.invalid {
					t.Fatalf("Failed to run: %v", err)
				}

				return
			}

			if tc.invalid {
				t.Fatalf("Should not have been able to set value, but got: %v (%T)", result, result)
			}

			if !cmp.Equal(tc.expected, result) {
				t.Fatalf("Expected %v (%T), but got %v (%T)", tc.expected, tc.expected, result, result)
			}
		})
	}
}
