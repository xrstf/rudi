# do

`do` evaluates all expressions in sequence, sharing a context between them. This
effectively forms a sub program and is useful for combining with other functions
that require exactly 1 expression.

`do` purposefully makes its child expressions have side effects, so that setting
variables has an effect on subsequent expressions.

## Examples

* `(do true 42)` -> `42`
* `(do (set $foo 1) (+ $foo 2))` -> `3`

## Forms

### `(do expr+)`

* `expr` is 1 or more expressions

`do` evaluates sets up a new context and then evaluates all `expr` in sequence,
sharing the context between them, forming a sub program. The return value of `do`
is the return value of the last expression in it.

When any expression encounters an error, `do` stops evaluation and returns the
error.

## Context

`do` shares a context between all child expressions.
