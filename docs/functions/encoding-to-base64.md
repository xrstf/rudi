# to-base64

`to-base64` encodes string with base64. See also the inverse function
[`from-base64`](encoding-from-base64.md).

## Examples

* `(to-base64 "")` -> `""`
* `(to-base64 "hello")` -> `"aGVsbG8="`

## Forms

### `(to-base64 data)`

* `data` is an arbitrary expression.

`to-base64` evaluates the given expression and coalesces the result to a
string. If either of those steps fail, an error is returned. Otherwise the
function will encode the string with base64 and return the result.
