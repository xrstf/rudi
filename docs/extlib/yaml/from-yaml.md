# from-yaml

This function decodes a YAML string into a Go datastructure.

## Examples

* `(from-yaml "foo: 23")` ➜ `{"foo" 23}`
* `(from-yaml "~")` ➜ `null`

## Forms

### `(from-yaml markup:string)` ➜ `any`

This is the only form of this function. It decodes a YAML string and returns the
result. If invalid YAML is provided, an error is thrown.
