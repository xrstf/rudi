# *

`*` returns the product of all given arguments. Arguments must evaluate to numeric
values. `mult` is an alias for this function.

## Examples

* `(* 1 2 3)` ➜ `6`
* `(* 1 1.5)` ➜ `1.5`

## Forms

### `(* expr:number…)` ➜ `number`

* `expr` is 1 or more expressions

`*` evaluates each of the given expressions in sequence. If an expression returns
an error, `*` returns that error and stops evaluating further expressions.

All values are multiplied together and the final product is returned.

## Context

`*` uses one scope per expression, so nothing is shared (like variables) between
expressions, and nothing is leaking out.
