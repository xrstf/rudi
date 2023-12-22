// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package options

import (
	"fmt"
	"reflect"

	"github.com/spf13/pflag"
)

type enumValue interface {
	fmt.Stringer
	IsValid() bool
}

type enumFlag struct {
	target enumValue
	values []enumValue
}

func newEnumFlag[T enumValue, V enumValue](value T, possibleValues ...V) *enumFlag {
	values := make([]enumValue, len(possibleValues))
	for i, pv := range possibleValues {
		values[i] = pv
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
	newValue := f.stringToEnumValue(s)
	if !newValue.IsValid() {
		return fmt.Errorf("invalid value %q, must be one of %v", s, f.values)
	}

	// replace value in the target
	reflect.ValueOf(f.target).Elem().Set(reflect.ValueOf(newValue))

	return nil
}

func (f *enumFlag) String() string {
	return f.target.String()
}

func (*enumFlag) Type() string {
	return "string"
}

func (f *enumFlag) stringToEnumValue(s string) enumValue {
	tt := reflect.TypeOf(f.target).Elem()      // e.g. turn *Coalescer type into Coalescer
	newValue := reflect.ValueOf(s).Convert(tt) // convert string to Coalescer

	return newValue.Interface().(enumValue)
}
