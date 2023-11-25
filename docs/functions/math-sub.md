# -

`-` returns the difference of all given arguments. Arguments must evaluate to
numeric values. `sub` is an alias for this function.

## Examples

* `(- 1 2 3)` -> `-4` ((1 - 2) - 3)
* `(- 1 1.5)` -> `-0.5`

## Forms

### `(- expr+)`

* `expr` is 1 or more expressions

`-` evaluates each of the given expressions in sequence. If an expression returns
an error, `-` returns that error and stops evaluating further expressions.

The first value is taken as the base value, every further value is then subtracted
from the value. The final result is then returned.

## Context

`-` uses one scope per expression, so nothing is shared (like variables) between
expressions, and nothing is leaking out.
