// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package pathexpr

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/jsonpath"
	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

func Traverse(ctx types.Context, value any, path *ast.PathExpression) (any, error) {
	if path == nil {
		return value, nil
	}

	jp, err := ToJSONPath(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("invalid path expression: %w", err)
	}

	return jsonpath.Get(value, jp)
}

func ToJSONPath(ctx types.Context, path *ast.PathExpression) (jsonpath.Path, error) {
	p := jsonpath.Path{}
	for i := range path.Steps {
		step := path.Steps[i]

		if step.Filter != nil {
			p = append(p, &filterStep{
				ctx:  ctx,
				expr: step.Filter,
			})
		} else if step.Expression != nil {
			value, err := ctx.Runtime().EvalExpression(ctx, step.Expression)
			if err != nil {
				return nil, err
			}

			p = append(p, &singleStep{
				coalescer: ctx.Coalesce(),
				value:     value,
			})
		} else {
			return nil, errors.New("invalid path step")
		}
	}

	return p, nil
}

type singleStep struct {
	coalescer coalescing.Coalescer
	value     any
}

var _ jsonpath.SingleStep = &singleStep{}

func (fs *singleStep) ToIndex() (int, bool) {
	index, err := fs.coalescer.ToInt64(fs.value)
	if err != nil {
		return 0, false
	}

	return int(index), true
}

func (fs *singleStep) ToKey() (string, bool) {
	key, err := fs.coalescer.ToString(fs.value)
	if err != nil {
		return "", false
	}

	return key, true
}

type filterStep struct {
	ctx  types.Context
	expr ast.Expression
}

var _ jsonpath.FilterStep = &filterStep{}

func (fs *filterStep) Keep(key any, value any) (bool, error) {
	vars := types.NewVariables().Set("key", key)

	doc, err := types.NewDocument(value)
	if err != nil {
		return false, fmt.Errorf("cannot use %T in filter expressions", value)
	}

	stepCtx := fs.ctx.NewShallowScope(&doc, vars)

	result, err := fs.ctx.Runtime().EvalExpression(stepCtx, fs.expr)
	if err != nil {
		return false, err
	}

	converted, err := fs.ctx.Coalesce().ToBool(result)
	if err != nil {
		return false, fmt.Errorf("expression result %T cannot be converted to bool", result)
	}

	return converted, nil
}
