// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type unknownType struct{}

type customObjGetter struct {
	value any
}

var _ ObjectReader = customObjGetter{}

func (g customObjGetter) GetObjectKey(name string) (any, error) {
	if name == "value" {
		return g.value, nil
	}

	return nil, fmt.Errorf("cannot get property %q", name)
}

type customVecGetter struct {
	magic int
	value any
}

var _ VectorReader = customVecGetter{}

func (g customVecGetter) GetVectorItem(index int) (any, error) {
	if index == g.magic {
		return g.value, nil
	}

	return nil, fmt.Errorf("cannot get index %d", index)
}

func TestGetSingle(t *testing.T) {
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
			path:    Path{KeyStep("foo")},
			invalid: true,
		},
		{
			value:   nil,
			path:    Path{0},
			invalid: true,
		},
		{
			value:   "scalar",
			path:    Path{KeyStep("foo")},
			invalid: true,
		},
		{
			value:   func() {},
			path:    Path{KeyStep("foo")},
			invalid: true,
		},
		{
			value:   unknownType{},
			path:    Path{KeyStep("foo")},
			invalid: true,
		},

		// simply object access without recursion

		{
			value:    map[string]any{"foo": "bar"},
			path:     Path{KeyStep("foo")},
			expected: "bar",
		},
		{
			value:   map[string]any{"foo": "bar"},
			path:    Path{KeyStep("nonexisting")},
			invalid: true,
		},
		{
			value:   map[string]any{"foo": "bar"},
			path:    Path{IndexStep(0)},
			invalid: true,
		},

		// simply slice access without recursion

		{
			value:    []any{"foo", "bar", "baz"},
			path:     Path{IndexStep(0)},
			expected: "foo",
		},
		{
			value:    []any{"foo", "bar", "baz"},
			path:     Path{IndexStep(2)},
			expected: "baz",
		},
		{
			value:   []any{"foo", "bar", "baz"},
			path:    Path{IndexStep(-1)},
			invalid: true,
		},
		{
			value:   []any{"foo", "bar", "baz"},
			path:    Path{IndexStep(3)},
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
			path:     Path{KeyStep("foo"), IndexStep(1)},
			expected: "b",
		},
		{
			value:   map[string]any{"foo": []any{"a", "b"}},
			path:    Path{KeyStep("foo"), IndexStep(2)},
			invalid: true,
		},
		{
			value:   map[string]any{"foo": []any{"a", "b"}},
			path:    Path{KeyStep("foo"), IndexStep(1), KeyStep("deep")},
			invalid: true,
		},
		{
			value:    map[string]any{"foo": []any{"a", "b", map[string]any{"deep": "value"}}},
			path:     Path{KeyStep("foo"), IndexStep(2), KeyStep("deep")},
			expected: "value",
		},
		{
			value:   map[string]any{"foo": []any{"a", "b", map[string]any{"deep": "value"}}},
			path:    Path{KeyStep("foo"), IndexStep(2), KeyStep("missing")},
			invalid: true,
		},

		// descend into custom types

		// {
		// 	value:    customObjGetter{value: map[string]any{"foo": "bar"}},
		// 	path:     Path{KeyStep("value"), KeyStep("foo")},
		// 	expected: "bar",
		// },
		// {
		// 	value:   customObjGetter{value: map[string]any{"foo": "bar"}},
		// 	path:    Path{KeyStep("unknown")},
		// 	invalid: true,
		// },
		// {
		// 	value:   customObjGetter{value: nil},
		// 	path:    Path{IndexStep(0)},
		// 	invalid: true,
		// },

		// {
		// 	value:    customVecGetter{magic: 7, value: map[string]any{"foo": "bar"}},
		// 	path:     Path{IndexStep(7), KeyStep("foo")},
		// 	expected: "bar",
		// },
		// {
		// 	value:   customVecGetter{magic: 7, value: map[string]any{"foo": "bar"}},
		// 	path:    Path{IndexStep(2)},
		// 	invalid: true,
		// },
		// {
		// 	value:   customVecGetter{magic: 7, value: nil},
		// 	path:    Path{KeyStep("objectstep")},
		// 	invalid: true,
		// },
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

type keySelector []string

func (ks keySelector) Keep(key any, _ any) (bool, error) {
	keyString, ok := key.(string)
	if !ok {
		return false, errors.New("keySelector is meant to only work on objects")
	}

	for _, k := range ks {
		if k == keyString {
			return true, nil
		}
	}

	return false, nil
}

type indexSelector []int

func (is indexSelector) Keep(key any, _ any) (bool, error) {
	index, ok := key.(int)
	if !ok {
		return false, errors.New("indexSelector is meant to only work on vectors")
	}

	for _, i := range is {
		if i == index {
			return true, nil
		}
	}

	return false, nil
}

func TestGetFiltered(t *testing.T) {
	testcases := []struct {
		value    any
		path     Path
		expected any
		invalid  bool
	}{
		{
			value:    nil,
			path:     Path{keySelector{"foo"}},
			expected: []any{},
		},

		{
			value:    "a string",
			path:     Path{keySelector{"foo"}},
			expected: []any{},
		},

		{
			value:    "a string",
			path:     Path{indexSelector{2}},
			expected: []any{},
		},

		{
			value: map[string]any{
				"foo":   "a",
				"hello": "b",
				"bla":   "c",
			},
			path:     Path{keySelector{"foo"}},
			expected: []any{"a"},
		},

		// Result is ordered by key to ensure consistency.
		// (ergo, bla comes before foo).

		{
			value: map[string]any{
				"foo":   "a",
				"hello": "b",
				"bla":   "c",
			},
			path:     Path{keySelector{"foo", "bla"}},
			expected: []any{"c", "a"},
		},

		{
			value: map[string]any{
				"foo":   []any{1, 2, 3},
				"hello": []any{4, 5, 6},
				"bla":   []any{7, 8, 9},
			},
			path:     Path{keySelector{"foo", "bla"}, IndexStep(1)},
			expected: []any{8, 2},
		},

		{
			value: map[string]any{
				"foo": map[string]any{
					"a": 1,
				},
				"hello": map[string]any{
					"a": 2,
				},
				"bla": map[string]any{},
			},
			path:     Path{keySelector{"foo", "bla"}, KeyStep("a")},
			expected: []any{1},
		},

		{
			value: map[string]any{
				"foo": map[string]any{
					"a": 1,
				},
				"hello": map[string]any{
					"a": 2,
				},
				"bla": []any{1, 2, 3},
			},
			path:    Path{keySelector{"foo", "bla"}, KeyStep("a")},
			invalid: true,
		},

		{
			value: map[string]any{
				"foo":   "a",
				"hello": "b",
				"bla":   "c",
			},
			path:     Path{keySelector{"nonexisting"}},
			expected: []any{},
		},

		// In dynamic paths, it's okay to continue to traverse in an empty
		// result.

		{
			value: map[string]any{
				"foo":   "a",
				"hello": "b",
				"bla":   "c",
			},
			path:     Path{keySelector{"nonexisting"}, IndexStep(1), KeyStep("foo")},
			expected: []any{},
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
