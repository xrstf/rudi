// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package docs

type Node int

const (
	DocumentNode Node = iota + 1
	BlockQuoteNode
	ParagraphNode
	ListNode
	HeadingNode
	H1Node
	H2Node
	H3Node
	H4Node
	H5Node
	H6Node
	StrikethroughNode
	EmphNode
	StrongNode
	HorizontalRuleNode
	ItemNode
	EnumerationNode
	TaskNode
	LinkNode
	LinkTextNode
	ImageNode
	ImageTextNode
	CodeNode
	CodeBlockNode
	TableNode
	DdNode
)

// These constants mirror what Glamour is configuring when rendering a Markdown,
// not necessarily _all_ possible chroma nodes.
// The value are not a sequence, but instead are valid values from chroma's
// ttyIndexed formatter, which will try to find the closest matching color for
// a given hex color when rendering, but this can lead to non-cistent results
// as this matching works by iterating over a loop. To prevent this from happening,
// we must choose colors that map perfectly on the tty256 palette.
// Randomly 0xd7 was chosen as the common prefix for all chroma constants,
// because enough colors are available in the palette with this prefix. As long
// as the values are clearly different from the regular Node values, their
// values do not really matter.

const (
	ChromaText                Node = 0xd70000
	ChromaError               Node = 0xd7005f
	ChromaComment             Node = 0xd70087
	ChromaCommentPreproc      Node = 0xd700af
	ChromaKeyword             Node = 0xd700d7
	ChromaKeywordReserved     Node = 0xd700ff
	ChromaKeywordNamespace    Node = 0xd75f00
	ChromaKeywordType         Node = 0xd75f5f
	ChromaOperator            Node = 0xd75f87
	ChromaPunctuation         Node = 0xd75faf
	ChromaName                Node = 0xd75fd7
	ChromaNameBuiltin         Node = 0xd75fff
	ChromaNameTag             Node = 0xd78700
	ChromaNameAttribute       Node = 0xd7875f
	ChromaNameClass           Node = 0xd78787
	ChromaNameDecorator       Node = 0xd787af
	ChromaNameFunction        Node = 0xd787d7
	ChromaLiteralNumber       Node = 0xd787ff
	ChromaLiteralString       Node = 0xd7af00
	ChromaLiteralStringEscape Node = 0xd7af5f
	ChromaGenericDeleted      Node = 0xd7af87
	ChromaGenericEmph         Node = 0xd7afaf
	ChromaGenericInserted     Node = 0xd7afd7
	ChromaGenericStrong       Node = 0xd7afff
	ChromaGenericSubheading   Node = 0xd7d700
	ChromaBackground          Node = 0xd7d75f

	// not configured by Chroma, so not configured by us either
	// ChromaNameConstant  Node = 0xd7d787
	// ChromaNameException Node = 0xd7d7af
	// ChromaNameOther     Node = 0xd7d7d7
	// ChromaLiteral       Node = 0xd7d7ff
	// ChromaLiteralDate   Node = 0xd7ff00.
)
