package deepcopy

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type customCopier struct {
	Value any
}

var _ Copier = customCopier{}

func (c customCopier) DeepCopy() (any, error) {
	copied, err := Clone(c.Value)
	if err != nil {
		return nil, err
	}

	return customCopier{
		Value: copied,
	}, nil
}

type customPtrCopier struct {
	Value any
}

var _ Copier = &customPtrCopier{}

func (c *customPtrCopier) DeepCopy() (any, error) {
	copied, err := Clone(c.Value)
	if err != nil {
		return nil, err
	}

	return &customPtrCopier{
		Value: copied,
	}, nil
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
				t.Fatalf("Unexpected result:\n\n%s\n", renderDiff(testcase.expected, cloned))
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
		t.Fatalf("Unexpected result:\n\n%s\n", renderDiff(input, cloned))
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
		t.Fatalf("Unexpected result:\n\n%s\n", renderDiff(input, cloned))
	}

	helloObj := input[3].(map[string]any)
	helloObj["deep"] = "new-value"

	if cmp.Equal(cloned, input) {
		t.Fatal("Changing the input changed the output, no actual deep cloning happened.")
	}
}

func TestCustomCopier(t *testing.T) {
	data := map[string]any{"foo": "bar"}
	cc := customCopier{Value: data}

	cloned, err := Clone(cc)
	if err != nil {
		t.Fatalf("Failed to clone: %v", err)
	}

	if !cmp.Equal(cc, cloned) {
		t.Fatalf("Unexpected result:\n\n%s\n", renderDiff(cc, cloned))
	}
}

func TestCustomPtrCopier(t *testing.T) {
	data := map[string]any{"foo": "bar"}
	cc := &customPtrCopier{Value: data}

	cloned, err := Clone(cc)
	if err != nil {
		t.Fatalf("Failed to clone: %v", err)
	}

	if !cmp.Equal(cc, cloned) {
		t.Fatalf("Unexpected result:\n\n%s\n", renderDiff(cc, cloned))
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
