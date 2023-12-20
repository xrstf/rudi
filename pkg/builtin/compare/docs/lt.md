# lt?

`lt?` returns whether the first argument is numerically smaller than the second.

## Examples

* `(lt? 3 2)` ➜ `false`
* `(lt? 2.0 3)` ➜ `true`

## Forms

### `(lt? left:any right:any)` ➜ `bool`

* `left` is an arbitrary expression, except for identifiers.
* `right` is likewise an arbitrary expression, except for identifiers.

`lt?` evaluates both expressions and coalesces their results to be numbers. If
either evaluation or conversion fail, an error is returned. If the two arguments
are valid numbers, the result of `left < right` is returned.

## Context

`lt?` executes both expressions in their own contexts, so nothing is shared.
