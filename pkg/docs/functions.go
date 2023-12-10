package docs

import (
	"embed"
	"strings"
)

type FunctionProvider interface {
	// Documentation returns the Markdown documentation for a function. If no
	// documentation is available, an error shall be returned.
	Documentation(funcName string) (string, error)
}

type funcProvider struct {
	fs *embed.FS
}

var _ FunctionProvider = &funcProvider{}

func NewFunctionProvider(fs *embed.FS) FunctionProvider {
	return &funcProvider{
		fs: fs,
	}
}

func (p funcProvider) Documentation(funcName string) (string, error) {
	contents, err := p.fs.ReadFile(Normalize(funcName) + ".md")
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func Normalize(funcName string) string {
	funcName = strings.ReplaceAll(funcName, "?", "")
	funcName = strings.ReplaceAll(funcName, "!", "")

	return funcName
}
