// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInvalidSets(t *testing.T) {
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
			path:    Path{KeyStep("foo")},
			invalid: true,
		},
		{
			name:    "invalid step",
			dest:    "value",
			path:    Path{true},
			invalid: true,
		},
		{
			name:    "cannot set anything in types that do not implement the Writer interfaces",
			dest:    unknownType{},
			path:    Path{KeyStep("foo")},
			invalid: true,
		},
		{
			name:    "cannot set anything in types that do not implement the Writer interfaces",
			dest:    unknownType{},
			path:    Path{IndexStep(0)},
			invalid: true,
		},

		// custom types

		// {
		// 	name: "can set in custom object writer",
		// 	dest: &customObjWriter{
		// 		Value: "old",
		// 	},
		// 	path:     Path{KeyStep("value")},
		// 	newValue: "new-value",
		// 	expected: &customObjWriter{
		// 		Value: "new-value",
		// 	},
		// },
		// {
		// 	name: "can set deeper in custom object writer",
		// 	dest: &customObjWriter{
		// 		Value: map[string]any{
		// 			"foo": "old",
		// 		},
		// 	},
		// 	path:     Path{KeyStep("value"), KeyStep("foo")},
		// 	newValue: "new-value",
		// 	expected: &customObjWriter{
		// 		Value: map[string]any{
		// 			"foo": "new-value",
		// 		},
		// 	},
		// },
		// {
		// 	name: "vector steps on custom object getter should fail",
		// 	dest: &customObjWriter{
		// 		Value: "old",
		// 	},
		// 	path:     Path{0},
		// 	newValue: "new-value",
		// 	invalid:  true,
		// },

		// {
		// 	name: "can set in custom vector writer",
		// 	dest: &customVecWriter{
		// 		Magic: 7,
		// 		Value: "old",
		// 	},
		// 	path:     Path{IndexStep(7)},
		// 	newValue: "new-value",
		// 	expected: &customVecWriter{
		// 		Magic: 7,
		// 		Value: "new-value",
		// 	},
		// },
		// {
		// 	name: "can set deeper in custom vector writer",
		// 	dest: &customVecWriter{
		// 		Magic: 7,
		// 		Value: map[string]any{
		// 			"foo": "old",
		// 		},
		// 	},
		// 	path:     Path{IndexStep(7), KeyStep("foo")},
		// 	newValue: "new-value",
		// 	expected: &customVecWriter{
		// 		Magic: 7,
		// 		Value: map[string]any{
		// 			"foo": "new-value",
		// 		},
		// 	},
		// },
		// {
		// 	name: "object steps on custom vector getter should fail",
		// 	dest: &customVecWriter{
		// 		Magic: 7,
		// 		Value: "old",
		// 	},
		// 	path:     Path{KeyStep("foo")},
		// 	newValue: "new-value",
		// 	invalid:  true,
		// },
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

type setTestcase struct {
	name         string
	objJSON      string
	path         Path
	newValue     any
	expectedJSON string
	invalid      bool
}

func (tc *setTestcase) Run(t *testing.T) {
	(&patchTestcase{
		name:    tc.name,
		objJSON: tc.objJSON,
		path:    tc.path,
		patch: func(_ *testing.T, _ bool, _ any, _ any) (any, error) {
			return tc.newValue, nil
		},
		expectedJSON: tc.expectedJSON,
		invalid:      tc.invalid,
	}).Run(t)
}

