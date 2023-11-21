// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import "go.xrstf.de/rudi/pkg/eval/types"

func NewContext(doc types.Document, funcs types.Functions, variables types.Variables) types.Context {
	return types.NewContext(doc, funcs, variables)
}

func NewFunctions() types.Functions {
	return types.NewFunctions()
}

func NewVariables() types.Variables {
	return types.NewVariables()
}

func NewDocument(data any) (types.Document, error) {
	return types.NewDocument(data)
}
