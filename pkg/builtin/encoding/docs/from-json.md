# from-json

This function decodes a JSON string into a Go datastructure.

## Examples

* `(from-json "{\"foo\": 23}")` ➜ `{"foo" 23}`
* `(from-json "true")` ➜ `true`

## Forms

### `(from-json markup:string)` ➜ `any`

This is the only form of this function. It decodes a JSON string and returns the
result. If invalid JSON is provided, an error is thrown.
