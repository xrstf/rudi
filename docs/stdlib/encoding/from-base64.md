# from-base64

`from-base64` decodes a base64-encoded string. See also the inverse function
[`to-base64`](to-base64.md).

## Examples

* `(from-base64 "")` ➜ `""`
* `(from-base64 "invalid")` ➜ error
* `(from-base64 "aGVsbG8=")` ➜ `"hello"`

## Forms

### `(from-base64 encoded:string)` ➜ `string`

* `encoded` is an arbitrary expression.

`from-base64` evaluates the given expression and coalesces the result to a
string. If either of those steps fail, an error is returned. Otherwise the
function will attempt to decode the value, returning the decoded data on success
and an error otherwise.
