// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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

func TestGet(t *testing.T) {
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
		{
			value:   emptyStruct{},
			path:    Path{"foo"},
			invalid: true,
		},

		// simply map access without recursion

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

		// simply slice access without recursion

		{
			value:    []any{"foo", "bar", "baz"},
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

		// custom object getter

		{
			value:    customObjGetter{value: map[string]any{"foo": "bar"}},
			path:     Path{"value", "foo"},
			expected: "bar",
		},
		{
			value:   customObjGetter{value: map[string]any{"foo": "bar"}},
			path:    Path{"unknown"},
			invalid: true,
		},
		{
			value:   customObjGetter{value: nil},
			path:    Path{0},
			invalid: true,
		},

		// custom vector getter

		{
			value:    customVecGetter{magic: 7, value: map[string]any{"foo": "bar"}},
			path:     Path{7, "foo"},
			expected: "bar",
		},
		{
			value:   customVecGetter{magic: 7, value: map[string]any{"foo": "bar"}},
			path:    Path{2},
			invalid: true,
		},
		{
			value:   customVecGetter{magic: 7, value: nil},
			path:    Path{"objectstep"},
			invalid: true,
		},

		// struct scalar fields

		{
			value:   ExampleStruct{},
			path:    Path{"DoesNotExist"},
			invalid: true,
		},
		{
			value:   ExampleStruct{},
			path:    Path{"privateField"},
			invalid: true,
		},
		{
			value:    ExampleStruct{StringField: ""},
			path:     Path{"StringField"},
			expected: "",
		},
		{
			value:   ExampleStruct{StringField: ""},
			path:    Path{"StringField", "dummy"},
			invalid: true,
		},
		{
			value:    ExampleStruct{StringField: "foo"},
			path:     Path{"StringField"},
			expected: "foo",
		},
		{
			value:    ExampleStruct{StringPointerField: nil},
			path:     Path{"StringPointerField"},
			expected: func() *string { return nil }(),
		},
		{
			value:    ExampleStruct{StringPointerField: stringPtr},
			path:     Path{"StringPointerField"},
			expected: stringPtr,
		},
		{
			value:    ExampleStruct{FuncField: nil},
			path:     Path{"FuncField"},
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
			path:    Path{"FuncField", "test"},
			invalid: true,
		},
		{
			value:   ExampleStruct{FuncField: dummyFieldFunc},
			path:    Path{"FuncField", 0},
			invalid: true,
		},

		// root value is pointer

		{
			value:    &ExampleStruct{StringField: "foo"},
			path:     Path{"StringField"},
			expected: "foo",
		},
		{
			value:   func() *ExampleStruct { return nil }(),
			path:    Path{"StringField"},
			invalid: true,
		},

		// struct map fields

		{
			value:    ExampleStruct{},
			path:     Path{"StringMapField"},
			expected: func() map[string]string { return nil }(),
		},
		{
			value:    ExampleStruct{StringMapField: map[string]string{}},
			path:     Path{"StringMapField"},
			expected: map[string]string{},
		},
		{
			value:   ExampleStruct{StringMapField: map[string]string{}},
			path:    Path{"StringMapField", 0},
			invalid: true,
		},
		{
			value:   ExampleStruct{StringMapField: map[string]string{}},
			path:    Path{"StringMapField", "test"},
			invalid: true,
		},
		{
			value:    ExampleStruct{StringMapField: map[string]string{"test": "value"}},
			path:     Path{"StringMapField", "test"},
			expected: "value",
		},
		{
			value:    ExampleStruct{EmptyInterfaceMapField: map[string]any{"test": "value"}},
			path:     Path{"EmptyInterfaceMapField", "test"},
			expected: "value",
		},
		{
			value:    ExampleStruct{StructMapField: map[string]OtherStruct{"foo": {}}},
			path:     Path{"StructMapField", "foo", "StringPointerField"},
			expected: func() *string { return nil }(),
		},
		{
			value:    ExampleStruct{StructMapField: map[string]OtherStruct{"foo": {StringPointerField: stringPtr}}},
			path:     Path{"StructMapField", "foo", "StringPointerField"},
			expected: stringPtr,
		},
		{
			value:    ExampleStruct{StructPointerMapField: map[string]*OtherStruct{"foo": {StringField: "bar"}}},
			path:     Path{"StructPointerMapField", "foo", "StringField"},
			expected: "bar",
		},

		// embedded structs

		{
			value:    StructWithEmbed{StringField: "foo"},
			path:     Path{"StringField"},
			expected: "foo",
		},
		{
			value:    StructWithEmbed{BaseStruct: BaseStruct{StringField: "foo"}},
			path:     Path{"StringField"},
			expected: "",
		},
		{
			value:    StructWithEmbed{BaseStruct: BaseStruct{StringField: "foo"}},
			path:     Path{"BaseStruct", "StringField"},
			expected: "foo",
		},
		{
			value:    StructWithEmbed{StringField: "foo", BaseStruct: BaseStruct{StringField: "bar"}},
			path:     Path{"StringField"},
			expected: "foo",
		},
		{
			value:    StructWithEmbed{StringField: "foo", BaseStruct: BaseStruct{StringField: "bar"}},
			path:     Path{"BaseStruct", "StringField"},
			expected: "bar",
		},
		{
			value:    StructWithEmbed{BaseStruct: BaseStruct{BaseStringField: "bar"}},
			path:     Path{"BaseStringField"},
			expected: "bar",
		},

		// interface fields

		{
			value:    ExampleStruct{CustomEmptyInterfaceField: map[string]*OtherStruct{"foo": {StringField: "bar"}}},
			path:     Path{"CustomEmptyInterfaceField", "foo", "StringField"},
			expected: "bar",
		},
		{
			value:    ExampleStruct{CustomEmptyInterfaceSlicePointerField: &customEmptyInterfacesValues},
			path:     Path{"CustomEmptyInterfaceSlicePointerField", 0, "foo", "StringField"},
			expected: "bar",
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v @ %v", tc.value, tc.path), func(t *testing.T) {
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

func TestGetFunction(t *testing.T) {
	bla := func() int {
		return dummyFieldFuncReturnValue
	}

	testcases := []struct {
		value any
		path  Path
	}{
		{
			value: ExampleStruct{FuncField: dummyFieldFunc},
			path:  Path{"FuncField"},
		},
		{
			value: ExampleStruct{FuncPointerField: &bla},
			path:  Path{"FuncPointerField"},
		},
	}

	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			result, err := Get(tc.value, tc.path)
			if err != nil {
				t.Fatalf("Failed to run: %v", err)
			}

			if result == nil {
				t.Fatal("Expected func, got nil.")
			}

			f, ok := result.(func() int)
			if !ok {
				fp, ok := result.(*func() int)
				if !ok {
					t.Fatalf("Expected dummy func, got %v (%T).", result, result)
				}

				f = *fp
			}

			if val := f(); val != dummyFieldFuncReturnValue {
				t.Fatalf("Got unexpected function, expected return value %d, got %d.", dummyFieldFuncReturnValue, val)
			}
		})
	}
}
