# to-upper

`to-upper` removes a copy of a string with all lowercase characters replaced
with their uppercase equivalent. This function uses Go's regular strings
package and so this function should not be used where Unicode characters are
involved, unexpected results might happen.

See also [`to-lower`](to-lower.md).

## Examples

* `(to-upper "foo")` ➜ `"FOO"`

## Forms

### `(to-upper value:string)` ➜ `string`

* `value` is an arbitrary expression.

Returns a copy of the value with all bytes being turned to their uppercase
equivalent.
