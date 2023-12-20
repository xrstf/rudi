// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"go.xrstf.de/rudi/cmd/rudi/batteries"
	"go.xrstf.de/rudi/pkg/docs"
)

const (
	rudiReadme   = "../../docs/README.md"
	stdlibReadme = "../../docs/stdlib/README.md"
	extlibReadme = "../../docs/extlib/README.md"
	consoleHelp  = "../../cmd/rudi/cmd/console/help.md"
)

func main() {
	if err := updateFile(rudiReadme, ""); err != nil {
		log.Fatalf("Failed to update %s: %v", rudiReadme, err)
	}

	if err := updateFile(stdlibReadme, "../"); err != nil {
		log.Fatalf("Failed to update %s: %v", stdlibReadme, err)
	}

	if err := updateFile(extlibReadme, "../"); err != nil {
		log.Fatalf("Failed to update %s: %v", extlibReadme, err)
	}

	if err := updateFile(consoleHelp, "../"); err != nil {
		log.Fatalf("Failed to update %s: %v", consoleHelp, err)
	}
}

func updateFile(filename string, linkPrefix string) error {
	// read current file
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filename, err)
	}
	body := string(content)

	// inject standard topics TOC
	body = inject(body, renderTopicsTOC(linkPrefix), "TOPICS")

	// inject help topics TOC (like the other topics, but lists names instead of titles)
	body = inject(body, renderHelpTopicsTOC(), "HELP_TOPICS")

	// inject stdlib TOC
	builtInModules := []docs.Module{}
	builtInModules = append(builtInModules, batteries.SafeBuiltInModules...)
	builtInModules = append(builtInModules, batteries.UnsafeBuiltInModules...)
	body = inject(body, renderLibraryTOC(builtInModules, linkPrefix+"stdlib/"), "STDLIB")

	// inject extlib TOC
	body = inject(body, renderLibraryTOC(batteries.ExtendedModules, linkPrefix+"extlib/"), "EXTLIB")

	// inject help lib TOC
	allModules := builtInModules
	allModules = append(allModules, batteries.ExtendedModules...)
	body = inject(body, renderHelpLibraryTOC(allModules), "HELP_LIB")

	// write updated file
	if err := os.WriteFile(filename, []byte(body), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}

	return nil
}

func inject(body string, injected string, marker string) string {
	beginMarker := fmt.Sprintf("<!-- BEGIN_%s_TOC -->", marker)
	endMarker := fmt.Sprintf("<!-- END_%s_TOC -->", marker)

	fullReplacement := fmt.Sprintf("%s\n%s%s", beginMarker, injected, endMarker)
	regex := regexp.MustCompile(`(?s)` + beginMarker + `.+` + endMarker)

	return regex.ReplaceAllString(body, fullReplacement)
}
