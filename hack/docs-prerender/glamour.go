// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"regexp"
	"strings"

	"go.xrstf.de/rudi/cmd/rudi/docs"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
)

func init() {
	// Configure a dummy style which contains fake color information. The colors actually encode
	// the type of node, which is then later used (at runtime) to inject the actual colors.
	glamour.DefaultStyles["synthetic"] = syntheticStyle
}

var markdownLink = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

func renderMarkdown(markup string, indent int) string {
	// turn local links into plain text
	matches := markdownLink.FindAllStringSubmatch(markup, -1)
	for _, match := range matches {
		link := match[2]

		if !strings.Contains(link, "http") {
			markup = strings.ReplaceAll(markup, match[0], match[1])
		}
	}

	r, _ := glamour.NewTermRenderer(
		// actual style is injected later at runtime when rendering the Rudimark
		glamour.WithStandardStyle("synthetic"),
		glamour.WithWordWrap(100),
	)

	rendered, err := r.Render(markup)
	if err != nil {
		panic(err)
	}

	return rendered
}

func ptrTo[T any](v T) *T {
	return &v
}

const defaultListIndent = uint(2)
const defaultMargin = uint(1)
const defaultCodeBlockMargin = uint(2)

func color(element docs.Node) *string {
	return ptrTo(fmt.Sprintf("#%06x", int(element)))
}

var (
	syntheticStyle = &ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "\n",
				BlockSuffix: "\n",
			},
			Margin: ptrTo(defaultMargin),
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: color(docs.BlockQuoteNode),
			},
			Indent:      ptrTo(uint(1)),
			IndentToken: ptrTo("│ "),
		},
		Paragraph: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: color(docs.ParagraphNode),
			},
		},
		List: ansi.StyleList{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Color: color(docs.ListNode),
				},
			},
			LevelIndent: defaultListIndent,
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:       color(docs.HeadingNode),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: " ",
				Suffix: " ",
				Color:  color(docs.H1Node),
			},
		},
		H2: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "## ",
				Color:  color(docs.H2Node),
			},
		},
		H3: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "### ",
				Color:  color(docs.H3Node),
			},
		},
		H4: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "#### ",
				Color:  color(docs.H4Node),
			},
		},
		H5: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "##### ",
				Color:  color(docs.H5Node),
			},
		},
		H6: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "###### ",
				Color:  color(docs.H6Node),
			},
		},
		Strikethrough: ansi.StylePrimitive{
			Color: color(docs.StrikethroughNode),
		},
		Emph: ansi.StylePrimitive{
			Color: color(docs.EmphNode),
		},
		Strong: ansi.StylePrimitive{
			Color: color(docs.StrongNode),
		},
		HorizontalRule: ansi.StylePrimitive{
			Color:  color(docs.HorizontalRuleNode),
			Format: "\n--------\n",
		},
		Item: ansi.StylePrimitive{
			BlockPrefix: "• ",
			Color:       color(docs.ItemNode),
		},
		Enumeration: ansi.StylePrimitive{
			BlockPrefix: ". ",
			Color:       color(docs.EnumerationNode),
		},
		Task: ansi.StyleTask{
			StylePrimitive: ansi.StylePrimitive{
				Color: color(docs.TaskNode),
			},
			Ticked:   "[✓] ",
			Unticked: "[ ] ",
		},
		Link: ansi.StylePrimitive{
			Color: color(docs.LinkNode),
		},
		LinkText: ansi.StylePrimitive{
			Color: color(docs.LinkTextNode),
		},
		Image: ansi.StylePrimitive{
			Color: color(docs.ImageNode),
		},
		ImageText: ansi.StylePrimitive{
			Color:  color(docs.ImageTextNode),
			Format: "Image: {{.text}} →",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: " ",
				Suffix: " ",
				Color:  color(docs.CodeNode),
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Color: color(docs.CodeBlockNode),
				},
				Margin: ptrTo(defaultCodeBlockMargin),
			},
			Chroma: &ansi.Chroma{
				// not configured by Chroma, so not configured by us either
				// NameConstant:  ansi.StylePrimitive{Color: color(docs.ChromaNameConstant)},
				// NameException: ansi.StylePrimitive{Color: color(docs.ChromaNameException)},
				// NameOther:     ansi.StylePrimitive{Color: color(docs.ChromaNameOther)},
				// Literal:       ansi.StylePrimitive{Color: color(docs.ChromaLiteral)},
				// LiteralDate:   ansi.StylePrimitive{Color: color(docs.ChromaLiteralDate)},

				Text:                ansi.StylePrimitive{Color: color(docs.ChromaText)},
				Error:               ansi.StylePrimitive{Color: color(docs.ChromaError)},
				Comment:             ansi.StylePrimitive{Color: color(docs.ChromaComment)},
				CommentPreproc:      ansi.StylePrimitive{Color: color(docs.ChromaCommentPreproc)},
				Keyword:             ansi.StylePrimitive{Color: color(docs.ChromaKeyword)},
				KeywordReserved:     ansi.StylePrimitive{Color: color(docs.ChromaKeywordReserved)},
				KeywordNamespace:    ansi.StylePrimitive{Color: color(docs.ChromaKeywordNamespace)},
				KeywordType:         ansi.StylePrimitive{Color: color(docs.ChromaKeywordType)},
				Operator:            ansi.StylePrimitive{Color: color(docs.ChromaOperator)},
				Punctuation:         ansi.StylePrimitive{Color: color(docs.ChromaPunctuation)},
				Name:                ansi.StylePrimitive{Color: color(docs.ChromaName)},
				NameBuiltin:         ansi.StylePrimitive{Color: color(docs.ChromaNameBuiltin)},
				NameTag:             ansi.StylePrimitive{Color: color(docs.ChromaNameTag)},
				NameAttribute:       ansi.StylePrimitive{Color: color(docs.ChromaNameAttribute)},
				NameClass:           ansi.StylePrimitive{Color: color(docs.ChromaNameClass)},
				NameDecorator:       ansi.StylePrimitive{Color: color(docs.ChromaNameDecorator)},
				NameFunction:        ansi.StylePrimitive{Color: color(docs.ChromaNameFunction)},
				LiteralNumber:       ansi.StylePrimitive{Color: color(docs.ChromaLiteralNumber)},
				LiteralString:       ansi.StylePrimitive{Color: color(docs.ChromaLiteralString)},
				LiteralStringEscape: ansi.StylePrimitive{Color: color(docs.ChromaLiteralStringEscape)},
				GenericDeleted:      ansi.StylePrimitive{Color: color(docs.ChromaGenericDeleted)},
				GenericEmph:         ansi.StylePrimitive{Color: color(docs.ChromaGenericEmph)},
				GenericInserted:     ansi.StylePrimitive{Color: color(docs.ChromaGenericInserted)},
				GenericStrong:       ansi.StylePrimitive{Color: color(docs.ChromaGenericStrong)},
				GenericSubheading:   ansi.StylePrimitive{Color: color(docs.ChromaGenericSubheading)},
				Background:          ansi.StylePrimitive{Color: color(docs.ChromaBackground)},
			},
		},
		Table: ansi.StyleTable{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Color: color(docs.TableNode),
				},
			},
			CenterSeparator: ptrTo("┼"),
			ColumnSeparator: ptrTo("│"),
			RowSeparator:    ptrTo("─"),
		},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: "\n🠶 ",
			Color:       color(docs.DdNode),
		},
	}
)
