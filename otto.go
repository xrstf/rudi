// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package otto

import (
	"fmt"

	"go.xrstf.de/otto/pkg/eval"
	"go.xrstf.de/otto/pkg/eval/builtin"
	"go.xrstf.de/otto/pkg/eval/types"
	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/parser"
)

// alias types
type Context = types.Context
type Variables = types.Variables
type Functions = types.Functions
type Function = types.Function
type Document = types.Document
type Program = ast.Program

func NewContext(doc Document, funcs Functions, variables Variables) Context {
	return types.NewContext(doc, funcs, variables)
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

func ParseScript(name string, script string) (*Program, error) {
	got, err := parser.Parse(name, []byte(script))
	if err != nil {
		return nil, ParseError{script: script, err: err}
	}

	program, ok := got.(Program)
	if !ok {
		// this should never happen
		return nil, fmt.Errorf("parsed input is not an ast.Program, but %T", got)
	}

	return &program, nil
}

func RunProgram(ctx Context, program *Program) (Context, any, error) {
	return eval.Run(ctx, program)
}
