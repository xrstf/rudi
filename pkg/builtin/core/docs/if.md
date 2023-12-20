# if

`if` allows to form conditions. It supports evaluating one expression if another
is true, and allows to optionally evaluate another expression if the
condition was false.

## Examples

* `(if true 42)` ➜ `42`
* `(if false 42)` ➜ `null`
* `(if true "yes" "no")` ➜ `"yes"`
* `(if false "yes" "no")` ➜ `"no"`
* `(if (gt? 4 2) "yes" "no")` ➜ `"yes"`

## Forms

### `(if condition:bool expr:expression)` ➜ `any`

* `condition` is any expression that evaluates to a bool.
* `expr` is any expression.

`if` evaluates the condition and if it returns `true`, the expression is
evaluated and its return value is the final return value of `if`.
If the condition is `false`, `if` returns `null`.

If the condition or expression return an error, `if` returns that error.

### `(if condition:bool expr-a:expression expr-b:expression)` ➜ `any`

* `condition` is any expression that evaluates to a bool.
* `expr-a` is any expression.
* `expr-b` is any expression.

`if` evaluates the condition and if it returns `true`, the `expr-a` is evaluated
and its return value is the final return value of `if`.
If the condition is `false`, `if` evaluates and returns `expr-b`.

`if` guarantees that only one of `expr-a` and `expr-b` is ever evaluated, as
expressions in Rudi can have side effects on the global document.

If `condition`, `expr-a` or `expr-b` return an error, `if` returns that error.

## Context

`if` evaluates all expressions as their own scopes, so if `condition` sets a
variable, this variable is not available in either of the positive / negative
expressions.
