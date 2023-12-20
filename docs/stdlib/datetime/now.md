# now

`now` returns the current date&time, formatted using Go's date formatting logic
(for example, a format string could be `"2006-01-02"`). The time using the
currently configured timezone.

See the [Go documentation](https://pkg.go.dev/time#pkg-constants) for more
information on the format string syntax.

## Examples

* `(now "2006-01-02")` ➜ `"2023-12-02"`

## Forms

### `(now format:string)` ➜ `string`

* `format` is an arbitrary expression.

`now` evaluates the format expression and coalesces the result to a string. If
either of those steps fail, an error is returned. Otherwise the current date &
time are formatted using the given format.
