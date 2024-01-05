// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package types

// All of these constants must be lowercased because the validation function
// normalizes the given user input to lowercase.

type Encoding string

func (e Encoding) String() string {
	return string(e)
}

func (e Encoding) IsValid() bool {
	for _, enc := range AllEncodings {
		if enc == e {
			return true
		}
	}

	return false
}

const (
	RawEncoding           Encoding = "raw"
	JsonEncoding          Encoding = "json"
	Json5Encoding         Encoding = "json5"
	YamlEncoding          Encoding = "yaml"
	YamlDocumentsEncoding Encoding = "yamldocs"
	TomlEncoding          Encoding = "toml"
)

var (
	AllEncodings = []Encoding{
		RawEncoding,
		JsonEncoding,
		Json5Encoding,
		YamlEncoding,
		YamlDocumentsEncoding,
		TomlEncoding,
	}

	InputEncodings = []Encoding{
		RawEncoding,
		JsonEncoding,
		Json5Encoding,
		YamlEncoding,
		YamlDocumentsEncoding,
		TomlEncoding,
	}

	OutputEncodings = []Encoding{
		RawEncoding,
		JsonEncoding,
		YamlEncoding,
		YamlDocumentsEncoding,
		TomlEncoding,
	}
)

type Coalescing string

func (c Coalescing) String() string {
	return string(c)
}

func (c Coalescing) IsValid() bool {
	for _, coal := range AllCoalescings {
		if coal == c {
			return true
		}
	}

	return false
}

const (
	StrictCoalescing   Coalescing = "strict"
	PedanticCoalescing Coalescing = "pedantic"
	HumaneCoalescing   Coalescing = "humane"
)

var AllCoalescings = []Coalescing{
	StrictCoalescing,
	PedanticCoalescing,
	HumaneCoalescing,
}

type VariableSource string

func (c VariableSource) String() string {
	return string(c)
}

func (c VariableSource) IsValid() bool {
	for _, coal := range AllVariableSources {
		if coal == c {
			return true
		}
	}

	return false
}

const (
	StringVariableSource      VariableSource = "string"
	FileVariableSource        VariableSource = "file"
	EnvironmentVariableSource VariableSource = "env"
)

var AllVariableSources = []VariableSource{
	StringVariableSource,
	FileVariableSource,
	EnvironmentVariableSource,
}
