// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package util

import (
	"fmt"

	"go.xrstf.de/rudi"
	"go.xrstf.de/rudi/cmd/rudi/batteries"
	"go.xrstf.de/rudi/cmd/rudi/types"
	"go.xrstf.de/rudi/pkg/coalescing"
)

func SetupRudiContext(opts *types.Options, files []any) (rudi.Context, error) {
	var (
		document rudi.Document
		err      error
	)

	if len(files) > 0 {
		document, err = rudi.NewDocument(files[0])
		if err != nil {
			return rudi.Context{}, fmt.Errorf("cannot use first input as document: %w", err)
		}
	} else {
		document, _ = rudi.NewDocument(nil)
	}

	vars := rudi.NewVariables().
		Set("files", files)

	var coalescer coalescing.Coalescer
	switch opts.Coalescing {
	case "strict":
		coalescer = coalescing.NewStrict()
	case "pedantic":
		coalescer = coalescing.NewPedantic()
	case "humane":
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
