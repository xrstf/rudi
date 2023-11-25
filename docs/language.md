# The Rudi Language

Rudi is a Lisp dialect. Each Rudi program consists of a series of statement, which are tuples
(`( ... )`), path expressions (`.foo.bar[0]`) or literal values (`{foo "bar"}`). Statements are
separated by whitespace. Each of these statements is evaluated in sequence, with the result of the
last statement being the result of the entire program.

```lisp
(set! .foo 42)
(set! $var (+ .bar (len .users)))
(map [1 2 3] [x] (+ $x $var))
```

Using path expressions or literal values as statements is usually only useful if your program ends
with a literal (often a constructed object) or if you use a path expression to select a sub-value
out of the global document.

Most functions can return errors and if an error happens, program execution stops unless the error
is caught using the `try` function.

## Data Types

Rudi knows the following data types, known as "literals":

* Null (`null`)
* Bools (`true` or `false`)
* Numbers (int64 like `42` or float64 like `42.42`)
* Strings (`"i am a string"`)
* Vectors (`[1 2 3]`, items are whitespace separated, but commas can also be used, like `[1, 2, 3]`;
  each element of a vector can be any type of expression, e.g. `[1 (+ 1 2)]` to create `[1 3]`)
* Objects (`{"key" "value" "otherkey" "othervalue"}`; keys and values can be any type of expressions,
  but the key expression must return a string; keys can also be identifiers, i.e. unquoted, so
  `{foo "bar"}` is the same as `{"foo" "bar"}`)

As you can see, Rudi is basically JSON extended with S-expressions.

In addition to these literal data types, Rudi understands a number of expressions:

## Expressions

### Null

Nulls are the empty value and are always written as `null`, like in JSON. Nulls equal each other,
so `(eq? null null)` is true.

### Bool

Booleans are the thruth values of `true` (positive) and `false` (negative).

### Number

Numbers are either 64-bit integers (int64) or 64-bit floating-points (float64). Whole numbers are
`0` or `183732627` or `-17`. Floating point numbers are integers followed by a decimal fraction, like
`0.4` or `-42.341`.

Integers and floats are not directly comparable, convert the int to a float to perform comparisons.
Even `0` is not equal to `0.0` without conversion.

### String

