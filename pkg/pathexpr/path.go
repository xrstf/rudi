// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package pathexpr

import "go.xrstf.de/rudi/pkg/lang/ast"

type Path []Step
type Step any

func FromEvaluatedPath(evaledPath ast.EvaluatedPathExpression) Path {
	p := Path{}
	for _, step := range evaledPath.Steps {
		if i := step.IntegerValue; i != nil {
			p = append(p, *i)
		} else if s := step.StringValue; s != nil {
			p = append(p, *s)
		}
	}

	return p
}

func toIntegerStep(s Step) (int, bool) {
	switch asserted := s.(type) {
	case int:
		return asserted, true
	case int32:
		return int(asserted), true
	case int64:
		return int(asserted), true
	default:
		return 0, false
	}
}

func toStringStep(s Step) (string, bool) {
	switch asserted := s.(type) {
	case string:
		return asserted, true
	default:
		return "", false
	}
}
