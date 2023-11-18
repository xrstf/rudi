// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"go.xrstf.de/otto"
	"gopkg.in/yaml.v3"
)

func loadFiles(opts *options, filenames []string) ([]any, error) {
	results := make([]any, len(filenames))

	for i, filename := range filenames {
		data, err := loadFile(opts, filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", filename, err)
		}

		results[i] = data
	}

	return results, nil
}

func loadFile(opts *options, filename string) (any, error) {
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

func setupOttoContext(files []any) (otto.Context, error) {
	var (
		document otto.Document
		err      error
	)

	if len(files) > 0 {
		document, err = otto.NewDocument(files[0])
		if err != nil {
			return otto.Context{}, fmt.Errorf("cannot use first input as document: %w", err)
		}
	} else {
		document, _ = otto.NewDocument(nil)
	}

	vars := otto.NewVariables().
		Set("files", files)

	ctx := otto.NewContext(document, otto.NewBuiltInFunctions(), vars)

	return ctx, nil
}
