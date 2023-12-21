// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package types

import (
	"fmt"
	"reflect"

	"github.com/spf13/pflag"
)

type enumFlag struct {
	target fmt.Stringer
	values []string
}

func newEnumFlag(value fmt.Stringer, possibleValues ...fmt.Stringer) *enumFlag {
	values := make([]string, len(possibleValues))
	for i, v := range possibleValues {
		values[i] = v.String()
	}

	return &enumFlag{
		target: value,
		values: values,
	}
}

func (f *enumFlag) Add(fs *pflag.FlagSet, longFlag string, shortFlag string, usage string) {
	fs.VarP(f, longFlag, shortFlag, fmt.Sprintf("%s, one of %v.", usage, f.values))
}

var _ pflag.Value = &enumFlag{}

func (f *enumFlag) Set(s string) error {
	exists := false
	for _, v := range f.values {
		if v == s {
			exists = true
		}
	}

	if !exists {
		return fmt.Errorf("invalid value %q, must be one of %v", s, f.values)
	}

	tt := reflect.TypeOf(f.target).Elem()      // e.g. turn *Coalescer type into Coalescer
	newValue := reflect.ValueOf(s).Convert(tt) // convert string to Coalescer

	// replace value in the target
	reflect.ValueOf(f.target).Elem().Set(newValue)

	return nil
}

func (f *enumFlag) String() string {
	return f.target.String()
}

func (*enumFlag) Type() string {
	return "string"
}
