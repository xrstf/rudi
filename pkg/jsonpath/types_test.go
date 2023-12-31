// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

func ptrTo[T any](v T) *T {
	return &v
}

type emptyStruct struct{}

type CustomEmptyInterface interface{}

type NonEmptyInterface interface {
	Foo()
}

type nonEmptyImplementation struct{}

var _ NonEmptyInterface = nonEmptyImplementation{}

func (nonEmptyImplementation) Foo() {}

const dummyFieldFuncReturnValue = 42

func dummyFieldFunc() int {
	return dummyFieldFuncReturnValue
}

type BaseStruct struct {
	privateBaseField int

	StringField     string
	BaseStringField string
}

type StructWithEmbed struct {
	BaseStruct

	StringField string
}

type ExampleStruct struct {
	privateField int

	StringField                           string
	StringPointerField                    *string
	BoolField                             bool
	BoolPointerField                      *bool
	FuncField                             func() int
	FuncPointerField                      *func() int
	EmptyInterfaceField                   any
	EmptyInterfaceSliceField              []any
	EmptyInterfaceSlicePointerField       *[]any
	CustomEmptyInterfaceField             CustomEmptyInterface
	CustomEmptyInterfaceSliceField        []CustomEmptyInterface
	CustomEmptyInterfaceSlicePointerField *[]CustomEmptyInterface
	NonEmptyInterfaceField                NonEmptyInterface
	NonEmptyInterfaceSliceField           []NonEmptyInterface
	NonEmptyInterfaceSlicePointerField    *[]NonEmptyInterface
	StructField                           OtherStruct
	StructPointerField                    *OtherStruct
	StructSliceField                      []OtherStruct
	StructSlicePointerField               *[]OtherStruct
	StructPointerSliceField               []*OtherStruct
	StructPointerSlicePointerField        *[]*OtherStruct
	StringMapField                        map[string]string
	EmptyInterfaceMapField                map[string]any
	CustomEmptyInterfaceMapField          map[string]CustomEmptyInterface
	NonEmptyInterfaceMapField             map[string]NonEmptyInterface
	StructMapField                        map[string]OtherStruct
	StructPointerMapField                 map[string]*OtherStruct
}

type OtherStruct struct {
	privateField int

	StringField        string
	StringPointerField *string
}
