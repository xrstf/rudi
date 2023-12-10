# /

`/` returns the quotient of dividing all arguments. Arguments must evaluate to
numeric values. `div` is an alias for this function.

To prevent ambiguity, this function always performs floating point divisions,
regardless if all its arguments are integer numbers.

## Examples

* `(/ 9 3 2)` -> `1.5` ((9.0 / 3.0) / 2.0)
* `(/ 1 0)` -> invalid: division by zero

## Forms

### `(/ expr+)`

* `expr` is 1 or more expressions

`/` evaluates each of the given expressions in sequence. If an expression returns
an error, `/` returns that error and stops evaluating further expressions.

The first value is taken as the dividend, every further value is then used as
a divisor. The final result is then returned.

## Context

`/` uses one scope per expression, so nothing is shared (like variables) between
expressions, and nothing is leaking out.
