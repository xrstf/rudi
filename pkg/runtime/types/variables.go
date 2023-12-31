// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package types

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

// WithMany is like With(), but for adding multiple new variables at once. This
// should be preferred to With() to prevent unnecessary DeepCopies.
func (v Variables) WithMany(vars map[string]any) Variables {
	if len(vars) == 0 {
		return v
	}

	out := v.DeepCopy()
	for k, v := range vars {
		out.Set(k, v)
	}
	return out
}

func (v Variables) DeepCopy() Variables {
	result := NewVariables()
	for key, val := range v {
		result[key] = val
	}
	return result
}
