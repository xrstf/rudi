// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package util

import (
	"errors"
	"fmt"
	"io"
	"os"

	"go.xrstf.de/rudi/cmd/rudi/types"

	"gopkg.in/yaml.v3"
)

func LoadFiles(opts *types.Options, filenames []string) ([]any, error) {
	results := make([]any, len(filenames))

	for i, filename := range filenames {
		data, err := LoadFile(opts, filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", filename, err)
		}

		results[i] = data
	}

	return results, nil
}

func LoadFile(opts *types.Options, filename string) (any, error) {
	if filename == "" {
		return nil, errors.New("no filename provided")
	}

	var input io.Reader

	if filename == "-" {
		input = os.Stdin
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		input = f
	}

	var doc any

	decoder := yaml.NewDecoder(input)
	if err := decoder.Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to parse document as YAML/JSON: %w", err)
	}

	return doc, nil
}
