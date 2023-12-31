// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package encoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"go.xrstf.de/rudi/cmd/rudi/types"

	"github.com/BurntSushi/toml"
	"github.com/titanous/json5"
	"gopkg.in/yaml.v3"
)

func decodeYaml(input io.Reader) ([]any, error) {
	decoder := yaml.NewDecoder(input)

	documents := []any{}
	for {
		var doc any
		if err := decoder.Decode(&doc); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("failed to parse file as YAML: %w", err)
		}

		documents = append(documents, doc)
	}

	return documents, nil
}

func Decode(input io.Reader, enc types.Encoding) (any, error) {
	var data any

	switch enc {
	case types.RawEncoding:
		content, err := io.ReadAll(input)
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}

		data = string(content)

	case types.JsonEncoding:
		decoder := json.NewDecoder(input)
		if err := decoder.Decode(&data); err != nil {
			return nil, fmt.Errorf("failed to parse file as JSON: %w", err)
		}

	case types.Json5Encoding:
		decoder := json5.NewDecoder(input)
		if err := decoder.Decode(&data); err != nil {
			return nil, fmt.Errorf("failed to parse file as JSON5: %w", err)
		}

	case types.YamlEncoding:
		decoded, err := decodeYaml(input)
		if err != nil {
			return nil, err
		}

		switch len(decoded) {
		case 0:
			data = nil
		case 1:
			data = decoded[0]
		default:
			data = decoded
		}

	case types.YamlDocumentsEncoding:
		decoded, err := decodeYaml(input)
		if err != nil {
			return nil, err
		}

		data = decoded

	case types.TomlEncoding:
		decoder := toml.NewDecoder(input)
		if _, err := decoder.Decode(&data); err != nil {
			return nil, fmt.Errorf("failed to parse file as TOML: %w", err)
		}

	default:
		return nil, fmt.Errorf("unexpected encoding %q", enc)
	}

	return data, nil
}
