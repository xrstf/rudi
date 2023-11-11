package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalProgram(p *ast.Program, rootObject *Object) (*Object, error) {
	for i, stmt := range p.Statements {
		newRootObject, err := evalStatement(&p.Statements[i], rootObject)
		if err != nil {
			return nil, fmt.Errorf("failed to eval statement %s: %w", stmt.String(), err)
		}

		rootObject = newRootObject
	}

	return rootObject, nil
}
