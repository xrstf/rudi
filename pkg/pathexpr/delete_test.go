// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package pathexpr

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type customObjDeleter struct {
	Value any
}

var _ ObjectKeyDeleter = &customObjDeleter{}

func (w customObjDeleter) GetObjectKey(name string) (any, error) {
	if name == "value" {
		return w.Value, nil
	}

	return nil, fmt.Errorf("cannot get property %q", name)
}

func (w *customObjDeleter) SetObjectKey(name string, value any) (any, error) {
	if name == "value" {
		w.Value = value
		return w, nil
	}

	return nil, fmt.Errorf("cannot set property %q", name)
}

func (w *customObjDeleter) DeleteObjectKey(name string) (any, error) {
	if name == "value" {
		w.Value = nil
		return w, nil
	}

	return nil, fmt.Errorf("cannot delete property %q", name)
}

type customVecDeleter struct {
	Magic int
	Value any
}

var _ VectorItemDeleter = &customVecDeleter{}

func (w customVecDeleter) GetVectorItem(index int) (any, error) {
	if index == w.Magic {
		return w.Value, nil
	}

	return nil, fmt.Errorf("cannot get index %d", index)
}

func (w *customVecDeleter) SetVectorItem(index int, value any) (any, error) {
	if index == w.Magic {
		w.Value = value
		return w, nil
	}

	return nil, fmt.Errorf("index %d out of bounds", index)
}

func (w *customVecDeleter) DeleteVectorItem(index int) (any, error) {
	if index == w.Magic {
		w.Value = nil
		return w, nil
	}

	return nil, fmt.Errorf("index %d out of bounds", index)
}

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
			name:     "delete everything",
			dest:     map[string]any{"foo": "bar", "other": "value"},
			path:     Path{},
			expected: nil,
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
			name:     "remove item from slice",
			dest:     []any{"foo", map[string]any{"foo": "bar"}, "bar"},
			path:     Path{1},
			expected: []any{"foo", "bar"},
		},
		{
			name:     "remove deeper item from slice",
			dest:     []any{"foo", map[string]any{"foo": []any{"a", "b", "c"}}, "bar"},
			path:     Path{1, "foo", 0},
			expected: []any{"foo", map[string]any{"foo": []any{"b", "c"}}, "bar"},
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

		// custom types

		{
			name: "can delete in custom objects",
			dest: &customObjDeleter{
				Value: "old",
			},
			path: Path{"value"},
			expected: &customObjDeleter{
				Value: nil,
			},
		},
		{
			name: "can delete in custom objects",
			dest: &customObjDeleter{
				Value: map[string]any{
					"foo":   "bar",
					"hello": "world",
					"list":  []any{1, 2, 3},
				},
			},
			path: Path{"value", "list", 1},
			expected: &customObjDeleter{
				Value: map[string]any{
					"foo":   "bar",
					"hello": "world",
					"list":  []any{1, 3},
				},
			},
		},

		{
			name: "can delete in custom vectors",
			dest: &customVecDeleter{
				Magic: 7,
				Value: "old",
			},
			path: Path{7},
			expected: &customVecDeleter{
				Magic: 7,
				Value: nil,
			},
		},
		{
			name: "can delete in custom vectors",
			dest: &customVecDeleter{
				Magic: 7,
				Value: map[string]any{
					"foo":   "bar",
					"hello": "world",
					"list":  []any{1, 2, 3},
				},
			},
			path: Path{7, "list", 1},
			expected: &customVecDeleter{
				Magic: 7,
				Value: map[string]any{
					"foo":   "bar",
					"hello": "world",
					"list":  []any{1, 3},
				},
			},
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