func TestSet(t *testing.T) {
	testcases := []setTestcase{
		{
			name:         "scalar root value can simply be changed",
			objJSON:      `null`,
			path:         Path{},
			newValue:     "foo",
			expectedJSON: `"foo"`,
		},
		{
			name:         "scalar root value can simply be changed",
			objJSON:      `"hello world"`,
			path:         Path{},
			newValue:     "new value",
			expectedJSON: `"new value"`,
		},
		{
			name:         "nils can be turned into objects",
			objJSON:      `null`,
			path:         Path{KeyStep("foo")},
			newValue:     "bar",
			expectedJSON: `{"foo": "bar"}`,
		},
		{
			name:         "nils can turn into slices",
			objJSON:      `null`,
			path:         Path{IndexStep(0)},
			newValue:     "bar",
			expectedJSON: `["bar"]`,
		},
		{
			name:     "only nils can type shift",
			objJSON:  `"a string"`,
			path:     Path{KeyStep("foo")},
			newValue: "bar",
			invalid:  true,
		},
		{
			name:     "only nils can type shift",
			objJSON:  `42`,
			path:     Path{KeyStep("foo")},
			newValue: "bar",
			invalid:  true,
		},
		{
			name:         "root object key can be updated",
			objJSON:      `{"foo": "bar"}`,
			path:         Path{KeyStep("foo")},
			newValue:     "new-value",
			expectedJSON: `{"foo": "new-value"}`,
		},
		{
			name:         "root object key can be added",
			objJSON:      `{"foo": "bar"}`,
			path:         Path{KeyStep("test")},
			newValue:     "new-value",
			expectedJSON: `{"foo": "bar", "test": "new-value"}`,
		},
		{
			name:         "root slice can be updated",
			objJSON:      `[1, 2, 3]`,
			path:         Path{IndexStep(1)},
			newValue:     "new-value",
			expectedJSON: `[1, "new-value", 3]`,
		},
		{
			name:     "handle out of bounds",
			objJSON:  `[1, 2, 3]`,
			path:     Path{IndexStep(-1)},
			newValue: "new-value",
			invalid:  true,
		},
		{
			name:         "can extend vectors",
			objJSON:      `[1, 2, 3]`,
			path:         Path{IndexStep(3)},
			newValue:     "new-value",
			expectedJSON: `[1, 2, 3, "new-value"]`,
		},
		{
			name:         "sub object key can be updated",
			objJSON:      `{"foo": "bar", "deeper": {"deep": "value", "other": "value"}}`,
			path:         Path{KeyStep("deeper"), KeyStep("deep")},
			newValue:     "new-value",
			expectedJSON: `{"foo": "bar", "deeper": {"deep": "new-value", "other": "value"}}`,
		},
		{
			name:         "sub slice key can be updated",
			objJSON:      `{"foo": "bar", "deeper": [1, 2, {"deep": "value"}]}`,
			path:         Path{KeyStep("deeper"), IndexStep(2), KeyStep("deep")},
			newValue:     "new-value",
			expectedJSON: `{"foo": "bar", "deeper": [1, 2, {"deep": "new-value"}]}`,
		},
		{
			name:     "cannot turn slice into object by accident",
			objJSON:  `{"foo": "bar", "deeper": [1, 2, {"deep": "value"}]}`,
			path:     Path{KeyStep("deeper"), KeyStep("whoops")},
			newValue: "new-value",
			invalid:  true,
		},
		{
			name:         "can change value types",
			objJSON:      `{"foo": "bar", "deeper": [1, 2, {"deep": "value"}]}`,
			path:         Path{KeyStep("deeper"), IndexStep(2)},
			newValue:     "new-value",
			expectedJSON: `{"foo": "bar", "deeper": [1, 2, "new-value"]}`,
		},
		{
			name:         "can extend as needed (simple)",
			objJSON:      `{"foo": "bar"}`,
			path:         Path{KeyStep("deep"), KeyStep("deeper")},
			newValue:     "new-value",
			expectedJSON: `{"foo": "bar", "deep": {"deeper": "new-value"}}`,
		},
		{
			name:         "can extend as needed (nulls)",
			objJSON:      `{"foo": null}`,
			path:         Path{KeyStep("foo"), KeyStep("deeper")},
			newValue:     "new-value",
			expectedJSON: `{"foo": {"deeper": "new-value"}}`,
		},
		{
			name:         "can extend as needed (root nulls)",
			objJSON:      `null`,
			path:         Path{KeyStep("foo"), KeyStep("deeper")},
			newValue:     "new-value",
			expectedJSON: `{"foo": {"deeper": "new-value"}}`,
		},
		{
			name:         "can extend as needed (deep)",
			objJSON:      `{"foo": "bar"}`,
			path:         Path{KeyStep("deep"), IndexStep(2), IndexStep(1), KeyStep("bar")},
			newValue:     "new-value",
			expectedJSON: `{"foo": "bar", "deep": [null, null, [null, {"bar": "new-value"}]]}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, tc.Run)
	}
}

func TestMultiSet(t *testing.T) {
	testcases := []setTestcase{
		{
			name:         "non-existing keys are not added if a multi step is used",
			objJSON:      `null`,
			path:         Path{keySelector{"foo"}},
			newValue:     "foo",
			expectedJSON: `null`,
		},
		{
			name:         "non-existing indexes are not added if a multi step is used",
			objJSON:      `null`,
			path:         Path{indexSelector{0}},
			newValue:     "foo",
			expectedJSON: `null`,
		},
		{
			name:    "multi steps yield an error if the type doesn't match",
			objJSON: `"a string"`,
			path:    Path{indexSelector{0}},
			invalid: true,
		},
		{
			name:         "can set multiple keys in a vector at once",
			objJSON:      `["a", "b", "c"]`,
			path:         Path{indexSelector{0, 2}},
			newValue:     "foo",
			expectedJSON: `["foo", "b", "foo"]`,
		},
		{
			name:         "will not add new indexes to a vector",
			objJSON:      `["a", "b", "c"]`,
			path:         Path{indexSelector{0, 2, 4}},
			newValue:     "foo",
			expectedJSON: `["foo", "b", "foo"]`,
		},
		{
			name:         "will not set new keys in object",
			objJSON:      `{}`,
			path:         Path{keySelector{"foo"}},
			newValue:     "foo",
			expectedJSON: `{}`,
		},
		{
			name:         "will not set new keys in object",
			objJSON:      `{"foo": "bar"}`,
			path:         Path{keySelector{"foo", "baz"}},
			newValue:     "new",
			expectedJSON: `{"foo": "new"}`,
		},
		{
			name:         "can overwrite simple object values",
			objJSON:      `{"foo": "bar"}`,
			path:         Path{keySelector{"foo"}},
			newValue:     "baz",
			expectedJSON: `{"foo": "baz"}`,
		},
		{
			name:         "non-multi steps will descend and create deeper values",
			objJSON:      `{"foo": null}`,
			path:         Path{keySelector{"foo"}, KeyStep("bar")},
			newValue:     "baz",
			expectedJSON: `{"foo": {"bar": "baz"}}`,
		},
		{
			name:         "non-multi steps will descend and create deeper values",
			objJSON:      `{"foo": null}`,
			path:         Path{keySelector{"foo"}, IndexStep(2), KeyStep("deep")},
			newValue:     "bar",
			expectedJSON: `{"foo": [null, null, {"deep": "bar"}]}`,
		},
		{
			name:         "can shrink result set with additional selectors",
			objJSON:      `{"foo": {"hello": "world"}}`,
			path:         Path{keySelector{"foo"}, keySelector{"bla"}},
			newValue:     "bar",
			expectedJSON: `{"foo": {"hello": "world"}}`,
		},
		{
			name:         "can shrink result set with additional selectors",
			objJSON:      `{"foo": {"hello": "world"}}`,
			path:         Path{keySelector{"foo"}, keySelector{"hello"}},
			newValue:     "bar",
			expectedJSON: `{"foo": {"hello": "bar"}}`,
		},
		{
			name:         "can overwrite deeper values found by multi steps",
			objJSON:      `{"foo": {"hello": "world"}}`,
			path:         Path{keySelector{"foo"}},
			newValue:     "bar",
			expectedJSON: `{"foo": "bar"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, tc.Run)
	}
}

type patchTestcase struct {
	name         string
	objJSON      string
	path         Path
	patch        func(t *testing.T, exists bool, key any, value any) (any, error)
	expectedJSON string
	invalid      bool
}

func (tc *patchTestcase) Run(t *testing.T) {
	var data interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(tc.objJSON)), &data); err != nil {
		t.Fatalf("Invalid obj in testcase: %v", err)
	}

	patcher := func(exists bool, key any, value any) (any, error) {
		return tc.patch(t, exists, key, value)
	}

	result, err := Patch(data, tc.path, patcher)
	if err != nil {
		if !tc.invalid {
			t.Fatalf("Failed to run: %v", err)
		}

		return
	}

	if tc.invalid {
		t.Fatalf("Should not have been able to patch value, but got: %v (%T)", result, result)
	}

	var expected interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(tc.expectedJSON)), &expected); err != nil {
		t.Fatalf("Invalid expected in testcase: %v", err)
	}

	if !cmp.Equal(expected, result) {
		t.Fatalf("Expected %v (%T), but got %v (%T)", expected, expected, result, result)
	}
}

func TestPatch(t *testing.T) {
	testcases := []patchTestcase{
		{
			name:    "scalar root value can simply be changed",
			objJSON: `null`,
			path:    Path{},
			patch: func(t *testing.T, exists bool, key any, val any) (any, error) {
				if !exists {
					t.Fatal("exists should have been true")
				}

				if val != nil {
					t.Fatalf("val should have been nil, but is %v (%T)", val, val)
				}

				return "foo", nil
			},
			expectedJSON: `"foo"`,
		},

		{
			name:    "can change an object's value",
			objJSON: `{"foo": "bar"}`,
			path:    Path{KeyStep("foo")},
			patch: func(t *testing.T, exists bool, key any, val any) (any, error) {
				if !exists {
					t.Fatal("exists should have been true")
				}

				s, ok := val.(string)
				if !ok {
					t.Fatalf("val should have been string, but is %v (%T)", val, val)
				}

				if s != "bar" {
					t.Fatalf("val should have been bar, but is %q", s)
				}

				return "foo", nil
			},
			expectedJSON: `{"foo": "foo"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, tc.Run)
	}
}
