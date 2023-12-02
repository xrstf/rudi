# to-float

`to-float` converts the given value to a floating point number. The function
always uses [humane coalescing](../coalescing.md#humane-coalescer).

## Examples

* `(to-float 0)` -> `0.0`
* `(to-float "0")` -> `0.0`
* `(to-float 1)` -> `1.0`
* `(to-float [])` -> error

## Forms

### `(to-float value)`

* `value` is an arbitrary expression.

`to-float` evaluates the given expression and then coalesces the result into a
floating point value. See the documentation for the humane coalescer for the
exact conversion rules.
