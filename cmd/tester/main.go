// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval"
	"go.xrstf.de/corel/pkg/lang/parser"
)

func main() {
	filename := "test.corel"
	if flag.NArg() > 0 {
		filename = flag.Arg(0)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", filename, err)
	}

	docData, err := os.ReadFile("document.json")
	if err != nil {
		log.Fatalf("Failed to read test document: %v", err)
	}

	var document interface{}
	if err := json.Unmarshal(docData, &document); err != nil {
		log.Fatalf("Failed to parse test document: %v", err)
	}

	got, err := parser.Parse(filename, content, parser.Debug(false))
	if err != nil {
		log.Fatal(caretError(err, string(content)))
	}

	program, ok := got.(ast.Program)
	if !ok {
		log.Fatalf("Parsed result is not a ast.Program, but %T", got)
	}

	fmt.Println("---[ INPUT ]-----------------------------------------")
	fmt.Println(string(content))
	// fmt.Println("---[ AST ]-------------------------------------------")
	// fmt.Printf("%#v\n", program)
	// fmt.Println("---[ PRINTED ]---------------------------------------")
	// fmt.Println(program.String())
	fmt.Println("---[ DOCUMENT ]--------------------------------------")
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(document)
	fmt.Println("---[ EVALUATED ]-------------------------------------")

	progContext := eval.NewContext(eval.Document{
		Data: document,
	}, eval.NewVariables().Set("global", document))

	fmt.Println(eval.Run(progContext, &program))
	fmt.Println("-----------------------------------------------------")
}

func caretError(err error, input string) string {
	if el, ok := err.(parser.ErrorLister); ok {
		var buffer bytes.Buffer
		for _, e := range el.Errors() {
			if parserErr, ok := e.(parser.ParserError); ok {
				_, col, off := parserErr.Pos()
				line := extractLine(input, off)
				if col >= len(line) {
					col = len(line) - 1
				} else {
					if col > 0 {
						col--
					}
				}
				if col < 0 {
					col = 0
				}
				pos := col
				for _, chr := range line[:col] {
					if chr == '\t' {
						pos += 7
					}
				}
				buffer.WriteString(fmt.Sprintf("%s\n%s\n%s\n", line, strings.Repeat(" ", pos)+"^", err.Error()))
			} else {
				return err.Error()
			}
		}
		return buffer.String()
	}
	return err.Error()
}

func extractLine(input string, initPos int) string {
	if initPos < 0 {
		initPos = 0
	}
	if initPos >= len(input) && len(input) > 0 {
		initPos = len(input) - 1
	}
	startPos := initPos
	endPos := initPos
	for ; startPos > 0; startPos-- {
		if input[startPos] == '\n' {
			if startPos != initPos {
				startPos++
				break
			}
		}
	}
	for ; endPos < len(input); endPos++ {
		if input[endPos] == '\n' {
			if endPos == initPos {
				endPos++
			}
			break
		}
	}
	return input[startPos:endPos]
}
