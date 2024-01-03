// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package types

import (
	"go.xrstf.de/rudi/pkg/lang/ast"
)

type Function interface {
	Evaluate(ctx Context, args []ast.Expression) (any, error)

	// Description returns a short, one-line description of the function; markdown
	// can be used to highlight other function names, like "behaves similar
	// to `foo`, but …".
	Description() string
}

type BangHandler interface {
	// All functions work fine with the default bang handler ("set!", "append!", ...), except
	// for "delete!", which requires special handling to make it work as expected. Custom bang
	// handlers are useful to introducing side effects explicitly (so it becomes very clear if a
	// function in Rudi has side effects or not).
	BangHandler(ctx Context, args []ast.Expression, value any) (any, error)
}

type TupleFunction func(ctx Context, args []ast.Expression) (any, error)

type basicFunc struct {
	f    TupleFunction
	desc string
}

// NewFunction creates the lowest of low level functions in Rudi and should rarely be used by
// integrators/library developers. Use the helpers in the root package to define functions
// using reflection and pattern matching instead.
func NewFunction(f TupleFunction, description string) Function {
	return basicFunc{
		f:    f,
		desc: description,
	}
}

var _ Function = basicFunc{}

func (f basicFunc) Evaluate(ctx Context, args []ast.Expression) (any, error) {
	return f.f(ctx, args)
}

func (f basicFunc) Description() string {
	return f.desc
}

type Functions map[string]Function

func NewFunctions() Functions {
	return Functions{}
}

func (f Functions) Get(name string) (Function, bool) {
	variable, exists := f[name]
	return variable, exists
}

// Set sets/replaces the function in the current set (in-place).
// The function returns the same Functions to allow fluent access.
func (f Functions) Set(name string, fun Function) Functions {
	f[name] = fun
	return f
}

// Set removes a function from the set.
// The function returns the same Functions to allow fluent access.
func (f Functions) Delete(name string) Functions {
	delete(f, name)
	return f
}

// Add adds all functions from other to the current set.
// The function returns the same Functions to allow fluent access.
func (f Functions) Add(other Functions) Functions {
	for name, fun := range other {
		f[name] = fun
	}
	return f
}

// Remove removes all functions from this set that are part of the other set,
// to enable constructs like AllFunctions.Remove(MathFunctions)
// The function returns the same Functions to allow fluent access.
func (f Functions) Remove(other Functions) Functions {
	for name := range other {
		f.Delete(name)
	}
	return f
}

func (f Functions) DeepCopy() Functions {
	result := NewFunctions()
	for key, val := range f {
		result[key] = val
	}
	return result
}
