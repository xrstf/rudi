# contains?

`contains?` returns `true` if the given haystack value contains the given needle
value.

`contains?` will first attempt to coalesce the first argument as a vector, and
fallback to string coalescing otherwise. This means depending on the coalescer
the behaviour of this function can change when the first argument is neither
vector nor string.

## Examples

* `(contains? ["foo"] "bar")` -> `false`
* `(contains? "foo" "f")` -> `true`
* `(contains? "9000" 9)` -> `true` with humane coalescing, error otherwise

## Forms

### `(contains? haystack needle)`

* `haystack` is an arbitrary expression.
* `needle` is an arbitrary expression.

`contains?` evaluates the haystack expression first. If the result coalesces to
a vector, the function will check if the vector contains the needle. Otherwise
string coalescing is attempted and if it succeeds, the function will check if
the needle is a substring of the hackstack. If the haystack is neither string
nor vector, an error is returned.

For vector checks, the needle is evaluated and compared to each element of the
haystack vector, using the current coalescer's equality rules. If an element was
found that is equal to the needle, `true` is returned, otherwise `false`.

For string checks, the needle is evaluated and coalesces to a string. If
successful, the function returns `true` if the needle string is contained in the
haystack string, otherwise `false`.

## Context

`contains?` executes all expressions in their own contexts, so nothing is shared.
