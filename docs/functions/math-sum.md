# +

`+` sums up all provided arguments. Arguments must evaluate to numeric values.
`sum` is an alias for this function.

## Examples

* `(+ 1 2 3)` -> `6`
* `(+ 1 1.5)` -> `2.5`

## Forms

### `(+ expr+)`

* `expr` is 1 or more expressions

`+` evaluates each of the given expressions in sequence. If the expression
evaluates to a number, it is added to the total sum. If an expression returns
an error, `+` returns that error and stops evaluating further expressions.

## Context

`+` uses one scope per expression, so nothing is shared (like variables) between
expressions, and nothing is leaking out.
