// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package encoding

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"go.xrstf.de/rudi/cmd/rudi/types"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

func newYamlEncoder(out io.Writer) *yaml.Encoder {
	encoder := yaml.NewEncoder(out)
	encoder.SetIndent(2)

	return encoder
}

func Encode(data any, enc types.Encoding, out io.Writer) error {
	var encoder interface {
		Encode(v any) error
	}

	switch enc {
	case types.JsonEncoding:
		encoder = json.NewEncoder(out)
		encoder.(*json.Encoder).SetIndent("", "  ")
	case types.YamlEncoding:
		encoder = newYamlEncoder(out)
	case types.YamlDocumentsEncoding:
		encoder = &yamldocsEncoder{out: out}
	case types.TomlEncoding:
		encoder = toml.NewEncoder(out)
		encoder.(*toml.Encoder).Indent = "  "
	default:
		encoder = &rawEncoder{out: out}
	}

	return encoder.Encode(data)
}

type rawEncoder struct {
	out io.Writer
}

func (e *rawEncoder) Encode(v any) error {
	_, err := fmt.Fprintln(e.out, v)
	return err
}

type yamldocsEncoder struct {
	out io.Writer
}

func (e *yamldocsEncoder) Encode(data any) error {
	rValue := reflect.ValueOf(data)
	rType := reflect.TypeOf(data)
	if rType.Kind() == reflect.Pointer {
		rValue = rValue.Elem()
		rType = rValue.Type()
	}

	encoder := newYamlEncoder(e.out)

	switch rType.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < rValue.Len(); i++ {
			value := rValue.Index(i).Interface()
			if err := encoder.Encode(value); err != nil {
				return err
			}
		}

		return nil
	}

	return encoder.Encode(data)
}
