// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.xrstf.de/rudi/cmd/rudi/encoding"
	"go.xrstf.de/rudi/cmd/rudi/options"
	"go.xrstf.de/rudi/cmd/rudi/types"
)

func LoadFiles(opts *options.Options, filenames []string) ([]any, error) {
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

func LoadFile(opts *options.Options, filename string) (any, error) {
	if filename == "" {
		return nil, errors.New("no filename provided")
	}

	if filename == "-" {
		return encoding.Decode(os.Stdin, opts.StdinFormat)
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return encoding.Decode(f, getFileFormat(filename))
}

func getFileFormat(filename string) types.Encoding {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".json":
		return types.JsonEncoding
	case ".tml", ".toml":
		return types.TomlEncoding
	default:
		return types.YamlEncoding
	}
}
