// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package docs

import "go.xrstf.de/rudi/pkg/eval/types"

// A module combines functions with their documentation. This type and concept only exist in
// the Rudi interpreter (cmd/rudi); the actual codebase separates the function code from their
// documentation into distinct packages, so that embedded Rudi usage does not need to embed all
// the documentation as well.
type Module struct {
	Name          string
	Functions     types.Functions
	Documentation FunctionProvider
}
