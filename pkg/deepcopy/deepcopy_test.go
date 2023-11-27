package deepcopy

import (
	"fmt"
	"strings"
	"testing"

	"go.xrstf.de/rudi/pkg/lang/ast"

	"github.com/google/go-cmp/cmp"
)

func ptrTo[T any](val T) *T {
	return &val
}

func TestCloneScalars(t *testing.T) {
	testcases := []struct {
		input    any
		expected any
	}{
		{
			input:    nil,
			expected: nil,
		},
		{
			input:    true,
			expected: true,
		},
		{
			input:    false,
			expected: false,
		},
		{
			input:    int(4),
			expected: int(4),
		},
		{
			input:    int32(4),
			expected: int32(4),
		},
		{
			input:    int64(0),
			expected: int64(0),
		},
		{
			input:    int64(-7),
			expected: int64(-7),
		},
		{
			input:    float32(-7.43),
			expected: float32(-7.43),
		},
		{
			input:    float64(-7.43),
			expected: float64(-7.43),
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    " foo bar ",
			expected: " foo bar ",
		},
		{
			input:    []any{1, 2, 3},
			expected: []any{1, 2, 3},
		},
		{
			input:    map[string]any{"foo": "bar", "hello": 42},
			expected: map[string]any{"foo": "bar", "hello": 42},
		},
	}

	for _, testcase := range testcases {
		t.Run(fmt.Sprintf("%T: %v", testcase.input, testcase.input), func(t *testing.T) {
			cloned, err := Clone(testcase.input)
			if err != nil {
				t.Fatalf("Failed to clone: %v", err)
			}

			if !cmp.Equal(cloned, testcase.expected) {
				t.Fatalf("Unpected result:\n\n%s\n", renderDiff(testcase.expected, cloned))
			}

			if &cloned == &testcase.input {
				t.Fatal("Both input and output data point to the same memory address, no actual cloning happened.")
			}
		})
	}
}

func TestCloneMap(t *testing.T) {
	input := map[string]any{"foo": "bar", "hello": 1}

	cloned, err := Clone(input)
	if err != nil {
		t.Fatalf("Failed to clone: %v", err)
	}

	input["new"] = "new-value"
	if _, ok := cloned["new"]; ok {
		t.Fatal("Both input and output data point to the same memory address, no actual cloning happened.")
	}
}

func TestCloneMapDeep(t *testing.T) {
	input := map[string]any{
		"foo": "bar",
		"hello": []any{
			1,
			2,
			"foo",
			map[string]any{
				"deep": "value",
				"keep": nil,
			},
			[]any{"sub", "list"},
		},
	}

	cloned, err := Clone(input)
	if err != nil {
		t.Fatalf("Failed to clone: %v", err)
	}

	if !cmp.Equal(cloned, input) {
		t.Fatalf("Unpected result:\n\n%s\n", renderDiff(input, cloned))
	}

	helloList := input["hello"].([]any)
	helloObj := helloList[3].(map[string]any)
	helloObj["deep"] = "new-value"

	if cmp.Equal(cloned, input) {
		t.Fatal("Changing the input changed the output, no actual deep cloning happened.")
	}
}

func TestCloneSlice(t *testing.T) {
	input := []any{1, 2, "foo"}

	cloned, err := Clone(input)
	if err != nil {
		t.Fatalf("Failed to clone: %v", err)
	}

	input[1] = "new"
	if cloned[1] == "new" {
		t.Fatal("Both input and output data point to the same memory address, no actual cloning happened.")
	}
}

func TestCloneSliceDeep(t *testing.T) {
	input := []any{
		1,
		2,
		"foo",
		map[string]any{
			"deep": "value",
			"keep": nil,
		},
		[]any{"sub", "list"},
	}

	cloned, err := Clone(input)
	if err != nil {
		t.Fatalf("Failed to clone: %v", err)
	}

	if !cmp.Equal(cloned, input) {
		t.Fatalf("Unpected result:\n\n%s\n", renderDiff(input, cloned))
	}

	helloObj := input[3].(map[string]any)
	helloObj["deep"] = "new-value"

	if cmp.Equal(cloned, input) {
		t.Fatal("Changing the input changed the output, no actual deep cloning happened.")
	}
}

