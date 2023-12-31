// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	"go.xrstf.de/rudi"
	"go.xrstf.de/rudi/cmd/rudi/batteries"
	"go.xrstf.de/rudi/cmd/rudi/cmd/console"
	"go.xrstf.de/rudi/cmd/rudi/cmd/help"
	"go.xrstf.de/rudi/cmd/rudi/cmd/script"
	"go.xrstf.de/rudi/cmd/rudi/options"
	"go.xrstf.de/rudi/cmd/rudi/util"

	"github.com/spf13/pflag"
)

// These variables get set by ldflags during compilation.
var (
	BuildTag    string
	BuildCommit string
	BuildDate   string // RFC3339 format ("2006-01-02T15:04:05Z07:00")
)

type moduleVersion struct {
	module  string
	version string
}

func printVersion() {
	fmt.Printf(
		"Rudi %s (%s), built with %s on %s\n",
		BuildTag,
		BuildCommit[:10],
		runtime.Version(),
		BuildDate,
	)

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	extlibVersions := []moduleVersion{}
	maxModuleLength := 0

	for _, extMod := range batteries.ExtendedModules {
		// This should never happen.
		if extMod.GoModule == "" {
			continue
		}

		for _, dep := range info.Deps {
			if dep.Path == extMod.GoModule {
				if len(dep.Path) > maxModuleLength {
					maxModuleLength = len(dep.Path)
				}

				extlibVersions = append(extlibVersions, moduleVersion{
					module:  dep.Path,
					version: dep.Version,
				})
			}
		}
	}

	if len(extlibVersions) == 0 {
		return
	}

	fmt.Println()
	fmt.Println("Extended Library:")
	fmt.Println()

	format := fmt.Sprintf("%%-%ds %%s\n", maxModuleLength)

	for _, v := range extlibVersions {
		fmt.Printf(format, v.module, v.version)
	}
}

func main() {
	opts := options.NewDefaultOptions()

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

	handler := util.SetupSignalHandler()

	if opts.Interactive || len(args) == 0 {
		if err := console.Run(handler, &opts, args, BuildTag); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		return
	}

	if err := script.Run(handler, &opts, args); err != nil {
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
