// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package rudi

import (
	"go.xrstf.de/rudi/pkg/eval/builtin"
	"go.xrstf.de/rudi/pkg/eval/types"
)

// alias types

type Context = types.Context
type Variables = types.Variables
type Functions = types.Functions
type Function = types.Function
type Document = types.Document

func NewContext(doc Document, variables Variables, funcs Functions) Context {
	return types.NewContext(doc, variables, funcs)
}

func NewFunctions() Functions {
	return types.NewFunctions()
}

func NewBuiltInFunctions() Functions {
	return builtin.Functions.DeepCopy()
}

func NewVariables() Variables {
	return types.NewVariables()
}

func NewDocument(data any) (Document, error) {
	return types.NewDocument(data)
}

func UnwrapType(val any) (any, error) {
	return types.UnwrapType(val)
}

func WrapNative(val any) (any, error) {
	return types.WrapNative(val)
}