Strings are lists of characters, like `"hello world"`. There is no separate type in Rudi for single
bytes. Strings can be empty (`""`) and leading/trailing whitespace is kept (`"  foo  "` will not be
trimmed automatically). Use `\"` to write a literal `"` inside a string, use `\\` to write a literal
`\` (e.g. `"C:\\dos\\run"`).

### Vector

Vectors are an ordered list of items. A vector starts with the literal `[` followed by a whitespace
separated list of expressions, followed by a literal `]`. Vectors can be empty (`[]`). Instead of
whitespace, commas can also be used, e.g. `[1,2,3]`.

The items in a vector do not need to share the same datatype, so `[1 "foo" true]` (number, string,
number) is valid.

### Object

Objects are unordered sets of key-value pairs. They always begin with a literal `{`, followed by an
arbitrary number of key-value pairs. Key and value are each any possible expressions, with limitations
listed below. Object declarations end with a literal `}`, for example:

```lisp
{"foo" "bar" "secondkey" "secondvalue"}
```

```lisp
{
  (to-upper "foo") (+ 1 2)
}
```

Keys and values are separated by whitespace (not with a `:` like in JSON). Likewise, key-value pairs
are separated from each other by whitespace. In effect, each object declaration needs to have an
even number of expression in it.

Empty objects are permitted (`{}`).

Note that objects are internally unordered and functions like [map](functions/lists-map.md) or
[range](functions/lists-range.md) will have a random iteration order. Due to the fact that Go's
JSON encoder sorts keys alphabetically when writing JSON, this should rarely be of concern.

### Statement

Statements are the top-level elements of an Rudi program and this distinction is mostly useful
internally to detect the top-level of a Rudi program. A statement can be a tuple (i.e. a function
call), a symbol (a variable and/or path expression) or a literal value.

```lisp
(set! .foo 42)
(set! $var (+ .bar (len .users)))
(map [1 2 3] [x] (+ $x $var))
```

### Tuple

A tuple always consists of the literal `(`, followed by one or more expressions, followed by a `)`.
Within tuples, any amount of whitespace is allowed. Each of the expressions are separated from another
by whitespace.

```lisp
(do something)
(do other things)
(add this to that)
```

Tuples represent "function calls". The first element of every tuple must be an **identifier**, which
is an unquoted string that refers to a built-in function. For example, the `to-upper` function can be
called by writing

```lisp
(to-upper "foo")
```

In general, any kind of expression can follow function names inside tuples. The number of expressions
depends on the function (`to-upper` requires exactly 1 argument, `concat` takes 2 or more arguments).
However there are special cases where only certain kinds of expressions are allowed, like on
functions with the bang modifier, which must use a symbol as their first argument (like
`(set! $var 42`).

Since tuples are expressions, tuples can be nested:

```lisp
(to-upper (to-lower "FOO"))
```

In the example above, the string `"FOO"` would first be lowercased, and the result of the inner
tuple (`"foo"`) would then be uppercased.

Tuples (i.e. functions) can return any of the known data types (numbers, strings, ...), but not
other expressions (a tuple cannot return an identifier, for example). This means the function name
cannot by dynamic, you cannot do `((concat "-" "to" "upper") "foo")` to call `(to-upper "foo")`.

### Bang Modifier

Functions in Rudi are stateless, meaning they compute a value and return it, without any side
effects. However that alone would be boring and not really helpful, so Rudi breaks the pure
functional approach and allows side effects.

For example, to set (define or update) a variable, the expression `(set $var 42)` would not do what
you might think: _Within_ the tuple, the `set` will define the variable and the its value to `42`,
but this will not affect the _next_ tuple after it. So for example the program `(set $var 42) $var`
will error out because `$var` is not defined in the second statement.

To "make changes stick", use the bang modifier (`!`): `(set! $var 42) $var` will actually return
`42`. The bang modifier can be used on any function, as long as the first argument is a symbol
(i.e. a variable or a bare path expression). If the modifier is used, the result of the function
expression is updated in the symbol, so `(set! $var 42)` will first calculate the desired value
(`42`, easy in this case) and then actually set the variable to that value.

The difference is more obvious with other functions:

```lisp
(set! $var "foo")      # assign "foo" to $var

(append $var "bar")    # returns a new string "foobar" without updating $var
$var                   # will print "foo"

(append! $var "bar")   # returns a new string "foobar" and updates $var
$var                   # will print "foobar"
```

Since the bang modifier makes a function modify the first argument, expressions like
`(append! "foo" "bar")` are not valid, as it's not clear where the intended side effect should go.

The bang modifier can be used with path expressions to set deeper values:

```lisp
(set! $var {foo "bar"}) # {"foo": "bar"}
(set! $var.foo "new")   # "new"
$var                    # {"foo": "new"}

(set! $var [1 2 3])    # [1 2 3]
(append! $var 4)       # 4
(set! $var[3] 5)       # 5
$var                   # [1 2 3 5]
```

The bang modifier can be used with any function, though you will be hardpressed to find meaningful
examples of `(if! .path.expr 42)` or `(eq?! .path 42)`.

### Symbol

Symbols are either variables or bare path expressions that reference the global document.

#### Variables

Rudi has support for runtime variables. A variable begins with the literal `$`, followed by a
case-sensitive name, like `$myVar`. Variables hold any of the possible data types.

Symbols are expressions and can therefore be used in most places:

```lisp
(add $myVar 5)
(concat "foo" [1 $var 2])
(set! $myVar (+ 1 4))
```

A path expression can follow a variable name, allowing easy access to sub fields:

```lisp
(set! $var [1, 2, 3])
(print $var[0])

(set! $var {foo [1, 2, {foo "bar"}]})
(print $var.foo[2].foo)
```

See further down for more details on path expressions.

#### Global Document Access

Rudi programs are meant to transform a document (usually an `Object`). To make it easy to access
the document and sub fields within, you can reference the global document by a single dot (`.`),
optionally (and often) with a path expression on it, like `.foo.bar`. So `.foo` would reference the
field `"foo"` in the global document, whereas `$var.foo` is the field `"foo"` in the variable `$var`.

Suppose a JSON document like

```json
{
   "foo": "bar",
   "list": [1, 2, 3]
}
```

is being processed by an Rudi program, then you could write

```lisp
.foo
.list[1]
```

to first select `"bar"`, then `2`. Bare path expressions work like variables, you can do anything
you can do with variables, like:

```lisp
(append! .list 4)
(to-upper! .foo)
```

#### Scopes

In general, side effects (i.e. functions with bang modifier) affect all following sibling and child
expressions, but not the parent. This is like doing

```go
foo := 42

if condition {
  foo := 7
}

println(foo) // prints 42
```

in Go. Variables are meant to be helpers and are so scoped to the scope where they are defined:

```lisp
(set! $var 42)
$var             # 42

# This set! function would set the value for the entire positive branch of the "if" tuple,
# but it will not leak outside of the "if".

(if true
  (set! $var "new-value"))    # "new-value"
$var                          # 42

# ... but the new variable is valid in its scope.

(if true
  (do
    (set! $var "new-value")
    (append $var "-suffix"))) # "new-value-suffix"
$var                          # 42
```

The exception from this rule is the global document. As the name implies, it is meant to be global
and to allow for effective, readable Rudi code, there is only one document and it can be modified
from anywhere.

If a Rudi program was loaded with

```json
{
   "foo": "bar",
   "list": [1, 2, 3]
}
```

then

```lisp
(if true (set! .foo "new-value"))    # "new-value"
.foo                                 # "new-value"

(if true (append! .list 4))    # 4
.list                          # [1 2 3 4]
```

### Path Expression

Rudi implements simple JSONPath-like expressions to allow descending into deeply nested objects.
Each path consists of a series of steps, with each step being either an object step (e.g. `.foo`)
or a vector step (e.g. `[42]`). Steps can be chained, like `.foo[42].bar.sub[1][2]`.

Path steps can also be computed (`[(+ 1 42)]`), which allows to use more complex expressions to
form steps, like `["string.with.dot"]` or even `[$var.index]`.

There is one special case: Paths that start with a vector step on the global document: For a variable
this would look like `$var[42]`, but for the global document this would be just `[42]`, which is
indistinguishable from "a vector with 1 element, the number 42". To resolve this ambiguity, bare
path expressions that start with a vector step must have a leading dot, like `.[42]`.

Path expressions must be traversable, or else an error is returned: Trying to descend with `.foo`
into a vector would result in an error, likewise using `[3]` to descend into a string is an error.
Use the [`has?`](functions/core-has.md) and [`try`](functions/core-try.md) functions to deal with
possibly misfitting path expressions.

Path expressions can be used on

* Symbols (`$var.foo` or `.document.key`)
* Vector nodes (`[1 2 3][1]`, first step of the path must be a vector step, i.e. `[1 2].foo` is invalid)
* Object nodes (`{foo "bar"}.foo`, first step of the path must be an object step, i.e. `{foo "bar"}[0]`
  is invalid)
* Tuples (`(map $obj to-upper).key`, requires that the tuple evaluated to a vector or object that
  can be traversed, otherwise an error is returned (e.g. `(+ 1 2).key` is invalid))

The evaluated value of any of these expressions is always _with_ the path expression applied, so
for example in `(+ (process $obj).userCount 32)` the `+` function will see 2 arguments like
`(+ $userCount 32)` because when processing the `process` tuple, the path expression on it is also
evaluated already.
