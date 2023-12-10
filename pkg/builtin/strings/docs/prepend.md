# prepend

`prepend` adds additional items to a vector or concatenates a string at the end
of another one, depending on the given arguments. See also this functions
cousin, [`append`](append.md).

`prepend` will prepend to a vector if the first argument coalesces to a vector;
otherwise string coalescing is attempted. This means depending on the coalescer
the behaviour of this function can change when the first argument is neither
vector nor string.

When prepending multiple elements, they are all prepended as a single list, not
individually, so `(prepend [1] 2 3)` yields `[2 3 1]`, not `[3 2 1]`.

## Examples

* `(prepend ["foo"] "bar" "x" 3)` -> `["x" 3 "foo" "bar"]`
* `(prepend "foo" "bar" "baz")` -> `"barbazfoo"`
* `(prepend "foo" 2)` -> `"2foo"` with humane coalescing, error otherwise
* `(prepend 2 3 4)` -> `"342"` with humane coalescing, error otherwise
* `(prepend null 1)` -> `[1]` because `null` can turn into empty vectors

## Forms

### `(prepend base prepends+)`

* `base` is an arbitrary expression.
* `prepends` are one or more additional arbitrary expressions.

`prepend` evaluates the base expression first. If the result coalesces to a
vector, all further arguments will be prepended to the vector. Otherwise, string
coalescing is attempted. If the argument coalesces to neither, an error is
returned.

All further `prepends` expressions are then evaluated and prepended. For vectors,
the types of the values does not matter, as vectors can hold any kind of data.
For string prepend mode, each of the further arguments must coalesce to strings,
otherwise an error is returned.

## Context

`prepend` executes all expressions in their own contexts, so nothing is shared.
