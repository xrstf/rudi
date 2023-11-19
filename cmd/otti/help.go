// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"strings"

	markdown "go.xrstf.de/go-term-markdown"
	"go.xrstf.de/otto/docs"
)

func renderHelpTopic(helpTopics []docs.Topic, selectedTopic string) error {
	for _, topic := range helpTopics {
		for _, cliName := range topic.CliNames {
			if strings.EqualFold(cliName, selectedTopic) {
				content, err := topic.Content()
				if err != nil {
					return fmt.Errorf("failed to render docs: %w", err)
				}

				printMarkdown(string(content))
				return nil
			}
		}
	}

	return fmt.Errorf("no help available for %q", selectedTopic)
}

func printMarkdown(markup string) {
	result := markdown.Render(markup, 80, 2)
	fmt.Println()
	fmt.Println(string(result))
}
