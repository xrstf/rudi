// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package util

import (
	"fmt"

	"go.xrstf.de/rudi"
	"go.xrstf.de/rudi/cmd/rudi/batteries"
	"go.xrstf.de/rudi/cmd/rudi/options"
	"go.xrstf.de/rudi/cmd/rudi/types"
	"go.xrstf.de/rudi/pkg/coalescing"
)

func SetupRudiContext(opts *options.Options, fileNames []string, fileContents []any) (rudi.Context, error) {
	var (
		document rudi.Document
		err      error
	)

	if len(fileContents) > 0 {
		document, err = rudi.NewDocument(fileContents[0])
		if err != nil {
			return rudi.Context{}, fmt.Errorf("cannot use first input as document: %w", err)
		}
	} else {
		document, _ = rudi.NewDocument(nil)
	}

	vars := rudi.NewVariables()
	for k, v := range opts.ExtraVariables {
		vars.Set(k, v)
	}

	// system-defined variables come last, so they override anything user specified
	vars.
		Set("files", fileContents).
		Set("filenames", fileNames)

	var coalescer coalescing.Coalescer
	switch opts.Coalescing {
	case types.StrictCoalescing:
		coalescer = coalescing.NewStrict()
	case types.PedanticCoalescing:
		coalescer = coalescing.NewPedantic()
	case types.HumaneCoalescing:
		coalescer = coalescing.NewHumane()
	default:
		return rudi.Context{}, fmt.Errorf("unknown coalescing mode %q", opts.Coalescing)
	}

	funcs := rudi.NewFunctions()

	for _, mod := range batteries.SafeBuiltInModules {
		funcs.Add(mod.Functions)
	}

	// Only add rudispace function support when explicitly enabled, as defining functions at
	// runtime can lead to non-terminating programs and resource exhaustion.
	if opts.EnableRudispaceFunctions {
		funcs.Add(batteries.RudifuncModule.Functions)
	}

	for _, mod := range batteries.ExtendedModules {
		funcs.Add(mod.Functions)
	}

	// No context set here, caller is expected to provide their own (the Rudi context is re-used
	// in the console, but the Go context should not be, hence the separation).
	//nolint:staticcheck
	return rudi.NewContext(nil, document, vars, funcs, coalescer), nil
}
