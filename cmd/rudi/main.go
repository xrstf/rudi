// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/pflag"
	"go.xrstf.de/rudi"
	"go.xrstf.de/rudi/cmd/rudi/cmd/console"
	"go.xrstf.de/rudi/cmd/rudi/cmd/help"
	"go.xrstf.de/rudi/cmd/rudi/cmd/script"
	"go.xrstf.de/rudi/cmd/rudi/types"
)

// These variables get set by ldflags during compilation.
var (
	BuildTag    string
	BuildCommit string
	BuildDate   string // RFC3339 format ("2006-01-02T15:04:05Z07:00")
)

func printVersion() {
	fmt.Printf(
		"Rudi %s (%s), built with %s on %s\n",
		BuildTag,
		BuildCommit[:10],
		runtime.Version(),
		BuildDate,
	)
}

func main() {
	opts := types.Options{}

	opts.AddFlags(pflag.CommandLine)
	pflag.Parse()

	if opts.ShowVersion {
		printVersion()
		return
	}

	if err := opts.Validate(); err != nil {
		fmt.Printf("Invalid command line: %v\n", err)
		os.Exit(2)
	}

	args := pflag.Args()

	if opts.ShowHelp || (len(args) > 0 && args[0] == "help") {
		if err := help.Run(&opts, args); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		return
	}

	if opts.Interactive || len(args) == 0 {
		if err := console.Run(&opts, args); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		return
	}

	if err := script.Run(&opts, args); err != nil {
		parseErr := &rudi.ParseError{}
		if errors.As(err, parseErr) {
			fmt.Println(parseErr.Snippet())
			fmt.Println(parseErr)
		} else {
			fmt.Printf("Error: %v\n", err)
		}

		os.Exit(1)
	}
}
