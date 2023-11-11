package eval

import (
	"go.xrstf.de/corel/pkg/lang/ast"
)

type Object struct {
	Data map[string]interface{}
}

func Run(p *ast.Program, rootObject Object) (interface{}, error) {
	result, err := evalProgram(p, &rootObject)
	if err != nil {
		return nil, err
	}

	return result, nil
}
