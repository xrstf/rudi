// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"strings"

	"go.xrstf.de/rudi/cmd/rudi/docs"
)

func renderTopicsTOC(linkPrefix string) string {
	topics := docs.Topics()

	var out strings.Builder

	for _, topic := range topics {
		line := fmt.Sprintf("* [%s](%s) – %s\n", topic.Title, linkPrefix+topic.Filename, topic.Description)
		out.WriteString(line)
	}

	return out.String()
}

func renderHelpTopicsTOC() string {
	topics := docs.Topics()

	var out strings.Builder

	for _, topic := range topics {
		line := fmt.Sprintf("* `%s` – %s\n", topic.CliNames[0], topic.Description)
		out.WriteString(line)
	}

	return out.String()
}
