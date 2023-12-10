# to-upper

`to-upper` removes a copy of a string with all lowercase characters replaced
with their uppercase equivalent. This function uses Go's regular strings
package and so this function should not be used where Unicode characters are
involved, unexpected results might happen.

See also [`to-lower`](to-lower.md).

## Examples

* `(to-upper "foo")` -> `"FOO"`

## Forms

### `(to-upper string)`

* `string` is an arbitrary expression.

`to-upper` evaluates the first argument and coalesces the result into a string.
When successful, the uppercased version of the string is returned.
