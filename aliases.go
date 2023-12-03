// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package rudi

import (
	"go.xrstf.de/rudi/pkg/builtin"
	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/eval/util"
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

// NewBuiltInFunctions returns a copy of all the built-in Rudi functions.
func NewBuiltInFunctions() Functions {
	return builtin.AllFunctions.DeepCopy()
}

// NewVariables returns an empty set of runtime variables.
func NewVariables() Variables {
	return types.NewVariables()
}

// NewDocument wraps any sort of data as a Rudi document.
func NewDocument(data any) (Document, error) {
	return types.NewDocument(data)
}

// RawFunction is a function that receives its raw, unevaluated child expressions as arguments.
// This is the lowest level a function can be, allowing to selectively evaluate the arguments to
// control side effects.
type RawFunction = util.RawFunction

// NewRawFunction wraps a raw function to be used in Rudi.
func NewRawFunction(f RawFunction, description string) util.Function {
	return util.NewRawFunction(f, description)
}

// LiteralFunction is a function that receives all of its arguments already evaluated, but not yet
// coalesced into specific types.
type LiteralFunction = util.LiteralFunction

// NewLiteralFunction wraps a literal function to be used in Rudi.
func NewLiteralFunction(f LiteralFunction, description string) util.Function {
	return util.NewLiteralFunction(f, description)
}
