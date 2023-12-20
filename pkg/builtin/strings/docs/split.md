# split

`split` splits a string using a separator into smaller substrings. It returns
a vector of these substrings.

## Examples

* `(split "" "")` ➜ `[]`
* `(split "," "")` ➜ `[""]`
* `(split "," "a,b,1")` ➜ `["a", "b", "1"]`
* `(split "" "hello")` ➜ `["h", "e", "l", "l", "o"]`

## Forms

### `(split separator:string value:string)` ➜ `vector`

* `separator` is an arbitrary expression that evaluates to a string.
* `value` is an arbitrary expression that evaluates to a string.

`split` splits the `value` and returns a vector with the substrings.

## Context

`split` executes separator and vector in their own contexts, so nothing is
shared.
