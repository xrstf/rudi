package docs

import (
	"embed"
	_ "embed"
	"fmt"
	"strings"

	"go.xrstf.de/rudi/pkg/eval/builtin"
)

//go:embed *
var embeddedFS embed.FS

type Topic struct {
	CliNames    []string
	Group       string
	Description string
	filename    string
}

func (t *Topic) Content() ([]byte, error) {
	return embeddedFS.ReadFile(t.filename)
}

func Topics() []Topic {
	topics := []Topic{
		{
			CliNames:    []string{"language", "lang", "rudi"},
			Group:       "General",
			Description: "A short introduction to the Rudi language",
			filename:    "language.md",
		},
	}

	for funcName := range builtin.Functions {
		var sanitized string
		switch funcName {
		case "+":
			sanitized = "sum"
		case "-":
			sanitized = "minus"
		case "*":
			sanitized = "multiply"
		case "/":
			sanitized = "divide"
		default:
			sanitized = strings.ReplaceAll(funcName, "?", "")
			sanitized = strings.ReplaceAll(sanitized, "!", "")
		}

		group, filename, found := getFunctionDocs(sanitized)
		if !found {
			continue
		}

		topics = append(topics, Topic{
			CliNames:    []string{funcName, sanitized},
			Group:       group + " Functions",
			Description: fmt.Sprintf("documentation for the %s function", funcName),
			filename:    filename,
		})
	}

	return topics
}

// Function names are global in Rudi; however the docs is logically split
// into groups like "core" or "math", which also make sense in the documentation
// (hence why names are like "core-if.md").
// In order not to introduce any sort of weird grouping/namespacing in the
// actual eval packages, we deduce the correct file by searching through _all_
// function files and search for *-[funcName].md
func getFunctionDocs(sanitizedFuncName string) (group string, filename string, found bool) {
	const functionsDir = "functions"

	entries, err := embeddedFS.ReadDir(functionsDir)
	if err != nil {
		panic(err)
	}

	suffix := fmt.Sprintf("-%s.md", sanitizedFuncName)

	for _, entry := range entries {
		filename := entry.Name()

		if strings.HasSuffix(filename, suffix) {
			group := strings.TrimSuffix(filename, suffix)

			return group, functionsDir + "/" + filename, true
		}
	}

	return "", "", false
}
