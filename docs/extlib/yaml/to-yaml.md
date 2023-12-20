# to-yaml

This function encodes a value as YAML.

## Examples

* `(to-yaml {foo 23})` ➜ `"foo: 23\n"`
* `(to-yaml null)` ➜ `"null\n"`

## Forms

### `(to-yaml value:any)` ➜ `string`

This is the only form of this function. It encodes a value as YAML. If encoding
fails, an error is thrown.
