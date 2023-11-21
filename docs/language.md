# The Rudi Language

Rudi is a Lisp dialect. Each Rudi program consists of a series of statement, which are tuples
(`( ... )`), which are executed in series.

## Data Types

Rudi knows the following data types:

* Null (`null`)
* Bool (`true` or `false`)
* Number (int64 like `42` or float64 like `42.42`)
* Strings (`"i am a string"`)
* Vectors (`[1 2 3]`, items are whitespace separated, but commas can also be used, like `[1, 2, 3]`;
  each element of a vector can be any type of expression, e.g. `[1 (+ 1 2)]` to create `[1 3]`)
* Objects (`{"key" "value" "otherkey" "othervalue"}`; keys and values can be any type of expressions,
  but the key expression must return a string; keys can also be identifiers, i.e. unquoted, so
  `{foo "bar"}` is the same as `{"foo" "bar"}`)

As you can see, Rudi is basically JSON extended with S-expressions.

In addition to these literal data types, Rudi understands a number of expressions:

## Expressions

### Statements

Statements are the top-level elements of an Rudi program and in fact just an alias for tuples. The
term "statement" is simply used to denote a top-level tuple.

```
(do something)
(do other things)
(add this to that)
```

## Tuples

A tuple always consists of the literal `(`, followed by one or more expressions, followed by a `)`.
Between tuples, any amount of whitespace is allowed. Each of the expressions are separated from another
by whitespace.

```
(do something)
(do other things)
(add this to that)
```

Tuples represent "function calls". The first element of every tuple must be an **identifier**, which
is an unquoted string that refers to a built-in function. For example, the `to-upper` function can be
called by writing

```
(to-upper "foo")
```

In general, any kind of expression can follow function names inside tuples. The number of expressions
depends on the function (`to-upper` requires exactly 1 argument, `concat` takes 2 or more arguments).

Since tuples are expressions, tuples can be nested:

```
(to-upper (to-lower "FOO"))
```

In the example above, the string `"FOO"` would first be lowercased, and the result of the inner
tuple (`"foo"`) would then be uppercased.

Tuples (i.e. functions) can return any of the known data types (numbers, strings, ...), but not
other expressions (a tuple cannot return an identifier, for example). This means the function name
cannot by dynamic, you cannot do `((concat "-" "to" "upper") "foo")` to call `(to-upper "foo")`.

## Symbols

Symbols are either variables or bare path expressions that reference the global document.

### Variables

Rudi has support for runtime variables. A variable holds any of the possible data types and has a
unique, case-sensitive name, like `$myVar`. Symbols are expressions and can therefore be used in
most places:

```
(add $myVar 5)
(concat "foo" [1 $var 2])
(set $myVar (+ 1 4))
```

A path expression can follow a variable name, allowing easy access to sub fields:

```
(set $var [1, 2, 3])
(print $var[0])

(set $var {foo [1, 2, {foo "bar"}]})
(print $var.foo[2].foo)
```

See further down for more details on path expressions.

### Global Document Access

Rudi programs are meant to transform a document (usually an `Object`). To make it easy to access
the document and sub fields within, path expressions can be used just like with variables:

Suppose a JSON document like

```json
{
   "foo": "bar",
   "list": [1, 2, 3]
}
```

is being processed by an Rudi program, then you could write

```
(print .foo)
(print .list[1])
```
