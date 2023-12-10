# to-int

`to-int` converts the given value to an integer number. The function always uses
[humane coalescing](../../coalescing.md#humane-coalescer).

## Examples

* `(to-int "0")` -> `0`
* `(to-int 1.0)` -> `1`
* `(to-int 1.2)` -> error, no lossless conversion possible
* `(to-int [])` -> error

## Forms

### `(to-int value)`

* `value` is an arbitrary expression.

`to-int` evaluates the given expression and then coalesces the result into an
integer value. See the documentation for the humane coalescer for the exact
conversion rules.