func TestCloneScalarPointers(t *testing.T) {
	testcases := []struct {
		input    any
		expected any
	}{
		{
			input:    ptrTo(true),
			expected: ptrTo(true),
		},
		{
			input:    ptrTo(false),
			expected: ptrTo(false),
		},
		{
			input:    ptrTo(int(4)),
			expected: ptrTo(int(4)),
		},
		{
			input:    ptrTo(int32(4)),
			expected: ptrTo(int32(4)),
		},
		{
			input:    ptrTo(int64(0)),
			expected: ptrTo(int64(0)),
		},
		{
			input:    ptrTo(float32(-7.43)),
			expected: ptrTo(float32(-7.43)),
		},
		{
			input:    ptrTo(float64(-7.43)),
			expected: ptrTo(float64(-7.43)),
		},
		{
			input:    ptrTo("foo bar"),
			expected: ptrTo("foo bar"),
		},
		{
			input:    ptrTo([]any{1, 2, 3}),
			expected: ptrTo([]any{1, 2, 3}),
		},
		{
			input:    ptrTo(map[string]any{"foo": "bar", "hello": 42}),
			expected: ptrTo(map[string]any{"foo": "bar", "hello": 42}),
		},
	}

	for _, testcase := range testcases {
		t.Run(fmt.Sprintf("%T: %v", testcase.input, testcase.input), func(t *testing.T) {
			cloned, err := Clone(testcase.input)
			if err != nil {
				t.Fatalf("Failed to clone: %v", err)
			}

			if !cmp.Equal(cloned, testcase.expected) {
				t.Fatalf("Unpected result:\n\n%s\n", renderDiff(testcase.expected, cloned))
			}

			if cloned == testcase.input {
				t.Fatal("Both input and output data point to the same memory address, no actual cloning happened.")
			}
		})
	}
}

func TestCloneMapPointer(t *testing.T) {
	input := map[string]any{"foo": "bar", "hello": 1}

	cloned, err := Clone(&input)
	if err != nil {
		t.Fatalf("Failed to clone: %v", err)
	}

	input["new"] = "new-value"
	if _, ok := (*cloned)["new"]; ok {
		t.Fatal("Both input and output data point to the same memory address, no actual cloning happened.")
	}
}

func TestCloneSlicePointer(t *testing.T) {
	input := []any{1, 2, "foo"}

	cloned, err := Clone(&input)
	if err != nil {
		t.Fatalf("Failed to clone: %v", err)
	}

	input[1] = "new"
	if (*cloned)[1] == "new" {
		t.Fatal("Both input and output data point to the same memory address, no actual cloning happened.")
	}
}

func TestCloneNodes(t *testing.T) {
	testcases := []struct {
		input    any
		expected any
	}{
		{
			input:    ast.Null{},
			expected: ast.Null{},
		},
		{
			input:    ast.Bool(true),
			expected: ast.Bool(true),
		},
		{
			input:    ast.Bool(false),
			expected: ast.Bool(false),
		},
		{
			input:    ast.Number{Value: 1},
			expected: ast.Number{Value: 1},
		},
		{
			input:    ast.Number{Value: -3.14},
			expected: ast.Number{Value: -3.14},
		},
		{
			input:    ast.String(""),
			expected: ast.String(""),
		},
		{
			input:    ast.String(" test "),
			expected: ast.String(" test "),
		},
		{
			input:    ast.Vector{Data: []any{1, "foo", ast.Bool(true), ast.Number{Value: 2}}},
			expected: ast.Vector{Data: []any{1, "foo", ast.Bool(true), ast.Number{Value: 2}}},
		},
		{
			input:    ast.Object{Data: map[string]any{"foo": "bar", "hello": ast.Bool(true)}},
			expected: ast.Object{Data: map[string]any{"foo": "bar", "hello": ast.Bool(true)}},
		},
	}

	for _, testcase := range testcases {
		t.Run(fmt.Sprintf("%T: %v", testcase.input, testcase.input), func(t *testing.T) {
			cloned, err := Clone(testcase.input)
			if err != nil {
				t.Fatalf("Failed to clone: %v", err)
			}

			if !cmp.Equal(cloned, testcase.expected) {
				t.Fatalf("Unpected result:\n\n%s\n", renderDiff(testcase.expected, cloned))
			}

			if &cloned == &testcase.input {
				t.Fatal("Both input and output data point to the same memory address, no actual cloning happened.")
			}
		})
	}
}

