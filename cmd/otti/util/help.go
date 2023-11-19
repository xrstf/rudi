// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package util

import (
	_ "embed"
	"fmt"
	"strings"

	markdown "go.xrstf.de/go-term-markdown"
	"go.xrstf.de/otto/docs"
)

func RenderMarkdown(markup string, indent int) string {
	return string(markdown.Render(markup, 80, indent))
}

func RenderHelpTopics(helpTopics []docs.Topic, indent int) string {
	width := 0
	for _, topic := range helpTopics {
		if l := len(topic.CliNames[0]); l > width {
			width = l
		}
	}

	format := fmt.Sprintf("* %%-%ds – %%s\n", width)

	var builder strings.Builder
	for _, topic := range helpTopics {
		builder.WriteString(fmt.Sprintf(format, topic.CliNames[0], topic.Description))
	}

	return RenderMarkdown(builder.String(), indent)
}

func RenderHelpTopic(helpTopics []docs.Topic, selectedTopic string, indent int) (string, error) {
	for _, topic := range helpTopics {
		for _, cliName := range topic.CliNames {
			if strings.EqualFold(cliName, selectedTopic) {
				content, err := topic.Content()
				if err != nil {
					return "", fmt.Errorf("failed to render docs: %w", err)
				}

				return RenderMarkdown(string(content), indent), nil
			}
		}
	}

	return "", fmt.Errorf("no help available for %q", selectedTopic)
}
