// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.xrstf.de/rudi/cmd/rudi/batteries"
	"go.xrstf.de/rudi/hack/docs-prerender/ansidoc"
	"go.xrstf.de/rudi/pkg/docs"
)

const (
	cmdsDirectory     = "../../cmd/rudi/cmd"
	docsDirectory     = "../../docs"
	embeddedDirectory = "../../cmd/rudi/docs/data"
)

func main() {
	// Dump docs for std/ext library modules, both to make processing all docs
	// easier (just have to implement 1 way to do it: walk the filesystem in docs/)
	// and to have easily browsable docs in GitHub in one central repository.
	dumpModules(batteries.SafeBuiltInModules, "stdlib", true)
	dumpModules(batteries.UnsafeBuiltInModules, "stdlib", false)
	dumpModules(batteries.ExtendedModules, "extlib", true)

	// find all Markdown files, except for READMEs
	docsFiles := map[string]string{}
	functionFiles := []string{}

	if err := filepath.WalkDir(docsDirectory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if name := d.Name(); name == "README.md" || filepath.Ext(name) != ".md" {
			return nil
		}

		relPath, _ := filepath.Rel(docsDirectory, path)

		if isFunctionDoc(relPath) {
			functionFiles = append(functionFiles, path)
		} else {
			docsFiles[path] = filepath.Base(path)
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	if err := filepath.WalkDir(cmdsDirectory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(d.Name()) != ".md" {
			return nil
		}

		relPath, _ := filepath.Rel(cmdsDirectory, path)
		command := filepath.Dir(relPath)

		docsFiles[path] = fmt.Sprintf("cmd-%s.md", command)

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	if err := os.RemoveAll(embeddedDirectory); err != nil {
		log.Fatal(err)
	}

	// Prerender each file and turn the original Markdown into Rudimark, a pre-computed
	// terminal-rendered Markdown, but with placeholders for actual color information.

	// For normal docs files, we follow the structure as given.
	for source, dest := range docsFiles {
		prerenderFile(source, dest)
	}

	// Functions, which are not namespaced in Rudi, are all dumped into the same directory,
	// making it easier to find them later when the user does "rudi help concat".
	for _, file := range functionFiles {
		prerenderFile(file, filepath.Join("functions", filepath.Base(file)))
	}

	log.Println("Done.")
}

func isFunctionDoc(path string) bool {
	return strings.HasPrefix(path, "stdlib"+string(filepath.Separator)) || strings.HasPrefix(path, "extlib"+string(filepath.Separator))
}

func dumpModules(mods []docs.Module, library string, wipe bool) {
	libDir := filepath.Join(docsDirectory, library)

	if wipe {
		// delete everything but the README.md
		entries, err := os.ReadDir(libDir)
		if err != nil {
			log.Fatal(err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				if err := os.RemoveAll(filepath.Join(libDir, entry.Name())); err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	for _, mod := range mods {
		// log.Printf("Dumping docs for module %s…", mod.Name)

		modDirName := filepath.Join(libDir, mod.Name)

		if err := os.MkdirAll(modDirName, 0755); err != nil {
			log.Fatalf("Failed to create %s: %v", modDirName, err)
		}

		for funcName := range mod.Functions {
			doc, err := mod.Documentation.Documentation(funcName)
			if err != nil {
				log.Printf("Warning: failed to get docs for %s/%s: %v", mod.Name, funcName, err)
				continue
			}

			// doc = strings.TrimSpace(doc)

			// Ensure that each function uses a filename that is derived from its function name,
			// regardless from where the original function doc provider (mod.Documentation) actually
			// got the Markdown from.
			// This is important so that the Rudi interpreter later can easily find the function just
			// based on its name alone.
			filename := fmt.Sprintf("%s/%s.md", modDirName, docs.Normalize(funcName))

			if err := os.WriteFile(filename, []byte(doc), 0644); err != nil {
				log.Fatalf("Failed to write %s: %v", filename, err)
			}
		}
	}
}

func prerenderFile(source string, dest string) {
	// log.Printf("Prerendering %s…", source)

	fullDest := filepath.Join(embeddedDirectory, dest) + ".gz"
	destDir := filepath.Dir(fullDest)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		log.Fatal(err)
	}

	content, err := os.ReadFile(source)
	if err != nil {
		log.Fatal(err)
	}

	prerendered, err := prerender(string(content))
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(fullDest)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	gzipwriter := gzip.NewWriter(f)
	defer gzipwriter.Close()

	io.WriteString(gzipwriter, prerendered)
}

func prerender(markdown string) (string, error) {
	rendered := renderMarkdown(markdown, 0)

	parsed, err := ansidoc.Parse(rendered)
	if err != nil {
		return "", err
	}

	parsed = ansidoc.Optimize(parsed)
	templated := ansidoc.Templatify(parsed, rendered)

	return templated, nil
}
