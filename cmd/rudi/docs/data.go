// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package docs

import (
	"compress/gzip"
	"embed"
	_ "embed"
	"fmt"
	"io"

	"go.xrstf.de/rudi/pkg/docs"
)

var Aliases = map[string]string{
	// math module

	"+": "add",
	"-": "sub",
	"*": "mult",
	"/": "div",
}

//go:embed data/*
var embeddedFS embed.FS

func ReadFile(filename string) (string, error) {
	filename = "data/" + filename + ".gz"

	f, err := embeddedFS.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	defer r.Close()

	content, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func RenderFile(filename string, painter PainterFunc) (string, error) {
	content, err := ReadFile(filename)
	if err != nil {
		return "", err
	}

	return Render(content, painter), nil
}

func ReadFunction(funcName string) (string, error) {
	realName, ok := Aliases[funcName]
	if ok {
		funcName = realName
	}

	return ReadFile(fmt.Sprintf("functions/%s.md", docs.Normalize(funcName)))
}

func RenderFunction(funcName string, painter PainterFunc) (string, error) {
	content, err := ReadFunction(funcName)
	if err != nil {
		return "", err
	}

	return Render(content, painter), nil
}

type Topic struct {
	Title       string
	CliNames    []string
	Description string
	Filename    string // in the format ReadFile expects
}

func Topics() []Topic {
	// Does not include the cmd/* documentation, as those are only for --help flag handling.
	return []Topic{
		{
			Title:       "The Rudi Language",
			CliNames:    []string{"language", "lang", "rudi"},
			Description: "A short introduction to the Rudi language",
			Filename:    "language.md",
		},
		{
			Title:       "Type Handling & Conversions",
			CliNames:    []string{"coalescing"},
			Description: "How Rudi handles, converts and compares values",
			Filename:    "coalescing.md",
		},
	}
}

func (t *Topic) Content() (string, error) {
	return ReadFile(t.Filename)
}

func (t *Topic) Render(painter PainterFunc) (string, error) {
	content, err := t.Content()
	if err != nil {
		return "", err
	}

	return Render(content, painter), nil
}
