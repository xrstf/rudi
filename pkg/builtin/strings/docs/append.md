# append

`append` adds additional items to a vector or concatenates a string at the end
of another one, depending on the given arguments. See also this functions
cousin, [`prepend`](prepend.md).

`append` will append to a vector if the first argument coalesces to a vector;
otherwise string coalescing is attempted. This means depending on the coalescer
the behaviour of this function can change when the first argument is neither
vector nor string.

## Examples

* `(append ["foo"] "bar" "x" 3)` -> `["foo" "bar" "x" 3]`
* `(append "foo" "bar" "baz")` -> `"foobarbaz"`
* `(append "foo" 2)` -> `"foo2"` with humane coalescing, error otherwise
* `(append 2 3 4)` -> `"234"` with humane coalescing, error otherwise
* `(append null 1)` -> `[1]` because `null` can turn into empty vectors

## Forms

### `(append base appends+)`

* `base` is an arbitrary expression.
* `appends` are one or more additional arbitrary expressions.

`append` evaluates the base expression first. If the result coalesces to a
vector, all further arguments will be appended to the vector. Otherwise, string
coalescing is attempted. If the argument coalesces to neither, an error is
returned.

All further `appends` expressions are then evaluated and appended. For vectors,
the types of the values does not matter, as vectors can hold any kind of data.
For string append mode, each of the further arguments must coalesce to strings,
otherwise an error is returned.

## Context

`append` executes all expressions in their own contexts, so nothing is shared.
