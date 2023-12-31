// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package docs

var (
	darkStyle = map[Node]TextStyle{
		DocumentNode: {
			Foreground: "252",
		},
		BlockQuoteNode: {},
		ParagraphNode: {
			Foreground: "252",
		},
		ListNode: {},
		HeadingNode: {
			Foreground: "39",
			Effects:    Bold,
		},
		H1Node: {
			Foreground: "228",
			Background: "63",
			Effects:    Bold,
		},
		H2Node: {},
		H3Node: {},
		H4Node: {},
		H5Node: {},
		H6Node: {
			Foreground: "35",
			Effects:    0,
		},
		StrikethroughNode: {
			Effects: Strikethrough,
		},
		EmphNode: {
			Effects: Italic,
		},
		StrongNode: {
			Effects: Bold,
		},
		HorizontalRuleNode: {
			Foreground: "240",
		},
		ItemNode:        {},
		EnumerationNode: {},
		TaskNode:        {},
		LinkNode: {
			Foreground: "30",
			Effects:    Underlined,
		},
		LinkTextNode: {
			Foreground: "35",
			Effects:    Bold,
		},
		ImageNode: {
			Foreground: "212",
		},
		ImageTextNode: {
			Foreground: "243",
		},
		CodeNode: {
			Foreground: "203",
			Background: "236",
		},
		CodeBlockNode: {
			Foreground: "244",
		},
		TableNode: {},
		DdNode:    {},

		// chroma

		ChromaText: {
			Foreground: "C4C4C4",
		},
		ChromaError: {
			Foreground: "F1F1F1",
			Background: "F05B5B",
		},
		ChromaComment: {
			Foreground: "676767",
		},
		ChromaCommentPreproc: {
			Foreground: "FF875F",
		},
		ChromaKeyword: {
			Foreground: "00AAFF",
		},
		ChromaKeywordReserved: {
			Foreground: "FF5FD2",
		},
		ChromaKeywordNamespace: {
			Foreground: "FF5F87",
		},
		ChromaKeywordType: {
			Foreground: "6E6ED8",
		},
		ChromaOperator: {
			Foreground: "EF8080",
		},
		ChromaPunctuation: {
			Foreground: "E8E8A8",
		},
		ChromaName: {
			Foreground: "C4C4C4",
		},
		ChromaNameBuiltin: {
			Foreground: "FF8EC7",
		},
		ChromaNameTag: {
			Foreground: "B083EA",
		},
		ChromaNameAttribute: {
			Foreground: "7A7AE6",
		},
		ChromaNameClass: {
			Foreground: "F1F1F1",
			Effects:    Underlined | Bold,
		},
		ChromaNameDecorator: {
			Foreground: "FFFF87",
		},
		ChromaNameFunction: {
			Foreground: "00D787",
		},
		ChromaLiteralNumber: {
			Foreground: "6EEFC0",
		},
		ChromaLiteralString: {
			Foreground: "C69669",
		},
		ChromaLiteralStringEscape: {
			Foreground: "AFFFD7",
		},
		ChromaGenericDeleted: {
			Foreground: "FD5B5B",
		},
		ChromaGenericEmph: {
			Effects: Italic,
		},
		ChromaGenericInserted: {
			Foreground: "00D787",
		},
		ChromaGenericStrong: {
			Effects: Bold,
		},
		ChromaGenericSubheading: {
			Foreground: "777777",
		},
		ChromaBackground: {
			Background: "373737",
		},
	}

	lightStyle = map[Node]TextStyle{
		DocumentNode: {
			Foreground: "234",
		},
		BlockQuoteNode: {},
		ParagraphNode: {
			Foreground: "234",
		},
		ListNode: {},
		HeadingNode: {
			Foreground: "27",
			Effects:    Bold,
		},
		H1Node: {
			Foreground: "228",
			Background: "63",
			Effects:    Bold,
		},
		H2Node: {},
		H3Node: {},
		H4Node: {},
		H5Node: {},
		H6Node: {
			Foreground: "35",
			Effects:    0,
		},
		StrikethroughNode: {
			Effects: Strikethrough,
		},
		EmphNode: {
			Effects: Italic,
		},
		StrongNode: {
			Effects: Bold,
		},
		HorizontalRuleNode: {
			Foreground: "249",
		},
		ItemNode:        {},
		EnumerationNode: {},
		TaskNode:        {},
		LinkNode: {
			Foreground: "36",
			Effects:    Underlined,
		},
		LinkTextNode: {
			Foreground: "29",
			Effects:    Bold,
		},
		ImageNode: {
			Foreground: "205",
		},
		ImageTextNode: {
			Foreground: "243",
		},
		CodeNode: {
			Foreground: "203",
			Background: "254",
		},
		CodeBlockNode: {
			Foreground: "242",
		},
		TableNode: {},
		DdNode:    {},

		// chroma

		ChromaText: {
			Foreground: "2A2A2A",
		},
		ChromaError: {
			Foreground: "F1F1F1",
			Background: "FF5555",
		},
		ChromaComment: {
			Foreground: "8D8D8D",
		},
		ChromaCommentPreproc: {
			Foreground: "FF875F",
		},
		ChromaKeyword: {
			Foreground: "279EFC",
		},
		ChromaKeywordReserved: {
			Foreground: "FF5FD2",
		},
		ChromaKeywordNamespace: {
			Foreground: "FB406F",
		},
		ChromaKeywordType: {
			Foreground: "7049C2",
		},
		ChromaOperator: {
			Foreground: "FF2626",
		},
		ChromaPunctuation: {
			Foreground: "FA7878",
		},
		ChromaName: {
			Foreground: "C4C4C4",
		},
		ChromaNameBuiltin: {
			Foreground: "C4C4C4",
		},
		ChromaNameTag: {
			Foreground: "581290",
		},
		ChromaNameAttribute: {
			Foreground: "8362CB",
		},
		ChromaNameClass: {
			Foreground: "212121",
			Effects:    Underlined | Bold,
		},
		ChromaNameDecorator: {
			Foreground: "A3A322",
		},
		ChromaNameFunction: {
			Foreground: "019F57",
		},
		ChromaLiteralNumber: {
			Foreground: "22CCAE",
		},
		ChromaLiteralString: {
			Foreground: "7E5B38",
		},
		ChromaLiteralStringEscape: {
			Foreground: "00AEAE",
		},
		ChromaGenericDeleted: {
			Foreground: "FD5B5B",
		},
		ChromaGenericEmph: {
			Effects: Italic,
		},
		ChromaGenericInserted: {
			Foreground: "00D787",
		},
		ChromaGenericStrong: {
			Effects: Bold,
		},
		ChromaGenericSubheading: {
			Foreground: "777777",
		},
		ChromaBackground: {
			Background: "373737",
		},
	}
)