func TestCloneObject(t *testing.T) {
	input := ast.Object{Data: map[string]any{"foo": "bar", "hello": 1}}

	cloned, err := Clone(input)
	if err != nil {
		t.Fatalf("Failed to clone: %v", err)
	}

	input.Data["new"] = "new-value"
	if _, ok := cloned.Data["new"]; ok {
		t.Fatal("Both input and output data point to the same memory address, no actual cloning happened.")
	}
}

func TestCloneVector(t *testing.T) {
	input := ast.Vector{Data: []any{1, 2, "foo"}}

	cloned, err := Clone(input)
	if err != nil {
		t.Fatalf("Failed to clone: %v", err)
	}

	input.Data[1] = "new"
	if cloned.Data[1] == "new" {
		t.Fatal("Both input and output data point to the same memory address, no actual cloning happened.")
	}
}

func TestCloneNodePointers(t *testing.T) {
	testcases := []struct {
		input    any
		expected any
	}{
		// TODO: these are hard to compare for Go as long as ast.Null is an empty struct.
		// {
		// 	input:    &ast.Null{},
		// 	expected: &ast.Null{},
		// },
		// {
		// 	input:    new(ast.Null),
		// 	expected: &ast.Null{},
		// },
		{
			input:    new(ast.Bool),
			expected: ptrTo(ast.Bool(false)),
		},
		{
			input:    ptrTo(ast.Bool(true)),
			expected: ptrTo(ast.Bool(true)),
		},
		{
			input:    ptrTo(ast.Bool(false)),
			expected: ptrTo(ast.Bool(false)),
		},
		{
			input:    new(ast.Number), // invalid Number
			expected: new(ast.Number),
		},
		{
			input:    &ast.Number{Value: 1},
			expected: &ast.Number{Value: 1},
		},
		{
			input:    &ast.Number{Value: -3.14},
			expected: &ast.Number{Value: -3.14},
		},
		{
			input:    new(ast.String),
			expected: ptrTo(ast.String("")),
		},
		{
			input:    ptrTo(ast.String("")),
			expected: ptrTo(ast.String("")),
		},
		{
			input:    ptrTo(ast.String(" test ")),
			expected: ptrTo(ast.String(" test ")),
		},
		{
			input:    ptrTo(ast.Vector{Data: []any{1, "foo", ast.Bool(true), ast.Number{Value: 2}}}),
			expected: ptrTo(ast.Vector{Data: []any{1, "foo", ast.Bool(true), ast.Number{Value: 2}}}),
		},
		{
			input:    ptrTo(ast.Object{Data: map[string]any{"foo": "bar", "hello": ast.Bool(true)}}),
			expected: ptrTo(ast.Object{Data: map[string]any{"foo": "bar", "hello": ast.Bool(true)}}),
		},
	}

	for _, testcase := range testcases {
		t.Run(fmt.Sprintf("%T: %v", testcase.input, testcase.input), func(t *testing.T) {
			cloned, err := Clone(testcase.input)
			if err != nil {
				t.Fatalf("Failed to clone: %v", err)
			}

			if !cmp.Equal(cloned, testcase.expected) {
				t.Fatalf("Unpected result:\n\n%s\n", renderDiff(testcase.expected, cloned))
			}

			if cloned == testcase.input {
				t.Fatal("Both input and output data point to the same memory address, no actual cloning happened.")
			}
		})
	}
}

func TestCloneObjectPointer(t *testing.T) {
	input := ast.Object{Data: map[string]any{"foo": "bar", "hello": 1}}

	cloned, err := Clone(&input)
	if err != nil {
		t.Fatalf("Failed to clone: %v", err)
	}

	input.Data["new"] = "new-value"
	if _, ok := (*cloned).Data["new"]; ok {
		t.Fatal("Both input and output data point to the same memory address, no actual cloning happened.")
	}
}

func TestCloneVectorPointer(t *testing.T) {
	input := ast.Vector{Data: []any{1, 2, "foo"}}

	cloned, err := Clone(&input)
	if err != nil {
		t.Fatalf("Failed to clone: %v", err)
	}

	input.Data[1] = "new"
	if (*cloned).Data[1] == "new" {
		t.Fatal("Both input and output data point to the same memory address, no actual cloning happened.")
	}
}

func renderDiff(expected any, actual any) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Expected type...: %T\n", expected))
	builder.WriteString(fmt.Sprintf("Expected value..: %#v\n", expected))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("Actual type.....: %T\n", actual))
	builder.WriteString(fmt.Sprintf("Actual value....: %#v\n", actual))

	return builder.String()
}
