//go:build integration

// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"text/template"
)

func identityTplFunc[T any](v T) T {
	fmt.Printf("v: %T = %#v\n", v, v)
	return v
}

var (
	script = strings.TrimSpace(`
{{ id (index .CustomEmptyInterfaceSlicePointerField 0) }}
`)

	funcs = template.FuncMap{
		"id": identityTplFunc[any],
	}
)

// TestGoTemplate exists to play around with how it parses field access.
func TestGoTemplate(t *testing.T) {
	tmpl, err := template.New("test").Funcs(funcs).Parse(script)
	if err != nil {
		t.Fatalf("Invalid template: %v", err)
	}

	customEmptyInterfacesValues := []CustomEmptyInterface{
		map[string]*OtherStruct{"foo": {StringField: "bar"}},
		ExampleStruct{},
	}

	data := ExampleStruct{
		CustomEmptyInterfaceField: map[string]string{
			"foo": "bar",
		},
		CustomEmptyInterfaceSlicePointerField: &customEmptyInterfacesValues,
	}

	if err = tmpl.Execute(os.Stdout, data); err != nil {
		t.Fatalf("Template failed: %v", err)
	}
}
