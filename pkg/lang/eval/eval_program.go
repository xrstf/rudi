package eval

import (
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func evalProgram(ctx Context, p *ast.Program) (interface{}, error) {
	finalCtx := ctx

	// This is all sorts of wonky and not really how the program execution should work.
	// But it compiles.
	var (
		result interface{}
		err    error
	)

	for i, stmt := range p.Statements {
		finalCtx, result, err = evalStatement(finalCtx, &p.Statements[i])
		if err != nil {
			return nil, fmt.Errorf("failed to eval statement %s: %w", stmt.String(), err)
		}
	}

	return result, nil
}
