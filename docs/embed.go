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
	Title       string
	CliNames    []string
	Group       string
	Description string
	IsFunction  bool
	Filename    string
}

func (t *Topic) Content() ([]byte, error) {
	return embeddedFS.ReadFile(t.Filename)
}

func Topics() []Topic {
	topics := []Topic{
		{
			Title:       "The Rudi Language",
			CliNames:    []string{"language", "lang", "rudi"},
			Group:       "General",
			Description: "A short introduction to the Rudi language",
			Filename:    "language.md",
		},
		{
			Title:       "Type Handling & Conversions",
			CliNames:    []string{"coalescing"},
			Group:       "General",
			Description: "A short introduction to the Rudi languageHow Rudi handles, converts and compares values",
			Filename:    "coalescing.md",
		},
	}

	ignoredFunctions := map[string]struct{}{
		"sum":  {},
		"sub":  {},
		"mult": {},
		"div":  {},
	}

	for funcName, function := range builtin.Functions {
		if _, ok := ignoredFunctions[funcName]; ok {
			continue
		}

		var sanitized string
		switch funcName {
		case "+":
			sanitized = "sum"
		case "-":
			sanitized = "sub"
		case "*":
			sanitized = "mult"
		case "/":
			sanitized = "div"
		default:
			sanitized = strings.ReplaceAll(funcName, "?", "")
			sanitized = strings.ReplaceAll(sanitized, "!", "")
		}

		group, filename, found := getFunctionDocs(sanitized)
		if !found {
			continue
		}

		topics = append(topics, Topic{
			Title:       funcName,
			CliNames:    []string{funcName, sanitized},
			Group:       ucFirst(group) + " Functions",
			Description: function.Description(),
			Filename:    filename,
			IsFunction:  true,
		})
	}

	return topics
}

func ucFirst(s string) string {
	if len(s) < 2 {
		return strings.ToUpper(s)
	}

	first := string(s[0])

	return strings.ToUpper(first) + s[1:]
}

// Function names are global in Rudi; however the docs is logically split
// into groups like "core" or "math", which also make sense in the documentation
// (hence why names are like "core-if.md").
// In order not to introduce any sort of weird grouping/namespacing in the
// actual eval packages, we deduce the correct file by searching through _all_
// function files and search for *-[funcName].md.
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
