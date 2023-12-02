# to-bool

`to-bool` converts the given value to a boolean. The function always uses
[humane coalescing](../coalescing.md#humane-coalescer).

## Examples

* `(to-bool 0)` -> `false`
* `(to-bool "0")` -> `false`
* `(to-bool 1)` -> `true`
* `(to-bool [])` -> `false`
* `(to-bool [0])` -> `true`

## Forms

### `(to-bool value)`

* `value` is an arbitrary expression.

`to-bool` evaluates the given expression and then coalesces the result into a
boolean value. See the documentation for the humane coalescer for the exact
conversion rules.
