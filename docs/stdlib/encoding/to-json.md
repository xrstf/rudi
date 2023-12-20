# to-json

This function encodes a value as JSON.

## Examples

* `(to-json {foo 23})` ➜ `"{\"foo\":23}"`
* `(to-json null)` ➜ `"null"`

## Forms

### `(to-json value:any)` ➜ `string`

This is the only form of this function. It encodes a value as JSON. If encoding
fails, an error is thrown.
