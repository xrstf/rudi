# like?

`like?` compares two arguments for equalities. Equality is defined using the
[humane coalescer](../../coalescing.md#humane-coalescer), so `(like? 1 "1")` yields
no type error, but true.

See also [`eq?`](eq.md), which uses the currently selected coalescer,
and [`identical?`](identical.md), which always uses strict
coalescing.

## Examples

* `(like? 1 "1")` -> `true`
* `(like? 1 "2")` -> `false`

## Forms

### `(like? left right)`

* `left` is an arbitrary expression, except for identifiers.
* `right` is likewise an arbitrary expression, except for identifiers.

Both expressions are evaluated and then compared using the current coalescer.
If evaluation of either of the expressions yields and error, that error is
returned.

Equality is not defined for all type combinations, even with the humane coalescer,
and so `like?` can return errors for invalid comparisons. See the coalescing
documentation for conversion rules.

## Context

`like?` executes both expressions in their own contexts, so nothing is shared.
