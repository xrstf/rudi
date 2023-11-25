// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	html "html/template"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"go.xrstf.de/rudi/docs"
)

const (
	filename    = "docs/README.md"
	beginMarker = `<!-- BEGIN_TOC -->`
	endMarker   = `<!-- END_TOC -->`
)

func main() {
	topics := docs.Topics()
	groups := getGroups(topics)

	rendered := renderTopics(topics, groups)
	rendered = fmt.Sprintf("%s\n%s\n%s", beginMarker, rendered, endMarker)

	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", filename, err)
	}

	regex := regexp.MustCompile(`(?s)` + beginMarker + `.+` + endMarker)
	output := regex.ReplaceAllString(string(content), rendered)

	os.WriteFile(filename, []byte(output), 0644)
}

func strSliceHas(haystack []string, needle string) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}

	return false
}

func getGroups(topics []docs.Topic) []string {
	// determine a sorted list of functions, with some groups
	// hardcoded to be at the top, regardless of their name
	prioritizedGroups := []string{
		"General",
		"Core Functions",
	}

	remainingGroups := []string{}

	for _, topic := range topics {
		if !strSliceHas(prioritizedGroups, topic.Group) {
			if !strSliceHas(remainingGroups, topic.Group) {
				remainingGroups = append(remainingGroups, topic.Group)
			}
		}
	}

	sort.Strings(remainingGroups)

	return append(prioritizedGroups, remainingGroups...)
}

func renderTopics(topics []docs.Topic, groups []string) string {
	var out strings.Builder

	for _, group := range groups {
		out.WriteString(fmt.Sprintf("## %s\n", group))
		out.WriteString("\n")

		topicNames := getTopicNames(topics, group)
		for _, topicName := range topicNames {
			topic := getTopic(topics, topicName)
			linkTitle := topicName

			if topic.IsFunction {
				linkTitle = fmt.Sprintf("`%s`", linkTitle)
			}

			out.WriteString(fmt.Sprintf("* [%s](%s) – %s\n", linkTitle, topic.Filename, topic.Description))
		}

		out.WriteString("\n")
	}

	return strings.TrimSpace(out.String())
}

func htmlencode(s string) string {
	return html.HTMLEscapeString(s)
}

func getTopicNames(topics []docs.Topic, group string) []string {
	names := []string{}

	for _, topic := range topics {
		if topic.Group != group {
			continue
		}

		primaryName := topic.CliNames[0]

		if !strSliceHas(names, primaryName) {
			names = append(names, primaryName)
		}
	}

	sort.Strings(names)

	return names
}

func getTopic(topics []docs.Topic, name string) docs.Topic {
	for _, topic := range topics {
		if topic.CliNames[0] == name {
			return topic
		}
	}

	panic("this should never happen")
}
