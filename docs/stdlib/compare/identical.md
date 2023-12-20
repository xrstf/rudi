# identical?

`identical?` compares two arguments for equalities. Equality is defined using the
[strict coalescer](../../coalescing.md#strict-coalescer), so `(identical? 1 "1")`
yields a type error because strings cannot be converted to numbers.

See also [`eq?`](eq.md), which uses the currently selected coalescer,
and [`like?`](like.md), which always uses humane coalescing.

## Examples

* `(identical? 1 "1")` ➜ `false`
* `(identical? 1 1.0)` ➜ `true`

## Forms

### `(identical? left:any right:any)` ➜ `bool`

* `left` is an arbitrary expression, except for identifiers.
* `right` is likewise an arbitrary expression, except for identifiers.

Both expressions are evaluated and then compared using the current coalescer.
If evaluation of either of the expressions yields and error, that error is
returned.

Equality is not defined for all type combinations and so `identical?` can return
errors for invalid comparisons. See the coalescing documentation for conversion
rules.

## Context

`identical?` executes both expressions in their own contexts, so nothing is
shared.
