# to-string

`to-string` converts the given value to a string. The function always uses
[humane coalescing](../../coalescing.md#humane-coalescer).

## Examples

* `(to-string null)` -> `""`
* `(to-string true)` -> `"true"`
* `(to-string false)` -> `"false"`
* `(to-string 0)` -> `"0"`
* `(to-string 1.1)` -> `"1.1"`
* `(to-string [])` -> error

## Forms

### `(to-string value)`

* `value` is an arbitrary expression.

`to-string` evaluates the given expression and then coalesces the result into a
string value. See the documentation for the humane coalescer for the exact
conversion rules.
