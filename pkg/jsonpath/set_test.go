// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type customObjWriter struct {
	Value any
}

var _ ObjectWriter = &customObjWriter{}

func (w customObjWriter) GetObjectKey(name string) (any, error) {
	if name == "value" {
		return w.Value, nil
	}

	return nil, fmt.Errorf("cannot get property %q", name)
}

func (w *customObjWriter) SetObjectKey(name string, value any) (any, error) {
	if name == "value" {
		w.Value = value
		return w, nil
	}

	return nil, fmt.Errorf("cannot set property %q", name)
}

type customVecWriter struct {
	Magic int
	Value any
}

var _ VectorWriter = &customVecWriter{}

func (w customVecWriter) GetVectorItem(index int) (any, error) {
	if index == w.Magic {
		return w.Value, nil
	}

	return nil, fmt.Errorf("cannot get index %d", index)
}

func (w *customVecWriter) SetVectorItem(index int, value any) (any, error) {
	if index == w.Magic {
		w.Value = value
		return w, nil
	}

	return nil, fmt.Errorf("index %d out of bounds", index)
}

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
			name:    "cannot set anything in types that do not implement the Writer interfaces",
			dest:    unknownType{},
			path:    Path{"foo"},
			invalid: true,
		},
		{
			name:    "cannot set anything in types that do not implement the Writer interfaces",
			dest:    unknownType{},
			path:    Path{0},
			invalid: true,
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

		// custom types

		{
			name: "can set in custom object writer",
			dest: &customObjWriter{
				Value: "old",
			},
			path:     Path{"value"},
			newValue: "new-value",
			expected: &customObjWriter{
				Value: "new-value",
			},
		},
		{
			name: "can set deeper in custom object writer",
			dest: &customObjWriter{
				Value: map[string]any{
					"foo": "old",
				},
			},
			path:     Path{"value", "foo"},
			newValue: "new-value",
			expected: &customObjWriter{
				Value: map[string]any{
					"foo": "new-value",
				},
			},
		},
		{
			name: "vector steps on custom object getter should fail",
			dest: &customObjWriter{
				Value: "old",
			},
			path:     Path{0},
			newValue: "new-value",
			invalid:  true,
		},

		{
			name: "can set in custom vector writer",
			dest: &customVecWriter{
				Magic: 7,
				Value: "old",
			},
			path:     Path{7},
			newValue: "new-value",
			expected: &customVecWriter{
				Magic: 7,
				Value: "new-value",
			},
		},
		{
			name: "can set deeper in custom vector writer",
			dest: &customVecWriter{
				Magic: 7,
				Value: map[string]any{
					"foo": "old",
				},
			},
			path:     Path{7, "foo"},
			newValue: "new-value",
			expected: &customVecWriter{
				Magic: 7,
				Value: map[string]any{
					"foo": "new-value",
				},
			},
		},
		{
			name: "object steps on custom vector getter should fail",
			dest: &customVecWriter{
				Magic: 7,
				Value: "old",
			},
			path:     Path{"foo"},
			newValue: "new-value",
			invalid:  true,
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
