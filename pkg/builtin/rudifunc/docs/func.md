# func!

This function can be used to define new functions at runtime.

## Unsafe Note

This function is considered unsafe because it allows the creation of programs
like:

```
(func! foo [] (bar))
(func! bar [] (foo))
(foo)
```

which will never terminate, breaking one of Rudi's promises. Because of this,
`func!` is not enabled by default and needs to be enabled using `--enable-funcs`
when using the Rudi interpreter or by including the `rudifunc` module when
embedding Rudi into other Go applications.

## Usage

`func!` creates a new function, which can be used in all subseqeuent statements
in the same Rudi program. Function creation is bound to its parent scope, so
defining a function within an `if` statement for example will make the function
available onside inside that statement, not globally.

Dynamically defined functions cannot overwrite statically defined functions, i.e.
you cannot define a custom `concat` function using `func!` if `concat` already
exists (which is highly likely, given that it's a safe built-in function).
Dynamically defined functions can be overwritten (redefined) though.

Since defining a new function is a side effect, `func` must always be used with
the bang modifier (`func!`). The behaviour of Rudi programs that use `func`
without the bang modifier is undefined.

## Examples

```
(func! myconcat [list] (concat "-" $list))
(myconcat ["1" "2" "3"]) # yields "1-2-3"
```

```
(func! fib [n]
  (if (lte? $n 1)
    0
    (if (eq? $n 2)
      1
      (+ (fib (- $n 1)) (fib (- $n 2))))))
```

Note that Rudi is not optimized for performance and so the function above is
very slow and will quickly exhaust space and time. If you need performance,
inject a Go function statically into Rudi instead of defining it at runtime.

## Forms

### `(func! name params body)`

* `name` is an identifier giving the function its name.
* `params` is a vector containing identifiers that hold the parameter names.
* `body` is a single expression (use `do` for multiple statements) that forms
  the function body.

This form will create a new function called `name` with as many parameters as
`params` has identifiers. `params` can be empty, but must otherwise contain only
unique identifiers.

`name` must be a bare identifier without the bang modifier (`!`). The bang
modifier is a Rudi runtime functionality that functions get by default, so
defining a function `foo` will automatically make `foo!` available as well, with
the same behaviour as built-in functions, for example:

```
(func! inc [n] (+ $n 1))
(set! $a 1)

(inc $a) # yields 2
$a       # yields 1

(inc! $a) # yields 2
$a        # yields 2
```

The newly defined function can then be called like any other built-in function
in Rudi, for example the above function `inc` can be used with
[`map`](../lists/map.md):

```
(func! inc [n] (+ $n 1))
(map [1 2 3] inc) # yields [2 3 4]
```

`func` should never be called without the bang modifier, as its behaviour is
undefined in thatcase and the return value is unusable.

`func!` always returns `null`.

## Context

Since defining a function is literally a side effect, `func!` must always be
called with the bang modifier.
