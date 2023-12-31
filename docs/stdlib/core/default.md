# default

`default` is used to apply default values, especially when reading data from
user-defined inputs. Because of this, `default` uses a loose definition when
determining emptyness by using the `empty?` function internally. This means
values like `0` or `""` are considered empty.

`default` is similar to `try`, but only returns the fallback value if the value
is empty-ish. If an error occurs, `default` does not fall back to the fallback
value, but returns the error instead.

`default` is a shortcut to writing `(if (empty? expr-a) expr-b expr-a)`

## Examples

* `(default "" "fallback")` ➜ `"fallback"`
* `(default "set" "fallback")` ➜ `"set"`
* `(default (+ "invalid") "fallback")` ➜ error

The function is also nice when combined with the bang modifier to apply default
values to a variable:

* `(default! $var.foo "default-value")`

## Forms

### `(default candidate:any fallback:any)` ➜ `any`

* `candidate` is an arbitrary expression.
* `fallback` is an arbitrary expression.

`default` evaluates the candidate expression and returns the evaluated fallback
value if the returned candidate value is empty-ish. If the candidate expression
returns an error, the error is returned and the fallback expressions is not
evaluated.

## Context

`default` executes candidate and fallback in their own scopes, so variables from
either expression are not visible in the other and neither leak outside of
`default`.
