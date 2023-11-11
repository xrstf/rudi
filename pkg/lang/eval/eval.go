package eval

import (
	"go.xrstf.de/corel/pkg/lang/ast"
)

type Object struct {
	Data map[string]interface{}
}

func Run(ctx Context, p *ast.Program) (interface{}, error) {
	result, err := evalProgram(ctx, p)
	if err != nil {
		return nil, err
	}

	return result, nil
}
