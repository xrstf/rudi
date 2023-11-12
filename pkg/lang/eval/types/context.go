// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package types

type Document struct {
	Data any
}

func NewDocument(data any) Document {
	return Document{
		Data: data,
	}
}

type Context struct {
	document  Document
	variables Variables
}

func NewContext(doc Document, variables Variables) Context {
	return Context{
		document:  doc,
		variables: variables,
	}
}

func (c Context) GetDocument() Document {
	return c.document
}

func (c Context) GetVariable(name string) (any, bool) {
	return c.variables.Get(name)
}

func (c Context) WithVariable(name string, val any) Context {
	return Context{
		document:  c.document,
		variables: c.variables.With(name, val),
	}
}

type Variables map[string]any

func NewVariables() Variables {
	return Variables{}
}

func (v Variables) Get(name string) (any, bool) {
	variable, exists := v[name]
	return variable, exists
}

// Set sets/replaces the variable value in the current set (in-place).
// The function returns the same variables to allow fluent access.
func (v Variables) Set(name string, val any) Variables {
	v[name] = val
	return v
}

// With returns a copy of the variables, with the new variable being added to it.
func (v Variables) With(name string, val any) Variables {
	return v.DeepCopy().Set(name, val)
}

func (v Variables) DeepCopy() Variables {
	result := NewVariables()
	for key, val := range v {
		result[key] = val
	}
	return result
}
