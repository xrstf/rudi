// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/spf13/pflag"
	"go.xrstf.de/otto/pkg/lang/eval/types"
	"gopkg.in/yaml.v3"
)

// These variables get set by ldflags during compilation.
var (
	BuildTag    string
	BuildCommit string
	BuildDate   string // RFC3339 format ("2006-01-02T15:04:05Z07:00")
)

func printVersion() {
	fmt.Printf(
		"Otti %s (%s), built with %s on %s\n",
		BuildTag,
		BuildCommit[:10],
		runtime.Version(),
		BuildDate,
	)
}

type options struct {
	interactive bool
	version     bool
}

func (o *options) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVarP(&o.interactive, "interactive", "i", o.interactive, "Start an interactive REPL to run expressions.")
	fs.BoolVarP(&o.version, "version", "V", o.version, "Show version and exit.")
}

func main() {
	opts := options{}

	opts.AddFlags(pflag.CommandLine)
	pflag.Parse()

	if opts.version {
		printVersion()
		return
	}

	args := pflag.Args()

	if opts.interactive {
		if err := replRun(&opts, args); err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
	}
}

func loadDocument(opts *options, filename string) (types.Document, error) {
	if filename == "" {
		return types.Document{}, errors.New("no filename provided")
	}

	var input io.Reader

	if filename == "-" {
		input = os.Stdin
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return types.Document{}, err
		}
		defer f.Close()

		input = f
	}

	var doc any

	decoder := yaml.NewDecoder(input)
	if err := decoder.Decode(&doc); err != nil {
		return types.Document{}, fmt.Errorf("failed to parse document as YAML/JSON: %w", err)
	}

	return types.NewDocument(doc)
}
