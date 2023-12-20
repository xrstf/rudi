# try

`try` is used to catch errors in expressions. It will evaluate an expression and
return its return value, unless the expression errors out, then a fallback value
is returned.

`try` is similar to `default`, but only returns the fallback value if an error
occurs, compared to `default` which tests for empty-ishness.

## Examples

* `(try (+ 1 2) "fallback")` ➜ `3` (no error occurred in `+`)
* `(try (+ 1 "invalid") "fallback")` ➜ `"fallback"`
* `(try (+ 1 "invalid") (+ "also invalid"))` ➜ error

## Forms

### `(try candidate:expression)` ➜ `any`

This is equivalent to `(try candidate null)`.

### `(try candidate:expression fallback:expression)` ➜ `any`

* `candidate` is an arbitrary expression.
* `fallback` is an arbitrary expression.

`try` evaluates the candidate expression and returns its return value upon
success. However when the candidate return an error, the fallback expression is
evaluated and its return value (or error) are returned.

## Context

`try` executes candidate and fallback in their own scopes, so variables from
either expression are not visible in the other and neither leak outside of `try`.
