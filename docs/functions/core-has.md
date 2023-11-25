# has?

`has?` returns `true` if the given path expression is valid (regardless of the
value it points to).

## Examples

* `(set! $var 42) (has? $var.foo)` -> `false`
* `(set! $var {foo "bar"}) (has? $var.foo)` -> `true`

## Forms

### `(has? expr)`

* `expr` is exactly 1 expression. This can be any expression that supports
  path expressions (like variables, objects, vectors, tuples).

The value of the expression is evaluated without applying the path expression
at first. If this evaluation results in an error, the error is returned from
`has?`. If the value was successfully computed, the path expression is evaluated
against it. The function then returns whether the path can be traversed
successfully.

## Context

`has?` evaluates the expression in its own scope, so variables defined in it
do not leak.
