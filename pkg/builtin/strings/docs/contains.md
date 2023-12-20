# contains?

`contains?` returns `true` if the given haystack value contains the given needle
value.

`contains?` will first attempt to coalesce the first argument as a vector, and
fallback to string coalescing otherwise. This means depending on the coalescer
the behaviour of this function can change when the first argument is neither
vector nor string.

## Examples

* `(contains? ["foo"] "bar")` ➜ `false`
* `(contains? "foo" "f")` ➜ `true`
* `(contains? "9000" 9)` ➜ `true` with humane coalescing, error otherwise

## Forms

### `(contains? haystack:vector needle:any)` ➜ `bool`

* `haystack` is an arbitrary expression.
* `needle` is an arbitrary expression.

For vector checks, the needle is evaluated and compared to each element of the
haystack vector, using the current coalescer's equality rules. If an element was
found that is equal to the needle, `true` is returned, otherwise `false`.

### `(contains? haystack:string needle:string)` ➜ `bool`

* `haystack` is an arbitrary expression.
* `needle` is an arbitrary expression.

For string checks, the needle is evaluated and coalesces to a string. If
successful, the function returns `true` if the needle string is contained in the
haystack string, otherwise `false`.

## Context

`contains?` executes all expressions in their own contexts, so nothing is shared.
