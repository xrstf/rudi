# to-lower

`to-lower` removes a copy of a string with all uppercase characters replaced
with their lowercase equivalent. This function uses Go's regular strings
package and so this function should not be used where Unicode characters are
involved, unexpected results might happen.

See also [`to-upper`](strings-to-upper.md).

## Examples

* `(to-lower "FOO")` -> `"foo"`

## Forms

### `(to-lower string)`

* `string` is an arbitrary expression.

`to-lower` evaluates the first argument and coalesces the result into a string.
When successful, the lowercased version of the string is returned.
