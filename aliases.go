// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package rudi

import (
	"go.xrstf.de/rudi/pkg/builtin"
	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/eval/functions"
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

// NewSafeBuiltInFunctions returns a copy of all the safe built-in Rudi functions. These are all the
// functions that do not break runtime guarantees like programs always terminating in a reasonable
// time. See also NewUnsafeBuiltInFunctions, which contains functions like func! that allow to define
// new functions within Rudi code, but could lead to infinite loops or resource exhaustion.
func NewSafeBuiltInFunctions() Functions {
	return builtin.SafeFunctions.DeepCopy()
}

// NewUnsafeBuiltInFunctions returns a copy of all the unsafe built-in Rudi functions. These are
// functions with extended side effects, please refer to the documentation or code for which
// functions exactly are considered "unsafe" in Rudi.
func NewUnsafeBuiltInFunctions() Functions {
	return builtin.UnsafeFunctions.DeepCopy()
}

// NewVariables returns an empty set of runtime variables.
func NewVariables() Variables {
	return types.NewVariables()
}

// NewDocument wraps any sort of data as a Rudi document.
func NewDocument(data any) (Document, error) {
	return types.NewDocument(data)
}

// NewFunctionBuilder is the recommended way to define new Rudi functions. The function builder can
// take multiple forms (e.g. if you have (foo INT) and (foo STRING)) and will create a function that
// automatically evaluates and coalesces Rudi expressions and matches them to the given forms. The
// first matching form is then evaluated.
func NewFunctionBuilder(forms ...any) *functions.Builder {
	return functions.NewBuilder(forms...)
}

// NewLowLevelFunction wraps a raw tuple function to be used in Rudi. This is mostly useful for
// defining really low-level functions and functions with special side effects. Most of the time,
// you'd want to use NewFunctionBuilder(), which will use reflection to make it much more straight
// forward to make a Go function available in Rudi.
func NewLowLevelFunction(f types.TupleFunction, description string) Function {
	return types.NewFunction(f, description)
}
