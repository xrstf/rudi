// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package batteries

import (
	coalescemod "go.xrstf.de/rudi/pkg/builtin/coalesce"
	coalescedocs "go.xrstf.de/rudi/pkg/builtin/coalesce/docs"
	comparemod "go.xrstf.de/rudi/pkg/builtin/compare"
	comparedocs "go.xrstf.de/rudi/pkg/builtin/compare/docs"
	coremod "go.xrstf.de/rudi/pkg/builtin/core"
	coredocs "go.xrstf.de/rudi/pkg/builtin/core/docs"
	datetimemod "go.xrstf.de/rudi/pkg/builtin/datetime"
	datetimedocs "go.xrstf.de/rudi/pkg/builtin/datetime/docs"
	encodingmod "go.xrstf.de/rudi/pkg/builtin/encoding"
	encodingdocs "go.xrstf.de/rudi/pkg/builtin/encoding/docs"
	hashingmod "go.xrstf.de/rudi/pkg/builtin/hashing"
	hashingdocs "go.xrstf.de/rudi/pkg/builtin/hashing/docs"
	listsmod "go.xrstf.de/rudi/pkg/builtin/lists"
	listsdocs "go.xrstf.de/rudi/pkg/builtin/lists/docs"
	logicmod "go.xrstf.de/rudi/pkg/builtin/logic"
	logicdocs "go.xrstf.de/rudi/pkg/builtin/logic/docs"
	mathmod "go.xrstf.de/rudi/pkg/builtin/math"
	mathdocs "go.xrstf.de/rudi/pkg/builtin/math/docs"
	stringsmod "go.xrstf.de/rudi/pkg/builtin/strings"
	stringsdocs "go.xrstf.de/rudi/pkg/builtin/strings/docs"
	typesmod "go.xrstf.de/rudi/pkg/builtin/types"
	typesdocs "go.xrstf.de/rudi/pkg/builtin/types/docs"
	"go.xrstf.de/rudi/pkg/docs"

	semvermod "go.xrstf.de/rudi-contrib/semver"
	semverdocs "go.xrstf.de/rudi-contrib/semver/docs"
	uuidmod "go.xrstf.de/rudi-contrib/uuid"
	uuiddocs "go.xrstf.de/rudi-contrib/uuid/docs"
	yamlmod "go.xrstf.de/rudi-contrib/yaml"
	yamldocs "go.xrstf.de/rudi-contrib/yaml/docs"
)

var (
	// BuiltInModules look alphabetically sorted, except that "core" is always the first item,
	// because it's the most important module and should be first in the documentation. Order here
	// does not matter otherwise anyway.
	BuiltInModules = []docs.Module{
		{
			Name:          "core",
			Functions:     coremod.Functions,
			Documentation: coredocs.Functions,
		},
		{
			Name:          "coalesce",
			Functions:     coalescemod.Functions,
			Documentation: coalescedocs.Functions,
		},
		{
			Name:          "compare",
			Functions:     comparemod.Functions,
			Documentation: comparedocs.Functions,
		},
		{
			Name:          "datetime",
			Functions:     datetimemod.Functions,
			Documentation: datetimedocs.Functions,
		},
		{
			Name:          "encoding",
			Functions:     encodingmod.Functions,
			Documentation: encodingdocs.Functions,
		},
		{
			Name:          "hashing",
			Functions:     hashingmod.Functions,
			Documentation: hashingdocs.Functions,
		},
		{
			Name:          "lists",
			Functions:     listsmod.Functions,
			Documentation: listsdocs.Functions,
		},
		{
			Name:          "logic",
			Functions:     logicmod.Functions,
			Documentation: logicdocs.Functions,
		},
		{
			Name:          "math",
			Functions:     mathmod.Functions,
			Documentation: mathdocs.Functions,
		},
		{
			Name:          "strings",
			Functions:     stringsmod.Functions,
			Documentation: stringsdocs.Functions,
		},
		{
			Name:          "types",
			Functions:     typesmod.Functions,
			Documentation: typesdocs.Functions,
		},
	}

	ExtendedModules = []docs.Module{
		{
			Name:          "semver",
			Functions:     semvermod.Functions,
			Documentation: semverdocs.Functions,
		},
		{
			Name:          "uuid",
			Functions:     uuidmod.Functions,
			Documentation: uuiddocs.Functions,
		},
		{
			Name:          "yaml",
			Functions:     yamlmod.Functions,
			Documentation: yamldocs.Functions,
		},
	}
)
