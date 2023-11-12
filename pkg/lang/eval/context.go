// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import "go.xrstf.de/corel/pkg/lang/eval/types"

func NewContext(doc types.Document, variables types.Variables) types.Context {
	return types.NewContext(doc, variables)
}

func NewVariables() types.Variables {
	return types.NewVariables()
}

func NewDocument(data any) types.Document {
	return types.NewDocument(data)
}
