// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package rudi

import (
	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/eval/builtin"
	"go.xrstf.de/rudi/pkg/eval/types"
)

// Context is the evaluation context for a Rudi program, consisting of
// the global document, variables and functions.
type Context = types.Context

// Variables is a map of Rudi variables.
type Variables = types.Variables

// Functions is a map of Rudi functions.
type Functions = types.Functions

// Function is a single Rudi function, available to be used inside a Rudi script.
type Function = types.Function

// Document is the global document that is being processed by a Rudi script.
type Document = types.Document

// Coalescer is responsible for type handling and equality rules. Build your own
// or use any of the predefined versions:
//
//   - coalescing.NewStrict() – mostly strict, but allows nulls to be converted
//     and allows ints to become floats
//   - coalescing.NewPedantic() – even more strict, allows absolutely no conversions
//   - coalescing.NewHumane() – gentle type handling that allows lossless
//     conversions like 1 => "1" or allowing (false == nil).
type Coalescer = coalescing.Coalescer

// NewContext wraps the document, variables and functions into a Context.
func NewContext(doc Document, variables Variables, funcs Functions, coalescer Coalescer) Context {
	return types.NewContext(doc, variables, funcs, coalescer)
}

// NewFunctions returns an empty set of runtime functions.
func NewFunctions() Functions {
	return types.NewFunctions()
}

// NewBuiltInFunctions returns a copy of the built-in Rudi functions.
func NewBuiltInFunctions() Functions {
	return builtin.Functions.DeepCopy()
}

// NewVariables returns an empty set of runtime variables.
func NewVariables() Variables {
	return types.NewVariables()
}

// NewDocument wraps any sort of data as a Rudi document.
func NewDocument(data any) (Document, error) {
	return types.NewDocument(data)
}

// Unwrap returns the native Go value for either native Go values or an
// Rudi AST node (like turning an ast.Number into an int64).
func Unwrap(val any) (any, error) {
	return types.UnwrapType(val)
}

// WrapNative returns the Rudi node equivalent of a native Go value, like turning
// a string into ast.String.
func WrapNative(val any) (any, error) {
	return types.WrapNative(val)
}
