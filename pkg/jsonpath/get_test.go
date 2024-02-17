// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type unknownType struct{}

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

func TestGetComplex(t *testing.T) {
	var nilDummyFunc func() int

	stringPtr := ptrTo("foo")
	customEmptyInterfacesValues := []CustomEmptyInterface{
		map[string]*OtherStruct{"foo": {StringField: "bar"}},
	}

	testcases := []struct {
		value    any
		path     Path
		expected any
		invalid  bool
	}{
		{
			value:   func() {},
			path:    Path{KeyStep("foo")},
			invalid: true,
		},
		{
			value:   emptyStruct{},
			path:    Path{KeyStep("foo")},
			invalid: true,
		},

		{
			value:   ExampleStruct{},
			path:    Path{KeyStep("DoesNotExist")},
			invalid: true,
		},
		{
			value:   ExampleStruct{},
			path:    Path{KeyStep("privateField")},
			invalid: true,
		},
		{
			value:    ExampleStruct{StringField: ""},
			path:     Path{KeyStep("StringField")},
			expected: "",
		},
		{
			value:   ExampleStruct{StringField: ""},
			path:    Path{KeyStep("StringField"), KeyStep("dummy")},
			invalid: true,
		},
		{
			value:    ExampleStruct{StringField: "foo"},
			path:     Path{KeyStep("StringField")},
			expected: "foo",
		},
		{
			value:    ExampleStruct{StringPointerField: nil},
			path:     Path{KeyStep("StringPointerField")},
			expected: func() *string { return nil }(),
		},
		{
			value:    ExampleStruct{StringPointerField: stringPtr},
			path:     Path{KeyStep("StringPointerField")},
			expected: stringPtr,
		},
		{
			value:    ExampleStruct{FuncField: nil},
			path:     Path{KeyStep("FuncField")},
			expected: nilDummyFunc,
		},
		// is its own testcase because cmp cannot compare functions
		// {
		// 	value:    ExampleStruct{FuncField: dummyFieldFunc},
		// 	path:     Path{"FuncField"},
		// 	expected: dummyFieldFunc,
		// },
		{
			value:   ExampleStruct{FuncField: dummyFieldFunc},
			path:    Path{KeyStep("FuncField"), KeyStep("test")},
			invalid: true,
		},
		{
			value:   ExampleStruct{FuncField: dummyFieldFunc},
			path:    Path{KeyStep("FuncField"), IndexStep(0)},
			invalid: true,
		},

		// root value is pointer

		{
			value:    &ExampleStruct{StringField: "foo"},
			path:     Path{KeyStep("StringField")},
			expected: "foo",
		},
		{
			value:   func() *ExampleStruct { return nil }(),
			path:    Path{KeyStep("StringField")},
			invalid: true,
		},

		// struct map fields

		{
			value:    ExampleStruct{},
			path:     Path{KeyStep("StringMapField")},
			expected: func() map[string]string { return nil }(),
		},
		{
			value:    ExampleStruct{StringMapField: map[string]string{}},
			path:     Path{KeyStep("StringMapField")},
			expected: map[string]string{},
		},
		{
			value:   ExampleStruct{StringMapField: map[string]string{}},
			path:    Path{KeyStep("StringMapField"), IndexStep(0)},
			invalid: true,
		},
		{
			value:   ExampleStruct{StringMapField: map[string]string{}},
			path:    Path{KeyStep("StringMapField"), KeyStep("test")},
			invalid: true,
		},
		{
			value:    ExampleStruct{StringMapField: map[string]string{"test": "value"}},
			path:     Path{KeyStep("StringMapField"), KeyStep("test")},
			expected: "value",
		},
		{
			value:    ExampleStruct{EmptyInterfaceMapField: map[string]any{"test": "value"}},
			path:     Path{KeyStep("EmptyInterfaceMapField"), KeyStep("test")},
			expected: "value",
		},
		{
			value:    ExampleStruct{StructMapField: map[string]OtherStruct{"foo": {}}},
			path:     Path{KeyStep("StructMapField"), KeyStep("foo"), KeyStep("StringPointerField")},
			expected: func() *string { return nil }(),
		},
		{
			value:    ExampleStruct{StructMapField: map[string]OtherStruct{"foo": {StringPointerField: stringPtr}}},
			path:     Path{KeyStep("StructMapField"), KeyStep("foo"), KeyStep("StringPointerField")},
			expected: stringPtr,
		},
		{
			value:    ExampleStruct{StructPointerMapField: map[string]*OtherStruct{"foo": {StringField: "bar"}}},
			path:     Path{KeyStep("StructPointerMapField"), KeyStep("foo"), KeyStep("StringField")},
			expected: "bar",
		},

		// embedded structs

		{
			value:    StructWithEmbed{StringField: "foo"},
			path:     Path{KeyStep("StringField")},
			expected: "foo",
		},
		{
			value:    StructWithEmbed{BaseStruct: BaseStruct{StringField: "foo"}},
			path:     Path{KeyStep("StringField")},
			expected: "",
		},
		{
			value:    StructWithEmbed{BaseStruct: BaseStruct{StringField: "foo"}},
			path:     Path{KeyStep("BaseStruct"), KeyStep("StringField")},
			expected: "foo",
		},
		{
			value:    StructWithEmbed{StringField: "foo", BaseStruct: BaseStruct{StringField: "bar"}},
			path:     Path{KeyStep("StringField")},
			expected: "foo",
		},
		{
			value:    StructWithEmbed{StringField: "foo", BaseStruct: BaseStruct{StringField: "bar"}},
			path:     Path{KeyStep("BaseStruct"), KeyStep("StringField")},
			expected: "bar",
		},
		{
			value:    StructWithEmbed{BaseStruct: BaseStruct{BaseStringField: "bar"}},
			path:     Path{KeyStep("BaseStringField")},
			expected: "bar",
		},

		// interface fields

		{
			value:    ExampleStruct{CustomEmptyInterfaceField: map[string]*OtherStruct{"foo": {StringField: "bar"}}},
			path:     Path{KeyStep("CustomEmptyInterfaceField"), KeyStep("foo"), KeyStep("StringField")},
			expected: "bar",
		},
		{
			value:    ExampleStruct{CustomEmptyInterfaceSlicePointerField: &customEmptyInterfacesValues},
			path:     Path{KeyStep("CustomEmptyInterfaceSlicePointerField"), IndexStep(0), KeyStep("foo"), KeyStep("StringField")},
			expected: "bar",
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
