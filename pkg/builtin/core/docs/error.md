# error

`error` creates a new error and returns it, i.e. this function always "fails".
Errors are constructed using Go's `fmt.Errorf`, so sprintf-style formatting is
available.

## Examples

* `(error "invalid choice")` ➜ error `"invalid choice"`
* `(error "too many replicas: %d" .replicas)` ➜ error `"too man replicas: 3"`

## Forms

### `(error message:string)` ➜ `error`

* `message` is an arbitrary expression.

`error` evaluates the the message and coalesces it to a string. When successful,
a new error with the message is created and returned.

### `(error fmt:string args:any…)` ➜ `error`

* `fmt` is an arbitrary expression.
* `args` are one ore more expressions.

`error` evaluates the the format and coalesces it to a string. When successful,
it evaluates all further arguments and passes their results straight into
`fmt.Errorf` and returns the newly created error.
