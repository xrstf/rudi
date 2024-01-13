// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import "go.xrstf.de/rudi/pkg/lang/ast"

type Path []Step

func (p Path) IsValid() bool {
	for _, s := range p {
		switch s.(type) {
		case SingularVectorStep,
			SingularObjectStep,
			MultiVectorStep,
			MultiObjectStep:
			continue
		default:
			return false
		}
	}

	return true
}

func (p Path) HasMultiSteps() bool {
	for _, s := range p {
		if isMultiStep(s) {
			return true
		}
	}

	return false
}

func isMultiStep(s Step) bool {
	switch s.(type) {
	case MultiVectorStep, MultiObjectStep:
		return true
	default:
		return false
	}
}

type Step any

type SingularVectorStep interface {
	Index() int
}

type IndexStep int

func (i IndexStep) Index() int {
	return int(i)
}

type SingularObjectStep interface {
	Key() string
}

type KeyStep string

func (k KeyStep) Key() string {
	return string(k)
}

type MultiVectorStep interface {
	Keep(index int, value any) (bool, error)
}

type MultiObjectStep interface {
	Keep(key string, value any) (bool, error)
}

func FromEvaluatedPath(evaledPath ast.EvaluatedPathExpression) Path {
	p := Path{}
	for _, step := range evaledPath.Steps {
		if i := step.IntegerValue; i != nil {
			p = append(p, IndexStep(*i))
		} else if s := step.StringValue; s != nil {
			p = append(p, KeyStep(*s))
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
