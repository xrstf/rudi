// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"sort"
	"strings"

	rudidocs "go.xrstf.de/rudi/cmd/rudi/docs"
	"go.xrstf.de/rudi/pkg/docs"
)

func renderLibraryTOC(lib []docs.Module, linkPrefix string) string {
	var out strings.Builder

	for i, module := range lib {
		out.WriteString(fmt.Sprintf("### %s\n", module.Name))
		out.WriteString("\n")

		functions := []string{}
		for funcName := range module.Functions {
			// Hack: ignore aliases, mostly because the math functions have names
			// that are incompatible with the filesystem (cannot have a "/.md" file).
			if _, exists := rudidocs.Aliases[funcName]; exists {
				continue
			}

			functions = append(functions, funcName)
		}
		sort.Strings(functions)

		for _, funcName := range functions {
			desc := module.Functions[funcName].Description()
			link := fmt.Sprintf("%s%s/%s.md", linkPrefix, module.Name, docs.Normalize(funcName))
			line := fmt.Sprintf("* [`%s`](%s) – %s\n", funcName, link, desc)

			out.WriteString(line)
		}

		if i < len(lib)-1 {
			out.WriteString("\n")
		}
	}

	return out.String()
}

func renderHelpLibraryTOC(lib []docs.Module) string {
	var out strings.Builder

	for i, module := range lib {
		out.WriteString(fmt.Sprintf("* **%s**\n", module.Name))

		functions := []string{}
		for funcName := range module.Functions {
			// Hack: ignore aliases, mostly because the math functions have names
			// that are incompatible with the filesystem (cannot have a "/.md" file).
			if _, exists := rudidocs.Aliases[funcName]; exists {
				continue
			}

			functions = append(functions, funcName)
		}
		sort.Strings(functions)

		for _, funcName := range functions {
			desc := module.Functions[funcName].Description()
			out.WriteString(fmt.Sprintf("  * `%s` – %s\n", funcName, desc))
		}

		if i < len(lib)-1 {
			out.WriteString("\n")
		}
	}

	return out.String()
}
