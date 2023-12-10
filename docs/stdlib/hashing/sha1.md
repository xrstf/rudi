# sha1

`sha1` returns the lowercase hex string presenration of the SHA-1 hash of the
given input value.

## Examples

* `(sha1 "")` -> `"da39a3ee5e6b4b0d3255bfef95601890afd80709"`
* `(sha1 "hello")` -> `"aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"`

## Forms

### `(sha1 data)`

* `data` is an arbitrary expression.

`sha1` evaluates the given expression and coalesces the result to a
string. If either of those steps fail, an error is returned. Otherwise the
function will calculate the SHA-1 hash and return the hex presentation (a string
of 40 characters).
