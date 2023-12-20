# to-lower

`to-lower` removes a copy of a string with all uppercase characters replaced
with their lowercase equivalent. This function uses Go's regular strings
package and so this function should not be used where Unicode characters are
involved, unexpected results might happen.

See also [`to-upper`](to-upper.md).

## Examples

* `(to-lower "FOO")` ➜ `"foo"`

## Forms

### `(to-lower string:string)` ➜ `string`

* `string` is an arbitrary expression.

Returns a copy of the value with all bytes being turned to their lowercase
equivalent.
