# eq?

`eq?` compares two arguments for equalities. Equality is defined based on the
currently active [coalescer](../coalescing.md), so for example with strict
coalescing, `(eq 1 "1")` is false, but with humane coalescing is true.

See also [`like?`](comparisons-like.md), which always uses humane coalescing,
and [`identical?`](comparisons-identical.md), which always uses strict
coalescing.

## Examples

* `(eq? "" "")` -> `true`
* `(eq? 1 2)` -> `false`
* `(eq? (+ 1 1) 2)` -> `true`

## Forms

### `(eq? left right)`

* `left` is an arbitrary expression, except for identifiers.
* `right` is likewise an arbitrary expression, except for identifiers.

Both expressions are evaluated and then compared using the current coalescer.
If evaluation of either of the expressions yields and error, that error is
returned.

Equality is not defined for all type combinations and so `eq?` can, depending
on the coalescer, return errors for invalid comparisons. With strict coalescing,
`(eq? true 2)` returns an error because numbers cannot be converted to booleans.
With humane coalescing, true would be returned instead and no error is generated.

## Context

`eq?` executes both expressions in their own contexts, so nothing is shared.
