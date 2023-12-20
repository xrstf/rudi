# concat

`concat` returns a string of all elements in a vector glued together with a
common string, for example `(concat "x" ["1" "2" "3"]) == "1x2x3"`. The glue
can be left empty, but must be a string. Only vectors are supported for
concatenation.

`concat` accepts more than 1 value for the source of strings to concatenate,
for example `(concat "glue" "a" "b" "c") == "agluebgluec"`. Strings are taken
as they are, vectors are unpacked (not recursively) and must contain only strings.

## Examples

* `(concat "," ["1" "2" "3"])` ➜ `"1,2,3"`
* `(concat "" ["1" "2" "3"])` ➜ `"123"`
* `(concat "," [])` ➜ `""`
* `(concat "," "a" "b")` ➜ `"a,b"`
* `(concat "," "a" ["b" "c"] "d" [] "e")` ➜ `"a,b,c,d,e"`
* `(concat "," "a" [["b"]])` ➜ invalid

## Forms

### `(concat glue:string element:any…)` ➜ `string`

* `glue` is an arbitrary expression that evaluates to a string.
* `element` is 1 or more arbitrary expressions that each evaluate to a string
  or a vector containing only strings.

`concat` combines all elements using the glue string. Each element can be either
a string or a vector. Vectors are only flattened to one level and must contain
only values that coalesce into strings. Empty vectors add nothing to the result.

## Context

`concat` executes glue and each element in their own contexts, so nothing is
shared.
